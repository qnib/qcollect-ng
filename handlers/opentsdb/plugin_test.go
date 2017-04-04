package main

import (
	"testing"
	"github.com/stretchr/testify/assert"

	"github.com/qnib/qcollect-ng/types"

	"time"
)

var (
	m = qtypes.NewExt("input", "file", "test", qtypes.Gauge, 10.0, make(map[string]string), time.Unix(1491227867, 0), false)
)

func TestConvertToOpenTSDBHandler(t *testing.T) {
	exp := "put test 1491227867 10.000000\n"
	assert.Equal(t, exp, m.ConvertToOpenTSDBHandler(""))
}
