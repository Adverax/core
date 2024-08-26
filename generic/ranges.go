package generic

type Range[T Ordered] struct {
	Min T
	Max T
}

func (r Range[T]) Contains(value T) bool {
	return r.Min <= value && value <= r.Max
}

func (r Range[T]) Overlaps(other Range[T]) bool {
	return r.Min <= other.Max && other.Min <= r.Max
}

func NewRange[T Ordered](min, max T) Range[T] {
	return Range[T]{Min: min, Max: max}
}
