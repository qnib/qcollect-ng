package qutils



import (
	"github.com/docker/docker/api/types"
	"github.com/fsouza/go-dockerclient"

)

func TransformFsouzaToDocker(stats docker.CPUStats) types.CPUStats {
	return types.CPUStats{
		CPUUsage: types.CPUUsage{
			TotalUsage:        stats.CPUUsage.TotalUsage,
			UsageInKernelmode: stats.CPUUsage.UsageInKernelmode,
			UsageInUsermode:   stats.CPUUsage.UsageInUsermode,
			PercpuUsage: stats.CPUUsage.PercpuUsage,
		},
		SystemUsage: stats.SystemCPUUsage,
		ThrottlingData: types.ThrottlingData{
			Periods:          stats.ThrottlingData.Periods,
			ThrottledPeriods: stats.ThrottlingData.ThrottledPeriods,
			ThrottledTime:    stats.ThrottlingData.ThrottledTime,
		},
	}
}

// DiffCPUUsage create a diff out ot two (plus knowledge about the system usage)
func DiffCPUUsage(pre types.CPUUsage, cur types.CPUUsage) types.CPUUsage {
	var cpuu types.CPUUsage
	cpuu.TotalUsage = cur.TotalUsage - pre.TotalUsage
	cpuu.UsageInKernelmode = cur.UsageInKernelmode - pre.UsageInKernelmode
	cpuu.UsageInUsermode = cur.UsageInUsermode - pre.UsageInUsermode
	pCPU := cur.PercpuUsage
	for idx, c := range pre.PercpuUsage {
		pCPU[idx] = (pCPU[idx] - c)
	}
	cpuu.PercpuUsage = pCPU
	return cpuu
}

// DiffThrottlingData diffs two ThrottlingData
func DiffThrottlingData(pre types.ThrottlingData, cur types.ThrottlingData) types.ThrottlingData {
	return types.ThrottlingData{
		// Number of periods with throttling active
		Periods: cur.Periods - pre.Periods,
		// Number of periods when the container hits its throttling limit.
		ThrottledPeriods: cur.ThrottledPeriods - pre.ThrottledPeriods,
		// Aggregate time the container was throttled for in nanoseconds.
		ThrottledTime: cur.ThrottledTime - pre.ThrottledTime,
	}
}

// DiffCPUStats create a diff out of two CPUStats
func DiffCPUStats(pre types.CPUStats, cur types.CPUStats) types.CPUStats {
	var cstat types.CPUStats
	cstat.SystemUsage = cur.SystemUsage - pre.SystemUsage
	cstat.ThrottlingData = DiffThrottlingData(pre.ThrottlingData, cur.ThrottlingData)
	cstat.CPUUsage = DiffCPUUsage(pre.CPUUsage, cur.CPUUsage)

	return cstat
}