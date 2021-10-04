package cache

type Rdb interface {
	GetFibonacciSequence(x uint32, y uint32) ([]string, uint32, error)
	SetFibonacciSequence(sequence *[]string, x uint32, y uint32, done chan struct{})
	PrepareRDB() error
}
