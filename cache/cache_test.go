package cache

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type CacheShould struct {
	suite.Suite
}

func TestCache(t *testing.T) {
	suite.Run(t, new(CacheShould))
}

func (that *CacheShould) TestPersistent() {
	c := New[string, string]()
	c.Set("hello", "word")
	item := c.Get("hello")
	assert.NotNil(that.T(), item)
	assert.Equal(that.T(), "word", item.val)
}

func (that *CacheShould) TestWithNotExpired_MustBeAccessable() {
	c := New[string, string](
		WithExpiration[string, string](time.Hour),
	)
	c.Set("hello1", "word1")
	time.Sleep(time.Microsecond)
	c.Set("hello2", "word2")
	time.Sleep(time.Microsecond)
	c.Set("hello3", "word3")
	time.Sleep(1000 * time.Millisecond)
	item := c.Get("hello1")
	assert.NotNil(that.T(), item)
	item = c.Get("hello2")
	assert.NotNil(that.T(), item)
	item = c.Get("hello3")
	assert.NotNil(that.T(), item)
}

func (that *CacheShould) TestWithExpired_MustBeNotAccessable() {
	c := New[string, string](
		WithExpiration[string, string](10 * time.Millisecond),
	)
	c.Set("hello1", "word1")
	time.Sleep(time.Microsecond)
	c.Set("hello2", "word2")
	time.Sleep(time.Microsecond)
	c.Set("hello3", "word3")
	time.Sleep(time.Microsecond)
	c.Assign("hello4", "word4", time.Hour)
	time.Sleep(20 * time.Millisecond)
	item := c.Get("hello1")
	assert.Nil(that.T(), item)
	item = c.Get("hello2")
	assert.Nil(that.T(), item)
	item = c.Get("hello3")
	assert.Nil(that.T(), item)
	item = c.Get("hello4")
	assert.NotNil(that.T(), item)
}

func (that *CacheShould) TestWithProlongation_MustBeAccessable() {
	c := New[string, string](
		WithExpiration[string, string](50*time.Millisecond),
		WithProlongation[string, string](),
	)
	c.Set("hello1", "word1")

	for i := 0; i < 10; i++ {
		c.Get("hello1")
		time.Sleep(20 * time.Millisecond)
	}

	item := c.Get("hello1")
	assert.NotNil(that.T(), item)
	expiration := time.Unix(0, item.expiration)
	d := expiration.Sub(time.Now())
	assert.True(that.T(), d > time.Millisecond*20)
}

func (that *CacheShould) TestWithCapacity_MustRemoveItems() {
	c := New[string, string](
		WithExpiration[string, string](time.Hour),
		WithCapacity[string, string](2),
	)
	c.Set("hello1", "word1")
	time.Sleep(10 * time.Millisecond)
	c.Set("hello2", "word2")
	time.Sleep(10 * time.Millisecond)
	c.Set("hello3", "word3")
	assert.Nil(that.T(), c.Get("hello1"))
	assert.NotNil(that.T(), c.Get("hello2"))
	assert.NotNil(that.T(), c.Get("hello3"))
}

func (that *CacheShould) TestWithMaxSize_MustRemoveItems() {
	c := New[string, string](
		WithExpiration[string, string](time.Hour),
		WithSize[string, string](10, func(item *Item[string, string]) int64 {
			return 4
		}),
	)
	for i := 0; i < 100; i++ {
		key := fmt.Sprintf("%d", i)
		c.Set(key, key)
	}
	assert.Equal(that.T(), 2, c.ItemCount())
}
