package cache

import "sort"

type indexSerial[K comparable, V any] struct {
	items []*Item[K, V]
}

func (that *indexSerial[K, V]) Len() int {
	return len(that.items)
}

func (that *indexSerial[K, V]) Swap(i, j int) {
	that.items[i], that.items[j] = that.items[j], that.items[i]
}

func (that *indexSerial[K, V]) Less(i, j int) bool {
	return that.items[i].id < that.items[j].id
}

func (that *indexSerial[K, V]) flush() {}

func (that *indexSerial[K, V]) truncate(iterator func(item *Item[K, V]) bool) {
	for len(that.items) > 0 {
		item := that.items[0]
		if iterator(item) {
			that.items = that.items[1:]
			continue
		}
		return
	}
}

func (that *indexSerial[K, V]) indexOf(item *Item[K, V]) int {
	low := 0
	high := len(that.items) - 1

	for low <= high {
		mid := low + (high-low)/2

		if that.items[mid].id < item.id {
			low = mid + 1
		} else if that.items[mid].id > item.id {
			high = mid - 1
		} else {
			if that.items[mid] == item {
				return mid
			}

			right := mid + 1
			for right <= high && that.items[right].id == item.id {
				if that.items[right] == item {
					return right
				}
				right++
			}

			left := mid - 1
			for left >= low && that.items[left].id == item.id {
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

func (that *indexSerial[K, V]) retract(item *Item[K, V]) {
	index := that.indexOf(item)
	if index != -1 {
		that.items = append(that.items[:index], that.items[index+1:]...)
	}
}

func (that *indexSerial[K, V]) assert(item *Item[K, V]) {
	index := sort.Search(
		len(that.items),
		func(i int) bool {
			return that.items[i].id >= item.id
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
