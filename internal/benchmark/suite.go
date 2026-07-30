package benchmark

import (
	"context"
	"fmt"
	"testing"
	"time"
)

type BenchmarkResult struct {
	Name     string
	Duration time.Duration
	Memory   uint64
	Error    error
}

type Suite struct {
	results []BenchmarkResult
}

func NewSuite() *Suite {
	return &Suite{results: make([]BenchmarkResult, 0)}
}

func (s *Suite) Run(ctx context.Context, name string, fn func(context.Context) error) BenchmarkResult {
	start := time.Now()

	result := BenchmarkResult{Name: name}

	if err := fn(ctx); err != nil {
		result.Error = err
		result.Duration = time.Since(start)
		s.results = append(s.results, result)
		return result
	}

	result.Duration = time.Since(start)
	s.results = append(s.results, result)
	return result
}

func (s *Suite) Results() []BenchmarkResult {
	return append([]BenchmarkResult(nil), s.results...)
}

func (s *Suite) Summary() string {
	var passed, failed int
	var totalDuration time.Duration

	for _, r := range s.results {
		totalDuration += r.Duration
		if r.Error != nil {
			failed++
		} else {
			passed++
		}
	}

	return fmt.Sprintf("benchmark: %d passed, %d failed, total %s", passed, failed, totalDuration)
}

func GoBenchmark(b *testing.B, name string, fn func(b *testing.B)) {
	b.Run(name, fn)
}
