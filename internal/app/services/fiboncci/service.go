package fiboncci

type Service interface {
	Calculate(x uint32, y uint32) ([]string, error)
}