package pksafe

import (
	"sync"

	"golang.org/x/exp/constraints"
)

type alertMapValueFunc[K constraints.Ordered, T any] func() (key K, value T)

type safeMap[K constraints.Ordered, T any] struct {
	mutex     sync.RWMutex
	data      map[K]T
	alertChan chan alertMapValueFunc[K, T]
}

type Map[K constraints.Ordered, T any] interface {
	Get(key K) (value T, exist bool)
	Has(key K) (exist bool)
	Insert(key K, value T) Map[K, T]
	Alert() <-chan alertMapValueFunc[K, T]
	Size() (size int)
	Delete(key K)
	Value() (value map[K]T)
	Clear()
}

func NewMap[K constraints.Ordered, T any]() Map[K, T] {
	return &safeMap[K, T]{
		data:      make(map[K]T),
		alertChan: make(chan alertMapValueFunc[K, T])}
}

func (m *safeMap[K, T]) Get(key K) (value T, exist bool) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	value, exist = m.data[key]
	return
}

func (m *safeMap[K, T]) Has(key K) (exist bool) {
	_, exist = m.Get(key)
	return
}

func (m *safeMap[K, T]) Insert(key K, value T) Map[K, T] {
	if len(m.alertChan) > 0 {
		<-m.alertChan
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.data[key] = value
	m.alertChan <- func() (key K, value T) { return key, value }
	return m
}

func (m *safeMap[K, T]) Alert() <-chan alertMapValueFunc[K, T] { return m.alertChan }

func (m *safeMap[K, T]) Size() (size int) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	size = len(m.data)
	return
}

func (m *safeMap[K, T]) Delete(key K) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	delete(m.data, key)
}

func (m *safeMap[K, T]) Value() (value map[K]T) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	value = make(map[K]T)
	for k, v := range m.data {
		value[k] = v
	}

	return
}

func (m *safeMap[K, T]) getKeys() (result []K) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	for key := range m.data {
		result = append(result, key)
	}
	return
}

func (m *safeMap[K, T]) Clear() {
	keys := m.getKeys()

	m.mutex.Lock()
	defer m.mutex.Unlock()
	for _, key := range keys {
		delete(m.data, key)
	}
}
