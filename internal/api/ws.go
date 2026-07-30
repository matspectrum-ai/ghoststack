package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type WSMessage struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data,omitempty"`
}

type WSClient struct {
	conn *websocket.Conn
	send chan []byte
	mu   sync.Mutex
}

type WSHub struct {
	mu       sync.RWMutex
	clients  map[*WSClient]bool
	upgrader websocket.Upgrader
}

func NewWSHub() *WSHub {
	return &WSHub{
		clients: make(map[*WSClient]bool),
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
}

func (h *WSHub) HandleWS(w http.ResponseWriter, r *http.Request) {
	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	client := &WSClient{
		conn: conn,
		send: make(chan []byte, 100),
	}

	h.mu.Lock()
	h.clients[client] = true
	h.mu.Unlock()

	go h.writePump(client)
	go h.readPump(client)
}

func (h *WSHub) readPump(client *WSClient) {
	defer func() {
		h.mu.Lock()
		delete(h.clients, client)
		h.mu.Unlock()
		client.conn.Close()
	}()

	client.conn.SetReadLimit(4096)
	client.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	client.conn.SetPongHandler(func(string) error {
		client.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, _, err := client.conn.ReadMessage()
		if err != nil {
			break
		}
	}
}

func (h *WSHub) writePump(client *WSClient) {
	ticker := time.NewTicker(30 * time.Second)
	defer func() {
		ticker.Stop()
		client.conn.Close()
	}()

	for {
		select {
		case message, ok := <-client.send:
			client.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				client.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := client.conn.WriteMessage(websocket.TextMessage, message); err != nil {
				return
			}
		case <-ticker.C:
			client.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := client.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (h *WSHub) Broadcast(msg WSMessage) {
	data, err := json.Marshal(msg)
	if err != nil {
		return
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	for client := range h.clients {
		select {
		case client.send <- data:
		default:
			close(client.send)
			delete(h.clients, client)
		}
	}
}

func (h *WSHub) BroadcastEvent(eventType string, payload any) {
	data, err := json.Marshal(payload)
	if err != nil {
		return
	}
	h.Broadcast(WSMessage{
		Type: eventType,
		Data: data,
	})
}

type MetricsCollector struct {
	mu      sync.RWMutex
	hub     *WSHub
	cpu     float64
	memory  uint64
	rxBytes int64
	txBytes int64
	uptime  int64
}

func NewMetricsCollector(hub *WSHub) *MetricsCollector {
	return &MetricsCollector{
		hub: hub,
	}
}

type MetricsSnapshot struct {
	CPU     float64 `json:"cpu"`
	Memory  uint64  `json:"memory"`
	RXBytes int64   `json:"rx_bytes"`
	TXBytes int64   `json:"tx_bytes"`
	Uptime  int64   `json:"uptime"`
}

func (mc *MetricsCollector) Collect(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			mc.snapshot()
		}
	}
}

func (mc *MetricsCollector) snapshot() {
	mc.mu.RLock()
	snap := MetricsSnapshot{
		CPU:     mc.cpu,
		Memory:  mc.memory,
		RXBytes: mc.rxBytes,
		TXBytes: mc.txBytes,
		Uptime:  mc.uptime,
	}
	mc.mu.RUnlock()

	mc.hub.BroadcastEvent("metrics", snap)
}

func (mc *MetricsCollector) Update(cpu float64, memory uint64, rx, tx int64) {
	mc.mu.Lock()
	defer mc.mu.Unlock()
	mc.cpu = cpu
	mc.memory = memory
	mc.rxBytes = rx
	mc.txBytes = tx
	mc.uptime = time.Now().Unix()
}

func (mc *MetricsCollector) Snapshot() MetricsSnapshot {
	mc.mu.RLock()
	defer mc.mu.RUnlock()
	return MetricsSnapshot{
		CPU:     mc.cpu,
		Memory:  mc.memory,
		RXBytes: mc.rxBytes,
		TXBytes: mc.txBytes,
		Uptime:  mc.uptime,
	}
}

func providerSSEHandler(hub *WSHub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "streaming unsupported", http.StatusInternalServerError)
			return
		}

		ctx := r.Context()
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		fmt.Fprintf(w, "data: {\"status\":\"connected\"}\n\n")
		flusher.Flush()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				fmt.Fprintf(w, "data: {\"type\":\"heartbeat\",\"ts\":%d}\n\n", time.Now().Unix())
				flusher.Flush()
			}
		}
	}
}
