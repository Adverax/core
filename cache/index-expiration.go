package cache

import "sort"

type indexExpiration[K comparable, V any] struct {
	items []*Item[K, V]
}

func (that *indexExpiration[K, V]) Len() int {
	return len(that.items)
}

func (that *indexExpiration[K, V]) Swap(i, j int) {
	that.items[i], that.items[j] = that.items[j], that.items[i]
}

func (that *indexExpiration[K, V]) Less(i, j int) bool {
	return that.items[i].expiration < that.items[j].expiration
}

func (that *indexExpiration[K, V]) flush() {}

func (that *indexExpiration[K, V]) truncate(iterator func(item *Item[K, V]) bool) {
	for len(that.items) > 0 {
		item := that.items[0]
		if iterator(item) {
			that.items = that.items[1:]
			continue
		}
		return
	}
}

func (that *indexExpiration[K, V]) indexOf(item *Item[K, V]) int {
	low := 0
	high := len(that.items) - 1

	for low <= high {
		mid := low + (high-low)/2

		if that.items[mid].expiration < item.expiration {
			low = mid + 1
		} else if that.items[mid].expiration > item.expiration {
			high = mid - 1
		} else {
			if that.items[mid] == item {
				return mid
			}

			right := mid + 1
			for right <= high && that.items[right].expiration == item.expiration {
				if that.items[right] == item {
					return right
				}
				right++
			}

			left := mid - 1
			for left >= low && that.items[left].expiration == item.expiration {
				if that.items[left] == item {
					return left
				}
				left--
			}

			return -1
		}
	}

	return -1
}

func (that *indexExpiration[K, V]) retract(item *Item[K, V]) {
	index := that.indexOf(item)
	if index != -1 {
		that.items = append(that.items[:index], that.items[index+1:]...)
	}
}

func (that *indexExpiration[K, V]) assert(item *Item[K, V]) {
	index := sort.Search(
		len(that.items),
		func(i int) bool {
			return that.items[i].expiration >= item.expiration
		},
	)

	if index == len(that.items) {
		that.items = append(that.items, item)
	} else {
		that.items = append(that.items, nil)
		copy(that.items[index+1:], that.items[index:])
		that.items[index] = item
	}
}
