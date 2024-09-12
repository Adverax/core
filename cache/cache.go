package cache

import (
	"core"
	"sync"
	"time"
)

type Item[K comparable, V any] struct {
	id         int64
	key        K
	val        V
	expiration int64
	size       int64
}

func (that *Item[K, V]) Key() K {
	return that.key
}

func (that *Item[K, V]) Value() V {
	return that.val
}

func (that *Item[K, V]) Expiration() int64 {
	return that.expiration
}

func (that *Item[K, V]) Expired() bool {
	return that.expired(time.Now().UnixNano())
}

func (that *Item[K, V]) expired(now int64) bool {
	return now > that.expiration
}

type Cache[K comparable, V any] struct {
	mu      sync.Mutex
	items   map[K]*Item[K, V]
	options *Options[K, V]
	counter int64
}

func (that *Cache[K, V]) Close() {
	that.options.close()
}

func (that *Cache[K, V]) Set(k K, v V) {
	that.Assign(k, v, that.options.expiration)
}

func (that *Cache[K, V]) Assign(k K, v V, d time.Duration) {
	that.mu.Lock()
	defer that.mu.Unlock()

	that.options.cleanup(that)

	oldItem := that.get(k)
	newItem := that.set(k, v, d)
	that.options.set(that, oldItem, newItem)
}

func (that *Cache[K, V]) set(k K, v V, d time.Duration) *Item[K, V] {
	that.counter++
	item := &Item[K, V]{
		id:         that.counter,
		key:        k,
		val:        v,
		expiration: time.Now().Add(d).UnixNano(),
	}
	that.items[k] = item
	return item
}

func (that *Cache[K, V]) Add(k K, v V) error {
	return that.Append(k, v, that.options.expiration)
}

func (that *Cache[K, V]) Append(k K, v V, d time.Duration) error {
	that.mu.Lock()
	defer that.mu.Unlock()

	that.options.cleanup(that)

	item := that.get(k)
	if item != nil {
		return core.ErrDuplicate
	}

	item = that.set(k, v, d)
	that.options.set(that, nil, item)
	return nil
}

func (that *Cache[K, V]) Replace(k K, v V, d time.Duration) error {
	that.mu.Lock()
	defer that.mu.Unlock()

	that.options.cleanup(that)

	oldItem := that.get(k)
	if oldItem == nil {
		return core.ErrNoMatch
	}

	newItem := that.set(k, v, d)
	that.options.set(that, oldItem, newItem)
	return nil
}

func (that *Cache[K, V]) Get(k K) *Item[K, V] {
	that.mu.Lock()
	defer that.mu.Unlock()

	that.options.cleanup(that)

	item := that.get(k)

	if item != nil && item.Expired() {
		delete(that.items, k)
		that.options.set(that, item, nil)
		return nil
	}

	return item
}

func (that *Cache[K, V]) get(k K) *Item[K, V] {
	item, found := that.items[k]
	if !found {
		return nil
	}

	that.options.get(that, item)
	return item
}

func (that *Cache[K, T]) Delete(k K) {
	that.mu.Lock()
	defer that.mu.Unlock()

	that.options.cleanup(that)

	item := that.get(k)
	if item != nil {
		delete(that.items, item.key)
		that.options.set(that, item, nil)
	}
}

func (that *Cache[K, V]) ItemCount() int {
	that.mu.Lock()
	defer that.mu.Unlock()

	that.options.cleanup(that)

	return len(that.items)
}

func (that *Cache[K, V]) Items() map[K]*Item[K, V] {
	that.mu.Lock()
	defer that.mu.Unlock()

	that.options.cleanup(that)

	m := make(map[K]*Item[K, V], len(that.items))
	for k, item := range that.items {
		m[k] = item
	}

	return m
}

func (that *Cache[K, V]) Flush() {
	that.mu.Lock()
	defer that.mu.Unlock()

	that.options.flush(that)
	that.items = map[K]*Item[K, V]{}
}

func New[K comparable, V any](options ...Option[K, V]) *Cache[K, V] {
	opts := &Options[K, V]{
		expiration: time.Hour,
	}

	for _, option := range options {
		option(opts)
	}
	opts.features.resolved = true
	opts.index = newIndex[K, V](opts.features)
	opts.feature = &baseFeature[K, V]{
		index: opts.index,
	}
	for _, option := range options {
		option(opts)
	}

	c := &Cache[K, V]{
		items:   make(map[K]*Item[K, V]),
		options: opts,
	}

	if opts.daemonInterval != 0 {
		opts.feature = &featureDaemon[K, V]{
			feature:  opts.feature,
			interval: opts.daemonInterval,
			done:     make(chan struct{}),
		}
	}

	return c
}
