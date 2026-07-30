package types

type Status string

const (
	StatusUp      Status = "up"
	StatusDown    Status = "down"
	StatusDegraded Status = "degraded"
)
