package tdx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEventLogRecoder(t *testing.T) {
	r := &TDXEventLogRecorder{}
	r.ProbeRecorder()
	l, err := r.FullEventLog()
	assert.Nil(t, err)
	assert.NotNil(t, l)
}
