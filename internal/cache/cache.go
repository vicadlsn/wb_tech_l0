package cache

import "sync"

// CacheEntry представляет пару ключ-значение, хранимую в кэше.
type CacheEntry[K comparable, V any] struct {
	key   K
	value V
}

// LRUCache реализует кэш с вытеснением по принципу Least Recently Used (LRU) через двусвязный список.
// При достижении максимальной вместимости кэша самый давно неиспользуемый элемент удаляется.
type LRUCache[K comparable, V any] struct {
	capacity int
	data     map[K]*Node[CacheEntry[K, V]]
	list     *DoubleLinkedList[CacheEntry[K, V]]
	mutex    sync.Mutex
}

// NewLRUCache создает новый LRU-кэш с заданной вместимостью.
func NewLRUCache[K comparable, V any](capacity int) *LRUCache[K, V] {
	return &LRUCache[K, V]{
		capacity: capacity,
		data:     make(map[K]*Node[CacheEntry[K, V]]),
		list:     NewDoubleLinkedList[CacheEntry[K, V]](),
	}
}

// Get возвращает значение по ключу и булево значение, указывающее на успешность поиска.
// Если элемент найден, он перемещается в начало списка.
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

// Put добавляет элемент в кэш или обновляет существующий.
// Если кэш заполнен, удаляется наименее недавно использованный элемент (в конце списка).
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
