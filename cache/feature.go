package cache

import "time"

type feature[K comparable, V any] interface {
	flush(c *Cache[K, V])
	get(c *Cache[K, V], item *Item[K, V])
	set(c *Cache[K, V], oldItem, newItem *Item[K, V])
	cleanup(c *Cache[K, V])
	close()
}

type baseFeature[K comparable, V any] struct {
	index index[K, V]
}

func (that *baseFeature[K, V]) close() {}

func (that *baseFeature[K, V]) flush(c *Cache[K, V]) {
	that.index.flush()
}

func (that *baseFeature[K, V]) get(c *Cache[K, V], item *Item[K, V]) {}

func (that *baseFeature[K, V]) set(c *Cache[K, V], oldItem, newItem *Item[K, V]) {
	if oldItem != nil {
		that.index.retract(oldItem)
	}
	if newItem != nil {
		that.index.assert(newItem)
	}
}

func (that *baseFeature[K, V]) cleanup(c *Cache[K, V]) {}

type featureExpiration[K comparable, V any] struct {
	feature[K, V]
}

func (that *featureExpiration[K, V]) cleanup(c *Cache[K, V]) {
	that.feature.cleanup(c)
	now := time.Now().UnixNano()
	c.options.index.truncate(
		func(item *Item[K, V]) bool {
			if !item.expired(now) {
				return false
			}
			delete(c.items, item.key)
			return true
		},
	)
}

type featureExpirationProlongation[K comparable, V any] struct {
	feature[K, V]
	prolongation time.Duration
	index        index[K, V]
}

func (that *featureExpirationProlongation[K, V]) get(c *Cache[K, V], item *Item[K, V]) {
	that.feature.get(c, item)
	that.index.retract(item)
	item.expiration = time.Now().Add(that.prolongation).UnixNano()
	that.index.assert(item)
}

type featureCapacity[K comparable, V any] struct {
	feature[K, V]
	index    index[K, V]
	capacity int
}

func (that *featureCapacity[K, V]) cleanup(c *Cache[K, V]) {
	that.feature.cleanup(c)
	that.index.truncate(func(item *Item[K, V]) bool {
		if len(c.items) <= that.capacity {
			return false
		}
		delete(c.items, item.key)
		return true
	})
}

type featureCapacityProlongation[K comparable, V any] struct {
	feature[K, V]
	index index[K, V]
}

func (that *featureCapacityProlongation[K, V]) get(c *Cache[K, V], item *Item[K, V]) {
	that.feature.get(c, item)
	that.index.retract(item)
	that.index.assert(item)
}

type featureSize[K comparable, V any] struct {
	feature[K, V]
	index   index[K, V]
	size    int64
	maxSize int64
	sizeOf  func(item *Item[K, V]) int64
}

func (that *featureSize[K, V]) cleanup(c *Cache[K, V]) {
	that.feature.cleanup(c)
	that.index.truncate(
		func(item *Item[K, V]) bool {
			if that.size <= that.maxSize {
				return false
			}
			that.size -= item.size
			delete(c.items, item.key)
			return true
		},
	)
}

func (that *featureSize[K, V]) set(c *Cache[K, V], oldItem, newItem *Item[K, V]) {
	if oldItem != nil {
		that.size -= oldItem.size
	}
	if newItem != nil {
		newItem.size = that.sizeOf(newItem)
		that.size += newItem.size
	}
	that.feature.set(c, oldItem, newItem)
}

type featureDaemon[K comparable, V any] struct {
	feature[K, V]
	interval time.Duration
	done     chan struct{}
}

func (that *featureDaemon[K, V]) cleanup(c *Cache[K, V]) {
	// nothing
}

func (that *featureDaemon[K, V]) start(c *Cache[K, V]) {
	go func() {
		ticker := time.NewTicker(that.interval)
		for {
			select {
			case <-that.done:
				return
			case <-ticker.C:
				that.autoCleanup(c)
			}
		}
	}()
}

func (that *featureDaemon[K, V]) autoCleanup(c *Cache[K, V]) {
	c.mu.Lock()
	defer c.mu.Unlock()

	that.feature.cleanup(c)
}

func (that *featureDaemon[K, V]) Close() {
	close(that.done)
}
