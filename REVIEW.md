---
phase: code-review
reviewed: 2026-07-30T12:00:00Z
depth: deep
files_reviewed: 18
files_reviewed_list:
  - internal/security/sandbox.go
  - internal/security/secureboot.go
  - internal/security/killswitch_impl.go
  - internal/security/audit_impl.go
  - internal/security/audit_structured.go
  - internal/api/ws.go
  - internal/api/server.go
  - internal/cli/start.go
  - internal/cli/diagnose.go
  - internal/plugins/ipc.go
  - internal/plugins/manager.go
  - internal/providers/wireguard.go
  - internal/providers/process.go
  - internal/providers/socks5.go
  - internal/providers/engine.go
  - internal/platform/linux/firewall.go
  - internal/networking/tun.go
  - internal/monitoring/metrics.go
findings:
  critical: 21
  warning: 15
  info: 5
  total: 41
status: issues_found
---

# Phase: Code Review Report

**Reviewed:** 2026-07-30T12:00:00Z
**Depth:** deep
**Files Reviewed:** 18
**Status:** issues_found

## Summary

This review covers 18 Go source files totaling ~2,800 lines. The codebase has significant issues across all categories: **21 critical** (including data races, deadlocks, compilation errors, command injection, nil pointer in seccomp, SOCKS5 protocol violation), **15 warnings** (ignored errors, missing headers, fragile patterns), and **5 info** items. The most severe systemic problems are in the TUN/networking layer (unsafe pointer misuse, wrong struct offsets causing runtime panics), the firewall (command injection through shell interpolation), the SOCKS5 implementation (completely broken address parsing), and the plugin subsystem (mutex deadlock, undefined variable).

---

## Critical Issues

### CR-01: Seccomp filter pointer is always NULL — sandbox is non-functional

**File:** `internal/security/sandbox.go:202-204`
**Issue:** `unsafePtr` returns `uintptr(0)` instead of the actual pointer to the seccomp BPF program. This means `PR_SET_SECCOMP` at line 126 receives a null filter pointer. The seccomp sandbox silently never does anything — all syscalls remain allowed. The function body is:

```go
func unsafePtr(p *unix.SockFprog) uintptr {
    return uintptr(0)   // <-- should be uintptr(unsafe.Pointer(p))
}
```

**Fix:**
```go
func unsafePtr(p *unix.SockFprog) uintptr {
    return uintptr(unsafe.Pointer(p))
}
```

---

### CR-02: `Broadcast` mutates client map while holding only RLock (data race)

**File:** `internal/api/ws.go:114-131`
**Issue:** `Broadcast()` iterates over `h.clients` under an `RLock` but calls `delete(h.clients, client)` when a client's send buffer is full. Writing to a map while other goroutines may be reading or writing it is a data race (undefined behavior). The `readPump` and `writePump` goroutines also read and write this map.

**Fix:** Upgrade to a write lock before mutating the map, or defer the deletion:
```go
func (h *WSHub) Broadcast(msg WSMessage) {
    data, err := json.Marshal(msg)
    if err != nil {
        return
    }
    h.mu.Lock() // Write lock needed
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
```

---

### CR-03: `Unload` calls `Disable` while holding mutex — deadlock

**File:** `internal/plugins/ipc.go:146-158`
**Issue:** `Unload()` acquires `p.mu.Lock()` (line 148), then calls `p.Disable(ctx)` (line 150). `Disable()` also tries to acquire `p.mu.Lock()` (line 121). Go's `sync.Mutex` is not reentrant — this causes a self-deadlock. The goroutine hangs forever.

**Fix:** Release the lock before calling Disable, or inline the Disable logic:
```go
func (p *subprocessPlugin) Unload(ctx context.Context) error {
    p.mu.Lock()
    // copy state
    p.mu.Unlock()

    p.Disable(ctx)  // no longer holding lock

    p.mu.Lock()
    defer p.mu.Unlock()
    if p.socketPath != "" {
        os.Remove(p.socketPath)
    }
    p.state = PluginStateRemoved
    return nil
}
```

---

### CR-04: Undefined variable `loader` in `manager.go:Load` — compilation error

**File:** `internal/plugins/manager.go:96`
**Issue:** `loader` is referenced on line 96 but is not defined anywhere in the function scope or at package level. The `defaultPluginLoader` type exists (line 241) but is never instantiated here. This code will not compile.

```go
m.plugins[pluginPath] = &pluginEntry{
    ...
    loader:  loader,  // <-- undefined
}
```

**Fix:**
```go
m.plugins[pluginPath] = &pluginEntry{
    state:    PluginStateLoaded,
    plugin:   plugin,
    manifest: manifest,
    loader:   &defaultPluginLoader{},
}
```

---

### CR-05: Map lookup key mismatch — `Initialize` always fails to find plugin

**File:** `internal/plugins/manager.go:108-115`
**Issue:** In `Load()` (line 92), the plugin is stored with key `pluginPath` (an absolute directory path). In `Initialize()` (line 110), the lookup uses `manifest.Entry` (a bare binary name like "my-plugin"). These are different values, so the lookup always returns `exists == false`, and `Initialize()` always returns "plugin not loaded".

**Fix:** Use `manifest.ID` as the key consistently, or store the lookup key when loading:
```go
// In Load:
m.plugins[manifest.ID] = &pluginEntry{...}

// In Initialize:
entry, exists := m.plugins[manifest.ID]
```

---

### CR-06: HTTP server shuts down after 2 seconds

**File:** `internal/cli/start.go:71-79` and `internal/api/server.go:46-49`
**Issue:** `apiCtx` is created with a 2-second timeout (`context.WithTimeout(cmd.Context(), 2*time.Second)`). This context is passed to `server.Start()`, which spawns a goroutine that calls `server.Close()` when `apiCtx.Done()` fires. After 2 seconds, the HTTP server shuts down. The application will serve requests for only 2 seconds.

**Fix:** Use the parent `cmd.Context()` (with no artificial timeout) for the server lifecycle:
```go
httpServer, err := server.Start(cmd.Context(), apiAddr)
```

---

### CR-07: `providerCfg` always nil — user configuration never passed to providers

**File:** `internal/cli/start.go:34-59`
**Issue:** `providerCfg` is declared as `var providerCfg map[string]any` (zero value = nil) and is never populated from the loaded config file. The dead code at lines 44-46 reads the config file again but assigns to `_`. So `engine.Start(startCtx, providerName, providerCfg)` always passes a nil config, and every provider defaults to zero values.

**Fix:** Extract provider config from the loaded `cfg`:
```go
cfg, err := config.Load(configPath)
if err != nil { ... }
providerCfg = cfg.Profiles[providerName].Config  // populate from config
```

---

### CR-08: SOCKS5 address parsing reads RSV and ATYP as port — completely wrong

**File:** `internal/providers/socks5.go:91-94`
**Issue:** The SOCKS5 request format is: `[ver, cmd, rsv=0x00, atyp, addr..., port_hi, port_lo]`. The code reads the port from `buf[2]` and `buf[3]`, but `buf[2]` is the RSV field (always 0x00) and `buf[3]` is the ATYP field (0x01=IPv4, 0x03=domain, 0x04=IPv6). The port is always decoded as 0 or a small number depending on ATYP. Additionally, the code only handles IPv4 addresses — domain names and IPv6 addresses produce garbage.

**Fix:** Handle address type properly:
```go
atyp := buf[3]
var addr string
switch atyp {
case 1: // IPv4
    addr = net.IP(buf[4:8]).String()
    port := binary.BigEndian.Uint16(buf[8:10])
    addr = net.JoinHostPort(addr, strconv.Itoa(int(port)))
case 3: // Domain name
    domainLen := int(buf[4])
    addr = string(buf[5 : 5+domainLen])
    port := binary.BigEndian.Uint16(buf[5+domainLen : 7+domainLen])
    addr = net.JoinHostPort(addr, strconv.Itoa(int(port)))
case 4: // IPv6
    addr = net.IP(buf[4:20]).String()
    port := binary.BigEndian.Uint16(buf[20:22])
    addr = net.JoinHostPort(addr, strconv.Itoa(int(port)))
default:
    return
}
```

---

### CR-09: SOCKS5 bidirectional relay deadlocks on half-close

**File:** `internal/providers/socks5.go:104-106`
**Issue:** `handleConn` launches `go relay(target, conn)` (reads target→conn) and then calls `relay(conn, target)` (reads conn→target). If the target closes the connection first, the first relay goroutine returns, but the second is blocked reading from `conn` (which the client hasn't closed yet). `conn.Close()` is deferred in `handleConn` but never runs because `handleConn` is blocked waiting for the second relay to return. Result: deadlocked goroutine leak.

**Fix:** Use `io.Copy` with synchronized close, or use a channel to signal closure:
```go
func (p *socks5Provider) handleConn(ctx context.Context, conn net.Conn) {
    defer conn.Close()
    // ... handshake ...
    var wg sync.WaitGroup
    wg.Add(2)
    go func() {
        defer wg.Done()
        relay(target, conn)
        target.Close()
    }()
    go func() {
        defer wg.Done()
        relay(conn, target)
        conn.Close()
    }()
    wg.Wait()
}
```

---

### CR-10: Netmask calculation panics on negative shift count

**File:** `internal/networking/tun.go:225-235`
**Issue:** When `maskVal > 8`, the code computes `0xff << (8 - maskVal)` which is a negative shift. In Go, shifting by a negative count causes a runtime panic:

```go
maskBytes[i] = byte(0xff << (8 - maskVal)) // maskVal=24 → 0xff << -16 → PANIC
```

This panics for any netmask larger than /8.

**Fix:** Use the correct netmask calculation:
```go
ones, _ := ipnet.Mask.Size()
var maskBytes [4]byte
for i := 0; i < 4; i++ {
    if ones >= 8 {
        maskBytes[i] = 0xff
        ones -= 8
    } else if ones > 0 {
        maskBytes[i] = byte(0xff << (8 - ones))
        ones = 0
    }
}
```

---

### CR-11: TUN `Addresses()` reads wrong byte offset for IP address

**File:** `internal/networking/tun.go:186-189`
**Issue:** `SIOCGIFADDR` returns the address in a `struct sockaddr` starting at byte 16 of the `ifreq` struct (after the 16-byte interface name). The `struct sockaddr` layout is: 2 bytes family + 14 bytes data. For AF_INET, the IPv4 address starts at byte 18. The code reads `ifr[4:8]` which is in the middle of the interface name.

**Fix:**
```go
// After ioctl, read address from correct offset
// struct ifreq: 16 bytes name + struct sockaddr { 2 bytes family + data }
ipBytes := ifr[20:24] // 16 + 2 (family) + 2 (port) = offset 20 for IPv4 addr
ip := net.IP(ipBytes).String()
```

---

### CR-12: TUN `Delete()` only closes file descriptor, doesn't destroy interface

**File:** `internal/networking/tun.go:274-287`
**Issue:** Closing the TUN file descriptor does not destroy the TUN interface on Linux. The interface (`ghost0`) persists until explicitly deleted via `ip link delete`. This leaks TUN interfaces — after multiple start/stop cycles, stale interfaces accumulate.

**Fix:** Execute `ip link delete <name>` or use an ioctl:
```go
func (t *tunDevice) Delete(ctx context.Context) error {
    t.mu.Lock()
    defer t.mu.Unlock()
    if t.fd < 0 {
        return nil
    }
    unix.Close(t.fd)
    exec.CommandContext(ctx, "ip", "link", "delete", t.name).Run()
    t.fd = -1
    t.name = ""
    t.up = false
    return nil
}
```

---

### CR-13: Unix socket listener closed before subprocess can connect

**File:** `internal/plugins/ipc.go:84-102`
**Issue:** A Unix domain socket listener is created (line 84), then the subprocess is started with `--socket` pointing to that socket path (line 93-99). Immediately after, the listener is closed (line 102). The subprocess hasn't started yet, so it will never be able to connect — the socket file is already gone. The `conn` field and `reader` field of `subprocessPlugin` are never populated.

**Fix:** Keep the listener alive or use a different IPC mechanism (stdin/stdout pipe):
```go
// Either pass the listener's fd to the child, or
// keep accepting connections in a goroutine:
go func() {
    conn, err := listener.Accept()
    if err != nil { return }
    p.conn = conn
    p.reader = bufio.NewScanner(conn)
    // ... handle IPC ...
}()
```

---

### CR-14: Firewall applies nftables setup without piping stdin

**File:** `internal/platform/linux/firewall.go:50-66`
**Issue:** The first `nft -f -` command (line 60) runs with no stdin connected. The `setup` variable (containing the actual nftables rules) is never passed to it. The command may fail or silently do nothing. If it fails, the function returns an error and the real setup via `echo ... | nft -f -` (line 64-66) is never executed.

**Fix:** Use proper stdin piping:
```go
cmd := exec.CommandContext(ctx, "nft", "-f", "-")
cmd.Stdin = strings.NewReader(setup)
if err := cmd.Run(); err != nil {
    return fmt.Errorf("nft setup: %w", err)
}
```
Also eliminate the echo + shell approach — use direct stdin piping.

---

### CR-15: Command injection via string interpolation in firewall shell commands

**File:** `internal/platform/linux/firewall.go:64-66, 80-91, 163-176`
**Issue:** Multiple places construct shell commands by string interpolation with user- or config-derived values, then pass them to `sh -c`:

```go
echo := fmt.Sprintf("echo '%s' | nft -f -", setup)                     // line 64
cmd = fmt.Sprintf("nft add rule inet ghoststack output ip daddr %s accept", dest)  // line 80
rule := "nft add rule inet ghoststack output oif != \"" + f.iface + "\" drop" // line 165
```

If any value contains a single quote, semicolon, or backtick, arbitrary commands execute. The `iface` field is derived from config, and `setup` contains nftables syntax that may embed user data.

**Fix:** Use direct command+args form:
```go
exec.CommandContext(ctx, "nft", "add", "rule", "inet", "ghoststack", "output", "ip", "daddr", dest, "accept")
```

---

### CR-16: All iptables command errors silently ignored

**File:** `internal/platform/linux/firewall.go:103-122`
**Issue:** Every `exec.CommandContext(...).Run()` call in `applyIptables` discards its error return. If the GHOSTSTACK chain already exists, if `-I OUTPUT` fails, or if any rule cannot be added, the function returns `nil` (success). The firewall may be in an incomplete state without any indication.

**Fix:** Check every error:
```go
if err := exec.CommandContext(ctx, ipt, "-N", "GHOSTSTACK").Run(); err != nil {
    // Chain may already exist; that's OK, but other errors are not
    ...
}
```

---

### CR-17: `process.go:Stop()` and `wait()` goroutine race on state

**File:** `internal/providers/process.go:86-94, 97-122`
**Issue:** When `Stop()` is called:
1. `cancel()` kills the process
2. `Stop()` acquires the lock, sets `pm.state = ProcessStopped` (line 120), releases lock
3. The `wait()` goroutine (started at line 81) finishes `cmd.Wait()` and acquires the lock
4. `wait()` sets `pm.state = ProcessFailed` (line 90) — overwriting `ProcessStopped`

After Stop returns, the state may be `ProcessFailed` instead of `ProcessStopped`. Future callers checking `State()` get the wrong value.

**Fix:** Check in `wait()` whether `Stop()` was called:
```go
func (pm *ProcessManager) wait(cmd *exec.Cmd) {
    err := cmd.Wait()
    pm.mu.Lock()
    defer pm.mu.Unlock()
    if pm.state != ProcessStopped { // Don't overwrite intentional stop
        pm.state = ProcessFailed
    }
    pm.err = err
    pm.pid = 0
    close(pm.done)
}
```

---

### CR-18: Metrics `collect()` is a no-op

**File:** `internal/monitoring/metrics.go:160-163`
**Issue:** The `collect()` method acquires the lock and immediately releases it. No actual metrics are collected — CPU, memory, and network I/O are never read from the system. The hardcoded values in `start.go:93` never change.

**Fix:** Implement actual metric collection (e.g., reading /proc/stat, /proc/meminfo, net/io counters):
```go
func (mc *MetricsCollector) collect() {
    mc.system.mu.Lock()
    defer mc.system.mu.Unlock()
    // Read CPU from /proc/stat, memory from /proc/meminfo, etc.
    mc.system.CPU = readCPUUsage()
    mc.system.Memory = readMemoryUsage()
}
```

---

### CR-19: `audit_structured.go` flush can cause duplicate log entries on error

**File:** `internal/security/audit_structured.go:74-89`
**Issue:** When `flush()` partially writes entries to disk and then encounters an error (line 83-84), it returns without clearing `l.entries` (line 87 is skipped). On the next call to `Log()` (which triggers another `flush()`), ALL entries are re-flushed — including those already written. This produces duplicate entries in the audit log file.

**Fix:** Track what was flushed or clear entries incrementally:
```go
func (l *StructuredAuditLogger) flush() error {
    remaining := l.entries
    l.entries = nil // Clear first; re-add on failure
    // ... write remaining to file ...
    if err != nil {
        l.entries = append(remaining, l.entries...) // restore on failure
        return err
    }
    return nil
}
```

---

### CR-20: `ipc.go` socket file not removed on error in Initialize

**File:** `internal/plugins/ipc.go:84-102`
**Issue:** The socket file is created by `net.Listen("unix", socketPath)` (line 84). If `p.proc.Start()` fails (line 93), the socket file is removed (line 98) — correct. But if `listener.Close()` at line 102 succeeds, the socket file persists on disk. When `Initialize` is called again, `net.Listen` will fail because the socket file already exists. The socket file is never cleaned up unless `Disable` is called.

**Fix:** Clean up the socket file after closing the listener:
```go
listener.Close()
os.Remove(socketPath) // Clean up even on success
```

---

### CR-21: `applyToProcess` TOCTOU + non-functional network isolation

**File:** `internal/security/sandbox.go:206-227`
**Issue:** Two bugs: (1) TOCTOU race — process existence is checked at line 208, but the process could exit before the isolation is applied at line 213. (2) `applyNetworkIsolation` calls `os.Chmod("/proc/pid/ns/net", 0)` which sets permissions on the *namespace file*, not the process's namespace membership. This cannot achieve network isolation — the process retains full network access. Achieving network isolation requires moving the process to a new network namespace.

**Fix:** Use `unix.Setns()` to move the process to a restricted network namespace, or use network cgroups.

---

## Warnings

### WR-01: `BroadcastEvent` silently discards JSON marshal errors

**File:** `internal/api/ws.go:134`
**Issue:** `data, _ := json.Marshal(payload)` — the error is discarded. If `payload` contains non-serializable types (e.g., channels, functions), the marshal fails silently and clients receive empty/invalid messages.

**Fix:** Log the error or return it:
```go
data, err := json.Marshal(payload)
if err != nil {
    log.Printf("broadcast marshal: %v", err)
    return
}
```

---

### WR-02: `audit_impl.go:List` returns oldest entries instead of most recent

**File:** `internal/security/audit_impl.go:51`
**Issue:** `copy(out, a.entries[:limit])` returns entries from index 0 (oldest). Audit consumers almost always want the most recent entries. The `StructuredAuditLogger.List` (audit_structured.go:57-72) correctly returns from the end.

**Fix:**
```go
start := len(a.entries) - limit
if start < 0 {
    start = 0
}
copy(out, a.entries[start:])
```

---

### WR-03: Audit log file created world-readable (0644)

**File:** `internal/security/audit_structured.go:75`
**Issue:** `os.OpenFile(l.output, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)` creates the audit log with world-readable permissions. Audit logs may contain sensitive information.

**Fix:** Use 0600:
```go
f, err := os.OpenFile(l.output, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
```

---

### WR-04: HTTP handlers missing `Content-Type: application/json`

**Files:** `internal/api/server.go:73, 87, 112, 122, 129, 136`
**Issue:** `json.NewEncoder(w).Encode(...)` is called without setting the Content-Type header. While browsers and many clients infer JSON from content, explicit headers are required for correctness.

**Fix:** Either set the header before encoding, or use a wrapper:
```go
w.Header().Set("Content-Type", "application/json")
json.NewEncoder(w).Encode(data)
```

---

### WR-05: `server.Close()` used instead of graceful `Shutdown()`

**File:** `internal/api/server.go:48`
**Issue:** `server.Close()` immediately terminates all connections without waiting for in-flight requests to complete. This can drop active WebSocket connections and HTTP responses.

**Fix:**
```go
go func() {
    <-ctx.Done()
    shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    server.Shutdown(shutdownCtx)
}()
```

---

### WR-06: `wg-quick down` and `killSwitch.Disable` errors ignored

**File:** `internal/providers/wireguard.go:169, 180`
**Issue:** Both `exec.CommandContext(ctx, "wg-quick", "down", p.configPath).Run()` and `p.killSwitch.Disable(ctx)` discard their error returns. If wg-quick fails or the kill switch cannot be disabled, the cleanup is incomplete and the user receives no feedback.

**Fix:** Check and log errors:
```go
if err := exec.CommandContext(ctx, "wg-quick", "down", p.configPath).Run(); err != nil {
    log.Printf("warning: wg-quick down failed: %v", err)
}
if err := p.killSwitch.Disable(ctx); err != nil {
    log.Printf("warning: kill switch disable failed: %v", err)
}
```

---

### WR-07: `TunDevice.Up` ignores MTU set error

**File:** `internal/networking/tun.go:125-128`
**Issue:** When `SIOCSIFMTU` ioctl fails, the code returns `nil` (success), silently accepting an incorrect MTU:
```go
if errno != 0 {
    return nil  // <-- Should return the error
}
```

**Fix:**
```go
if errno != 0 {
    return fmt.Errorf("ioctl SIOCSIFMTU: %v", errno)
}
```

---

### WR-08: `RemoveKillSwitch` ignores all errors

**File:** `internal/platform/linux/firewall.go:203-211`
**Issue:** Both the nftables flush command error and `f.Flush(ctx)` error are discarded. The function always returns nil.

**Fix:** Check and propagate errors:
```go
func (f *Firewall) RemoveKillSwitch(ctx context.Context) error {
    switch f.mode {
    case FirewallNftables:
        return exec.CommandContext(ctx, "nft", "flush", "rule", "inet", "ghoststack", "output").Run()
    default:
        return f.Flush(ctx)
    }
}
```

---

### WR-09: `diagnose.go` passes `nil` as context.Context to `Verify`

**File:** `internal/cli/diagnose.go:123`
**Issue:** `security.NewSecureBoot("").Verify(nil, exe)` passes nil for the required context parameter. While the current implementation doesn't use the context, this is fragile — any future use of the context (timeouts, cancellation) will panic.

**Fix:** Use `context.Background()`:
```go
failures, err := security.NewSecureBoot("").Verify(context.Background(), exe)
```

---

### WR-10: `ProcessManager` has unused `stdout`/`stderr` fields

**File:** `internal/providers/process.go:30-32`
**Issue:** The `stdout` and `stderr` fields are stored in `ProcessConfig.Start()` but are never read or exposed. They use memory but have no effect.

**Fix:** Either remove them or add accessor methods if they're needed externally.

---

### WR-11: `handleConfig` POST silently ignores config body

**File:** `internal/api/server.go:130-136`
**Issue:** The POST handler reads a JSON body into a map, does nothing with it, and returns `{"status":"config_updated"}`. The config is never applied. Additionally, there's no authentication or input validation on this endpoint.

**Fix:** Either implement actual config updates with validation, or remove the POST handler to avoid misleading clients.

---

### WR-12: `Engine.StopAll` only returns last error, loses previous errors

**File:** `internal/providers/engine.go:71-83`
**Issue:** When stopping multiple providers, if several fail, only the last error is returned. Previous failures are silently lost.

**Fix:** Collect all errors using `errors.Join` (Go 1.20+) or multierror:
```go
var errs []error
for name, provider := range e.active {
    if err := provider.Stop(ctx); err != nil {
        errs = append(errs, fmt.Errorf("stop provider %s: %w", name, err))
    }
    delete(e.active, name)
}
return errors.Join(errs...)
```

---

### WR-13: `SecureBoot.RotateSecret` computes hashes but never uses SecretRotationRecord

**File:** `internal/security/secureboot.go:84-85`
**Issue:** `prevHash` and `newHash` are computed with SHA-256 but the `SecretRotationRecord` struct (defined at lines 66-71) is never instantiated. The hashes are used nowhere — not in the audit log, not stored, not returned.

**Fix:** Either use the hashes meaningfully or remove the computation.

---

### WR-14: `start.go` dead code reads config file twice with ignored errors

**File:** `internal/cli/start.go:44-46`
**Issue:** After `config.Load(configPath)` at line 35, lines 44-46 read the same file again with `os.ReadFile`, parse it with `config.LoadFromString`, and assign the result to `_`. Errors are ignored. This is dead code that performs an unnecessary syscall.

**Fix:** Remove lines 44-46.

---

### WR-15: `start.go` passes canceled context to engine on error path

**File:** `internal/cli/start.go:77`
**Issue:** On API server start failure, `engine.StopAll(startCtx)` is called. But `startCtx` was created with a 30-second timeout at line 56, and `defer startCancel()` (line 57) doesn't run until function return. This specific issue is minor (startCtx is still alive), but if `startCtx` were cancelled earlier, StopAll would receive a cancelled context and potentially fail to stop providers.

**Fix:** Use a fresh short-lived context for cleanup:
```go
stopCtx, stopCancel := context.WithTimeout(context.Background(), 5*time.Second)
defer stopCancel()
engine.StopAll(stopCtx)
```

---

## Info

### IN-01: Unused `NewNoopSandbox` variable

**File:** `internal/security/sandbox.go:229-231`
**Issue:** `NewNoopSandbox` is defined but never referenced in any provided file. If this is intended for testing, consider moving it to a test file.

---

### IN-02: Unused `pid` variable from unnecessary `Gettid()` syscall

**File:** `internal/cli/diagnose.go:107-111`
**Issue:** `syscall.Gettid()` is called and stored in `pid`, but `pid` is only used in `_ = pid` on line 111. The syscall is wasted. The seccomp check via `PR_GET_SECCOMP` doesn't need the thread ID.

**Fix:** Remove the Gettid call and the `_ = pid` line.

---

### IN-03: Unused `emptyDir()` function

**File:** `internal/plugins/ipc.go:191-193`
**Issue:** `emptyDir()` always returns an empty string and is never called anywhere in the provided files.

**Fix:** Remove the dead function.

---

### IN-04: `SecretRotationRecord` struct defined but never instantiated

**File:** `internal/security/secureboot.go:66-71`
**Issue:** The `SecretRotationRecord` type is defined with proper fields but is never created anywhere in the provided files. The `RotateSecret` method computes `prevHash` and `newHash` but stores nothing in this struct.

**Fix:** Either use the struct in `RotateSecret` or remove it.

---

### IN-05: Magic numbers in seccomp BPF filter construction

**File:** `internal/security/sandbox.go:140-199`
**Issue:** The seccomp BPF filter uses hardcoded magic values (e.g., `0x20`, `0x15`, `0x06`, `0x7FFF0000`, `0x00050000`) without named constants. These correspond to BPF instruction codes (BPF_LD|BPF_W|BPF_ABS, BPF_JMP|BPF_JEQ|BPF_K, BPF_RET|BPF_K, SECCOMP_RET_ALLOW, SECCOMP_RET_KILL_PROCESS).

**Fix:** Use const definitions from `golang.org/x/sys/unix` or well-named constants:
```go
const (
    bpfLdWAbs = 0x20
    bpfJeqK   = 0x15
    bpfRetK   = 0x06
    secRetAllow = 0x7FFF0000
    secRetKill  = 0x00050000
)
```

---

_Reviewed: 2026-07-30T12:00:00Z_
_Reviewer: gsd-code-reviewer_
_Depth: deep_
