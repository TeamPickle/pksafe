package pksafe

type queue[T any] struct{ items chan T }

type Queue[T any] interface{}

func NewQueue[T any]() Queue[T]                   { return &queue[T]{items: make(chan T)} }
func NewQueueWithLimit[T any](limit int) Queue[T] { return &queue[T]{items: make(chan T, limit)} }

func (q *queue[T]) Append(item T) { q.items <- item }
func (q *queue[T]) Pop() T        { return <-q.items }
