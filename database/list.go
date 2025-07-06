package database

import "sync"

type List struct {
	mu  sync.RWMutex
	elements []string
}

func NewList() *List {
	return &List{
		elements: make([]string, 0),
	}
}

func (l *List) Type() string {
	return "list"
}

func (l *List) LPush(elements ...string) int {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.elements = append(elements, l.elements...)
	return len(l.elements)
}

func (l *List) RPush(elements ...string) int {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.elements = append(l.elements, elements...)
	return len(l.elements)
}

func (l *List) LPop() (string, bool) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if len(l.elements) == 0 {
		return "", false
	}
	val := l.elements[0]
	l.elements = l.elements[1:]
	return val, true
}

func (l *List) RPop() (string, bool) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if len(l.elements) == 0 {
		return "", false
	}
	val := l.elements[len(l.elements)-1]
	l.elements = l.elements[:len(l.elements)-1]
	return val, true
}

func (l *List) LLen() int {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return len(l.elements)
}