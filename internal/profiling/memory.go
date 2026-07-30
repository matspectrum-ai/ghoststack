package profiling

import "github.com/ghoststack/ghoststack/pkg/types"

type memoryPool struct {
	allocated int
}

func newMemoryPool(initial int) *memoryPool {
	return &memoryPool{allocated: initial}
}

func (p *memoryPool) Stats() types.MemoryPoolStats {
	return types.MemoryPoolStats{Allocated: p.allocated}
}

func (p *memoryPool) Reset() {
	p.allocated = 0
}
