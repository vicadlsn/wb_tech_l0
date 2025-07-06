package cache

import "sync"

type CacheEntry[K comparable, V any] struct {
	key   K
	value V
}

type LRUCache[K comparable, V any] struct {
	capacity int
	data     map[K]*Node[CacheEntry[K, V]]
	list     *DoubleLinkedList[CacheEntry[K, V]]
	mutex    sync.Mutex
}

func NewLRUCache[K comparable, V any](capacity int) *LRUCache[K, V] {
	return &LRUCache[K, V]{
		capacity: capacity,
		data:     make(map[K]*Node[CacheEntry[K, V]]),
		list:     NewDoubleLinkedList[CacheEntry[K, V]](),
	}
}

func (c *LRUCache[K, V]) Get(key K) (V, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if node, ok := c.data[key]; ok {
		c.list.MoveToFront(node)
		return node.data.value, true
	}

	var v V
	return v, false
}

func (c *LRUCache[K, V]) Put(key K, value V) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.capacity == 0 {
		return
	}

	if node, ok := c.data[key]; ok {
		node.data.value = value
		c.list.MoveToFront(node)
	} else {
		node := NewNode(CacheEntry[K, V]{key: key, value: value})
		c.list.PushFront(node)
		c.data[key] = node
	}

	if c.list.Size() > c.capacity {
		node := c.list.PopBack()
		delete(c.data, node.data.key)
	}
}
