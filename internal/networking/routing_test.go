package networking

import (
	"context"
	"errors"
	"testing"
)

func TestRouteTableAddRemoveList(t *testing.T) {
	table := newRouteTable()

	list, err := table.List(context.Background())
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(list) != 0 {
		t.Fatalf("expected empty list, got %v", list)
	}

	if err := table.Add(context.Background(), "10.0.0.0/24", "10.0.0.1"); err != nil {
		t.Fatalf("add: %v", err)
	}

	if err := table.Add(context.Background(), "10.0.0.0/24", "10.0.0.1"); !errors.Is(err, ErrRouteAlreadyExists) {
		t.Fatalf("expected ErrRouteAlreadyExists, got %v", err)
	}

	list, err = table.List(context.Background())
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(list) != 1 {
		t.Fatalf("expected 1 route, got %d", len(list))
	}

	if err := table.Remove(context.Background(), "10.0.0.0/24"); err != nil {
		t.Fatalf("remove: %v", err)
	}

	if err := table.Remove(context.Background(), "10.0.0.0/24"); !errors.Is(err, ErrRouteNotFound) {
		t.Fatalf("expected ErrRouteNotFound, got %v", err)
	}
}
