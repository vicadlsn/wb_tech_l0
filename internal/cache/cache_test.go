package cache_test

import (
	"testing"
	"webtechl0/internal/cache"
)

func TestLRUCache(t *testing.T) {
	c := cache.NewLRUCache[int, int](2)
	c.Put(1, 1)
	c.Put(2, 2)

	if val, ok := c.Get(1); !ok || val != 1 {
		t.Errorf("expected key 1 to have value 1, got %v", val)
	}

	if val, ok := c.Get(2); !ok || val != 2 {
		t.Errorf("expected key 2 to have value 2, got %v", val)
	}

	c.Get(1)
	c.Put(3, 3)
	if val, ok := c.Get(3); !ok || val != 3 {
		t.Errorf("expected key 3 to have value 3, got %v", val)
	}

	if _, ok := c.Get(2); ok {
		t.Errorf("expected key 2 to be evicted")
	}

	if val, ok := c.Get(1); !ok || val != 1 {
		t.Errorf("expected key 1 to be present with value 1")
	}

	c.Put(3, 30)
	if val, ok := c.Get(3); !ok || val != 30 {
		t.Errorf("expected key 3 to be 30")
	}

	c.Put(2, 2)
	if _, ok := c.Get(1); ok {
		t.Errorf("expected key 1 to be evicted")
	}
}
