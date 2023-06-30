package enums

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type CumulativeTransactionType string

const (
	Increment CumulativeTransactionType = "increment"
	Decrement                           = "decrement"
	Reset                               = "reset"
)

var CumulativeTransactionsTypes = NewEnum[CumulativeTransactionType](
	Increment,
	Decrement,
	Reset,
)

func BenchmarkEnum_GetAll(b *testing.B) {
	for i := 0; i < b.N; i++ {
		CumulativeTransactionsTypes.GetAll()
	}
}

func TestEnum_Of_WithMemberShouldBeValid(t *testing.T) {
	v, err := CumulativeTransactionsTypes.Of("increment")
	assert.NoError(t, err)
	assert.Equal(t, Increment, v)
}

func TestEnum_Of_WithNonMemberShouldBeError(t *testing.T) {
	_, err := CumulativeTransactionsTypes.Of("unknown")
	assert.Error(t, err)
}
