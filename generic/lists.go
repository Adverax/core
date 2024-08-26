package generic

import "sort"

// This is package for process sorted lists

type List[T Ordered] []T

func (that List[T]) Len() int {
	return len(that)
}

func (that List[T]) Less(i, j int) bool {
	return that[i] < that[j]
}

func (that List[T]) Swap(i, j int) {
	that[i], that[j] = that[j], that[i]
}

func (that List[T]) Clone() List[T] {
	res := make(List[T], len(that))
	copy(res, that)
	return res
}

// Check for include item
func (that List[T]) Contains(item T) bool {
	l := len(that)
	if l == 0 {
		return false
	}

	i := sort.Search(l, func(i int) bool { return that[i] >= item })
	if i == l {
		return false
	}

	return that[i] == item
}

// Append id to the list
func (that *List[T]) Include(item T) {
	l := len(*that)
	if l == 0 {
		*that = List[T]{item}
		return
	}

	i := sort.Search(l, func(i int) bool { return (*that)[i] >= item })
	if i == l {
		*that = append(*that, item)
		return
	}

	if (*that)[i] == item {
		return
	}

	*that = append(*that, item)
	copy((*that)[i+1:], (*that)[i:])
	(*that)[i] = item
}

// Remove id from the list
func (that *List[T]) Exclude(item T) {
	l := len(*that)
	if l == 0 {
		return
	}

	i := sort.Search(l, func(i int) bool { return (*that)[i] >= item })
	if i == l {
		return
	}

	if (*that)[i] != item {
		return
	}

	if l == 1 {
		*that = make(List[T], 0)
		return
	}

	*that = append((*that)[:i], (*that)[i+1:]...)
}

// Merge two lists
func (that List[T]) Add(bs List[T]) List[T] {
	la := len(that)
	lb := len(bs)
	if la == 0 {
		return bs
	}
	if lb == 0 {
		return that
	}

	a := 0
	b := 0
	c := make(List[T], 0, la+lb)
	for a < la && b < lb {
		if that[a] < bs[b] {
			c = append(c, that[a])
			a++
			continue
		}
		if that[a] > bs[b] {
			c = append(c, bs[b])
			b++
			continue
		}
		c = append(c, that[a])
		a++
		b++
	}
	if a < la {
		c = append(c, that[a:]...)
	}
	if b < lb {
		c = append(c, bs[b:]...)
	}
	return c
}

// Subtract list b from list a
func (that List[T]) Sub(bs List[T]) List[T] {
	la := len(that)
	lb := len(bs)
	if la == 0 {
		return nil
	}
	if lb == 0 {
		return that
	}

	a := 0
	b := 0
	c := make(List[T], 0, la)
	for a < la && b < lb {
		if that[a] < bs[b] {
			c = append(c, that[a])
			a++
			continue
		}
		if that[a] > bs[b] {
			b++
			continue
		}
		a++
		b++
	}
	if a < la {
		c = append(c, that[a:]...)
	}
	return c
}
