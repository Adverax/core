package hyphenations

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHyphen(t *testing.T) {
	type Test struct {
		Src   string
		Dst   []string
		Width int
	}

	tests := map[string]Test{
		"1": {
			Src:   "Олфен пластир трансдерм. 140 мг/12 годин  #10 (1/10шт)*469.00грн СЕРИЯ: E0396",
			Dst:   []string{"Олфен пластир трансдерм. 140 мг/12 годин", "#10 (1/10шт)*469.00грн СЕРИЯ: E0396"},
			Width: 40,
		},
		"2": {
			Width: 24,
			Src:   "Анальгін табл. 0,5 г  #10",
			Dst:   []string{"Анальгін табл. 0,5 г", "#10"},
		},
		"3": {
			Width: 24,
			Src:   "Прокладки урологічніBELLACONTROLDISCREETmicro.(18шт)  ?20",
			Dst:   []string{"Прокладки урологічніBEL-", "LACONTROLDISCREET-", "micro.(18шт)  ?20"},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			hyphen, err := NewBuilder().Build()
			assert.NoError(t, err)
			actual := hyphen.Split(test.Src, test.Width)
			assert.Equal(t, test.Dst, actual)
		})
	}
}
