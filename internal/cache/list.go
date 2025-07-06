package cache

type Node[T any] struct {
	data T
	prev *Node[T]
	next *Node[T]
}

func NewNode[T any](data T) *Node[T] {
	return &Node[T]{
		data: data,
		prev: nil,
		next: nil,
	}
}

type DoubleLinkedList[T any] struct {
	front *Node[T]
	back  *Node[T]
	size  int
}

func NewDoubleLinkedList[T any]() *DoubleLinkedList[T] {
	return &DoubleLinkedList[T]{
		front: nil,
		back:  nil,
		size:  0,
	}
}

func (l *DoubleLinkedList[T]) Size() int {
	return l.size
}

func (l *DoubleLinkedList[T]) PushFront(node *Node[T]) {
	if l.front == nil {
		l.front = node
		l.back = node
	} else {
		node.next = l.front
		node.prev = nil
		l.front.prev = node
		l.front = node
	}
	l.size++
}

func (l *DoubleLinkedList[T]) PopFront() *Node[T] {
	if l.size == 0 {
		return nil
	}
	ret := l.front

	if l.front == l.back {
		l.front = nil
	} else {
		l.front = l.front.next
		l.front.prev = nil
	}
	l.size--
	return ret
}

func (l *DoubleLinkedList[T]) PushBack(node *Node[T]) {
	if l.front == nil {
		l.front = node
		l.back = node
	} else {
		l.back.next = node
		node.prev = l.back
		node.next = nil
		l.back = node
	}
	l.size++
}

func (l *DoubleLinkedList[T]) PopBack() *Node[T] {
	if l.size == 0 {
		return nil
	}

	ret := l.back
	if l.front == l.back {
		l.front = nil
		l.back = nil
	} else {
		l.back = l.back.prev
		l.back.next = nil
	}
	l.size--
	return ret
}

func (l *DoubleLinkedList[T]) remove(node *Node[T]) {
	if l.size == 0 {
		return
	}

	if node.prev != nil {
		node.prev.next = node.next
	} else {
		l.front = node.next
	}

	if node.next != nil {
		node.next.prev = node.prev
	} else {
		l.back = node.prev
	}

	node.prev = nil
	node.next = nil
	l.size--
}

func (l *DoubleLinkedList[T]) MoveToFront(node *Node[T]) {
	if node == l.front {
		return
	}

	l.remove(node)
	l.PushFront(node)
}
