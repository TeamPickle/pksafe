package pksafe

import (
	"reflect"
	"sync"
)

type set[T any] struct {
	mutex     sync.RWMutex
	data      map[int]T
	alertChan chan T
}

type Set[T any] interface {
	Add(item T) Set[T]
	Alert() <-chan T
	Delete(item T)
	Clear()
	Has(item T) (exist bool)
	Size() (size int)
	Value() (values []T)
	ForEach(fn func(item T))
}

func NewSet[T any]() Set[T] { return &set[T]{data: make(map[int]T), alertChan: make(chan T)} }

func (s *set[T]) Add(item T) Set[T] {
	if len(s.alertChan) > 0 {
		<-s.alertChan
	}

	size := s.Size()

	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.data[size] = item
	s.alertChan <- item
	return s
}

func (s *set[T]) Alert() <-chan T { return s.alertChan }

func (s *set[T]) Delete(item T) {
	s.mutex.RLock()

	dataIndex := -1
	for index, value := range s.data {
		if reflect.DeepEqual(item, value) {
			dataIndex = index
			break
		}
	}

	s.mutex.RUnlock()
	if dataIndex == -1 {
		return
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()
	delete(s.data, dataIndex)
}

func (s *set[T]) Clear() {
	size := s.Size()

	s.mutex.Lock()
	defer s.mutex.Unlock()

	for i := 0; i < size; i++ {
		delete(s.data, i)
	}
}

func (s *set[T]) Has(item T) (exist bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	for _, value := range s.data {
		if reflect.DeepEqual(item, value) {
			exist = true
			return
		}
	}
	return
}

func (s *set[T]) Size() (size int) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	size = len(s.data)
	return
}

func (s *set[T]) Value() (values []T) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	size := len(s.data)

	for i := 0; i < size; i++ {
		if value, exist := s.data[i]; exist {
			values = append(values, value)
		}
	}

	return
}

func (s *set[T]) ForEach(fn func(item T)) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	size := len(s.data)
	for i := 0; i < size; i++ {
		if value, exist := s.data[i]; exist {
			fn(value)
		}
	}
}
