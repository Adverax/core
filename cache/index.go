package cache

type index[K comparable, V any] interface {
	truncate(iterator func(item *Item[K, V]) bool)
	assert(item *Item[K, V])
	retract(item *Item[K, V])
	flush()
}

func newIndex[K comparable, V any](features features) index[K, V] {
	if features.expiration {
		return new(indexExpiration[K, V])
	}
	if features.capacity || features.size {
		return new(indexSerial[K, V])
	}
	if features.prolongation {
		return new(indexExpiration[K, V])
	}
	return new(indexDummy[K, V])
}
