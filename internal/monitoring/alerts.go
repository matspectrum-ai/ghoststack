package monitoring

import "context"

type Alert struct {
	ID       string
	Severity string
	Message  string
}

type AlertManager struct {
	alerts []Alert
}

func NewAlertManager() *AlertManager {
	return &AlertManager{alerts: make([]Alert, 0)}
}

func (m *AlertManager) Add(ctx context.Context, alert Alert) error {
	m.alerts = append(m.alerts, alert)
	return nil
}

func (m *AlertManager) List(ctx context.Context) []Alert {
	return append([]Alert(nil), m.alerts...)
}
