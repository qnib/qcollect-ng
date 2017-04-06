package qutils

/*
import (
	"testing"
	"github.com/stretchr/testify/assert"

	"github.com/docker/docker/api/types"
)

func TestDiffCPUUsage(t *testing.T) {
	pre := types.CPUUsage{
		TotalUsage:        1000,
		UsageInKernelmode: 200,
		UsageInUsermode:   800,
		PercpuUsage: []uint64{
			600,
			200,
		},
	}
	cur := types.CPUUsage{
		TotalUsage:        1150,
		UsageInKernelmode: 300,
		UsageInUsermode:   850,
		PercpuUsage: []uint64{
			650,
			300,
		},
	}
	diff := types.CPUUsage{
		TotalUsage:        150,
		UsageInKernelmode: 100,
		UsageInUsermode:   50,
		PercpuUsage: []uint64{
			50,
			100,
		},
	}
	got := DiffCPUUsage(pre, cur)
	assert.Equal(t, diff, got)
}

func TestDiffThrottlingData(t *testing.T) {
	pre := types.ThrottlingData{
		Periods:          100,
		ThrottledPeriods: 100,
		ThrottledTime:    100,
	}
	cur := types.ThrottlingData{
		Periods:          130,
		ThrottledPeriods: 120,
		ThrottledTime:    110,
	}
	diff := types.ThrottlingData{
		Periods:          30,
		ThrottledPeriods: 20,
		ThrottledTime:    10,
	}
	got := DiffThrottlingData(pre, cur)
	assert.Equal(t, diff, got)
}

func TestDiffCPUStats(t *testing.T) {
	pre := types.CPUStats{
		CPUUsage: types.CPUUsage{
			TotalUsage:        1000,
			UsageInKernelmode: 200,
			UsageInUsermode:   800,
			PercpuUsage: []uint64{
				600,
				200,
			},
		},
		SystemUsage: 2000,
		ThrottlingData: types.ThrottlingData{
			Periods:          30,
			ThrottledPeriods: 20,
			ThrottledTime:    10,
		},
	}
	cur := types.CPUStats{
		CPUUsage: types.CPUUsage{
			TotalUsage:        1200,
			UsageInKernelmode: 300,
			UsageInUsermode:   900,
			PercpuUsage: []uint64{
				700,
				300,
			},
		},
		SystemUsage: 2500,
		ThrottlingData: types.ThrottlingData{
			Periods:          130,
			ThrottledPeriods: 120,
			ThrottledTime:    110,
		},
	}
	diff := types.CPUStats{
		CPUUsage: types.CPUUsage{
			TotalUsage:        200,
			UsageInKernelmode: 100,
			UsageInUsermode:   100,
			PercpuUsage: []uint64{
				100,
				100,
			},
		},
		SystemUsage: 500,
		ThrottlingData: types.ThrottlingData{
			Periods:          100,
			ThrottledPeriods: 100,
			ThrottledTime:    100,
		},
	}
	got := DiffCPUStats(pre, cur)
	assert.Equal(t, diff, got)
}
*/