package logrus

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestTemplateFormatter_Format(t *testing.T) {
	entry := &Entry{
		Logger: nil,
		Data: Fields{
			"my_key": "my_value",
		},
		Time:    time.Time{},
		Level:   2,
		Caller:  nil,
		Message: "hello",
		Buffer:  nil,
		Context: context.Background(),
	}

	f := &TemplateFormatter{
		TimestampFormat: "2006/01/02 15:04:05",
	}
	data, err := f.Format(entry)
	require.NoError(t, err)
	fmt.Println(string(data))
}
