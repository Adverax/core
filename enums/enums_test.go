package enums

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type TransactionType string

const (
	Increment TransactionType = "increment"
	Decrement                 = "decrement"
	Reset                     = "reset"
)

var CumulativeTransactionTypes = NewEnum[TransactionType](
	Increment,
	Decrement,
	Reset,
)

func BenchmarkEnum_GetAll(b *testing.B) {
	for i := 0; i < b.N; i++ {
		CumulativeTransactionTypes.GetAll()
	}
}

func TestEnum_Of_WithMemberShouldBeValid(t *testing.T) {
	v, err := CumulativeTransactionTypes.Of("increment")
	assert.NoError(t, err)
	assert.Equal(t, Increment, v)
}

func TestEnum_Of_WithNonMemberShouldBeError(t *testing.T) {
	_, err := CumulativeTransactionTypes.Of("unknown")
	assert.Error(t, err)
}
