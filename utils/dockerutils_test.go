package qutils

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"fmt"
)

func TestSanatizeContainerName(t *testing.T) {
	names := []string{"/crazy_lizard"}
	exp := "crazy_lizard"
	got := SanatizeContainerName(names)
	assert.Equal(t, exp, got)
}

func TestContainerNameExtractService(t *testing.T) {
	names := []string{"/influxdb_backend.1.f2hypsmktounqx1p1p85o81f5"}
	exp := DockerServiceTask{
		Name: "influxdb_backend",
		Slot: 1,
		TaskID: "f2hypsmktounqx1p1p85o81f5",
	}
	got, _ := ContainerNameExtractService(names)
	assert.Equal(t, exp, got)
	names = []string{"/crazy_lizard"}
	_, eGot := ContainerNameExtractService(names)
	fmt.Println(eGot)
	assert.NotNil(t, eGot)
}