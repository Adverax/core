package cache

type indexDummy[K comparable, V any] struct{}

func (that *indexDummy[K, V]) flush() {}

func (that *indexDummy[K, V]) truncate(iterator func(item *Item[K, V]) bool) {}

func (that *indexDummy[K, V]) assert(item *Item[K, V]) {}

func (that *indexDummy[K, V]) retract(item *Item[K, V]) {}
