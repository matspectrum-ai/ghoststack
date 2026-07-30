package types

type ProfileSnapshot struct {
	Timestamp  int64
	AllocMB    float64
	SysMB      float64
	Goroutines int
}

type MemoryStats interface {
	AllocatedBytes() int
}

type MemoryPoolStats struct {
	Allocated int
}

func (m MemoryPoolStats) AllocatedBytes() int {
	return m.Allocated
}
