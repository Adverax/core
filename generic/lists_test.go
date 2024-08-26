package generic

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestContains(t *testing.T) {
	type Test struct {
		items List[int64]
		item  int64
		has   bool
	}

	tests := map[string]Test{
		"find in empty list": {
			item:  1,
			items: []int64{},
			has:   false,
		},
		"find first item": {
			item:  1,
			items: []int64{1, 2, 3, 4},
			has:   true,
		},
		"find last item": {
			item:  5,
			items: []int64{1, 2, 3, 4, 5},
			has:   true,
		},
		"find middle item": {
			item:  3,
			items: []int64{1, 2, 3, 4, 5},
			has:   true,
		},
		"find illegal item": {
			item:  10,
			items: []int64{1, 2, 3, 4, 5},
			has:   false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			has := test.items.Contains(test.item)
			assert.Equal(t, test.has, has)
		})
	}
}

func TestInclude(t *testing.T) {
	type Test struct {
		src List[int64]
		dst List[int64]
		id  int64
	}

	tests := map[string]Test{
		"insert into empty list": {
			id:  1,
			dst: []int64{1},
		},
		"insert last item": {
			id:  5,
			src: []int64{1, 2, 3, 4},
			dst: []int64{1, 2, 3, 4, 5},
		},
		"insert first item": {
			id:  1,
			src: []int64{2, 3, 4, 5},
			dst: []int64{1, 2, 3, 4, 5},
		},
		"insert middle item": {
			id:  2,
			src: []int64{1, 3, 4, 5},
			dst: []int64{1, 2, 3, 4, 5},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			test.src.Include(test.id)
			assert.Equal(t, test.dst, test.src)
		})
	}
}

func TestExclude(t *testing.T) {
	type Test struct {
		src List[int64]
		dst List[int64]
		id  int64
	}

	tests := map[string]Test{
		"remove from empty list": {
			id:  1,
			dst: nil,
		},
		"remove first item": {
			id:  1,
			src: []int64{1, 2, 3, 4, 5},
			dst: []int64{2, 3, 4, 5},
		},
		"remove last item": {
			id:  5,
			src: []int64{1, 2, 3, 4, 5},
			dst: []int64{1, 2, 3, 4},
		},
		"remove item": {
			id:  2,
			src: []int64{1, 2, 3, 4, 5},
			dst: []int64{1, 3, 4, 5},
		},
		"remove illegal item": {
			id:  77,
			src: []int64{1, 2, 3, 4, 5},
			dst: []int64{1, 2, 3, 4, 5},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			test.src.Exclude(test.id)
			assert.Equal(t, test.dst, test.src)
		})
	}
}

func TestAdd(t *testing.T) {
	type Test struct {
		a  List[int64]
		b  List[int64]
		c  List[int64]
		id int64
	}

	tests := map[string]Test{
		"merge with empty b": {
			id: 1,
			a:  []int64{1, 2, 3, 4, 5},
			b:  nil,
			c:  []int64{1, 2, 3, 4, 5},
		},
		"merge with empty a": {
			id: 1,
			a:  nil,
			b:  []int64{1, 2, 3, 4, 5},
			c:  []int64{1, 2, 3, 4, 5},
		},
		"merge full": {
			id: 1,
			a:  []int64{1, 4, 5, 6, 7},
			b:  []int64{1, 2, 3, 4, 5},
			c:  []int64{1, 2, 3, 4, 5, 6, 7},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			c := test.a.Add(test.b)
			assert.Equal(t, test.c, c)
		})
	}
}

func TestSubtract(t *testing.T) {
	type Test struct {
		a  List[int64]
		b  List[int64]
		c  List[int64]
		id int64
	}

	tests := map[string]Test{
		"subtract with empty b": {
			id: 1,
			a:  []int64{1, 2, 3, 4, 5},
			b:  nil,
			c:  []int64{1, 2, 3, 4, 5},
		},
		"subtract with empty a": {
			id: 1,
			a:  nil,
			b:  []int64{1, 2, 3, 4, 5},
			c:  nil,
		},
		"subtract full": {
			id: 1,
			a:  []int64{1, 2, 3, 4, 5, 6, 7, 8, 9},
			b:  []int64{0, 2, 4, 6, 8, 10},
			c:  []int64{1, 3, 5, 7, 9},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			c := test.a.Sub(test.b)
			assert.Equal(t, test.c, c)
		})
	}
}
