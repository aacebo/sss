package async

import "sync"

type Map[K comparable, V any] struct {
	mu      sync.RWMutex
	content map[K]V
}

func NewMap[K comparable, V any]() Map[K, V] {
	return Map[K, V]{
		sync.RWMutex{},
		map[K]V{},
	}
}

func (self *Map[K, V]) Has(key K) bool {
	self.mu.RLock()
	_, ok := self.content[key]
	self.mu.RUnlock()
	return ok
}

func (self *Map[K, V]) Get(key K) V {
	self.mu.RLock()
	v := self.content[key]
	self.mu.RUnlock()
	return v
}

func (self *Map[K, V]) Set(key K, value V) {
	self.mu.Lock()
	self.content[key] = value
	self.mu.Unlock()
}

func (self *Map[K, V]) Del(key K) {
	self.mu.Lock()
	delete(self.content, key)
	self.mu.Unlock()
}

func (self *Map[K, V]) Len() int {
	self.mu.RLock()
	l := len(self.content)
	self.mu.RUnlock()
	return l
}

func (self *Map[K, V]) Map() map[K]V {
	return self.content
}

func (self *Map[K, V]) Slice() []V {
	arr := []V{}

	for _, v := range self.content {
		arr = append(arr, v)
	}

	return arr
}

func (self *Map[K, V]) ForEach(callback func(k K, v V)) {
	self.mu.RLock()

	for k, v := range self.content {
		callback(k, v)
	}

	self.mu.RUnlock()
}
