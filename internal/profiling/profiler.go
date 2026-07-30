package profiling

import (
	"context"
	"errors"
	"runtime"
	"time"

	"github.com/ghoststack/ghoststack/pkg/types"
)

type Profiler struct {
	enabled bool
	output  string
}

func NewProfiler(output string) *Profiler {
	return &Profiler{output: output}
}

func (p *Profiler) Enable(ctx context.Context) error {
	if p.enabled {
		return nil
	}
	if p.output == "" {
		return errors.New("profiler output path is empty")
	}
	p.enabled = true
	return nil
}

func (p *Profiler) Disable(ctx context.Context) error {
	if !p.enabled {
		return nil
	}
	p.enabled = false
	return nil
}

func (p *Profiler) Snapshot() (types.ProfileSnapshot, error) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	return types.ProfileSnapshot{
		Timestamp:  time.Now().UnixNano(),
		AllocMB:    float64(m.Alloc) / 1024 / 1024,
		SysMB:      float64(m.Sys) / 1024 / 1024,
		Goroutines: runtime.NumGoroutine(),
	}, nil
}
