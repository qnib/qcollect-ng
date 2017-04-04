package main


import (
	"C"
	"log"
	"github.com/zpatrick/go-config"
	"github.com/qnib/qcollect-ng/types"
	fTypes "github.com/qnib/qframe/types"
 	"runtime"
	"time"
)


func extractMetrics(stats runtime.MemStats) ([]qtypes.Metric) {
	var dims map[string]string
	now := time.Now()
	metrics := []qtypes.Metric{
		qtypes.NewExt("collector", "internal", "goroutine.count", qtypes.Counter, float64(runtime.NumGoroutine()), dims, now, false),
		qtypes.NewExt("collector", "internal", "memory.alloc.total", qtypes.Counter, float64(stats.TotalAlloc), dims, now, false),
		qtypes.NewExt("collector", "internal", "memory.lookups", qtypes.Counter, float64(stats.Lookups), dims, now, false),
		qtypes.NewExt("collector", "internal", "memory.mallocs", qtypes.Counter, float64(stats.Mallocs), dims, now, false),
		qtypes.NewExt("collector", "internal", "memory.frees", qtypes.Counter, float64(stats.Frees), dims, now, false),
		qtypes.NewExt("collector", "internal", "memory.pause.total.ns", qtypes.Counter, float64(stats.PauseTotalNs), dims, now, false),
		qtypes.NewExt("collector", "internal", "memory.gc.count", qtypes.Counter, float64(stats.NumGC), dims, now, false),
		qtypes.NewExt("collector", "internal", "memory.alloc.bytes", qtypes.Gauge, float64(stats.Alloc), dims, now, false),
		qtypes.NewExt("collector", "internal", "memory.sys.bytes", qtypes.Gauge, float64(stats.Sys), dims, now, false),
		qtypes.NewExt("collector", "internal", "memory.heap.alloc.bytes", qtypes.Gauge, float64(stats.HeapAlloc), dims, now, false),
		qtypes.NewExt("collector", "internal", "memory.heap.sys.bytes", qtypes.Gauge, float64(stats.HeapSys), dims, now, false),
		qtypes.NewExt("collector", "internal", "memory.heap.idle.bytes", qtypes.Gauge, float64(stats.HeapIdle), dims, now, false),
		qtypes.NewExt("collector", "internal", "memory.heap.inuse.bytes", qtypes.Gauge, float64(stats.HeapInuse), dims, now, false),
		qtypes.NewExt("collector", "internal", "memory.heap.objects.count", qtypes.Gauge, float64(stats.HeapObjects), dims, now, false),

	}
	/*
	gauges := map[string]float64{
		"StackInuse":   float64(m.StackInuse),
		"StackSys":     float64(m.StackSys),
		"MSpanInuse":   float64(m.MSpanInuse),
		"MSpanSys":     float64(m.MSpanSys),
		"MCacheInuse":  float64(m.MCacheInuse),
		"MCacheSys":    float64(m.MCacheSys),
		"BuckHashSys":  float64(m.BuckHashSys),
		"GCSys":        float64(m.GCSys),
		"OtherSys":     float64(m.OtherSys),
		"NextGC":       float64(m.NextGC),
		"LastGC":       float64(m.LastGC),
	}*/

	return metrics
}

func Run(qChan fTypes.QChan, cfg config.Config) {
	log.Println("[II] Start internal collector")

	tick, _ := cfg.IntOr("collector.internal.tick", 1000)
	ticker := time.NewTicker(time.Millisecond * time.Duration(tick)).C

	stats := new(runtime.MemStats)

	for {
		<-ticker
		runtime.ReadMemStats(stats)
		for _, m := range extractMetrics(*stats) {
			go qChan.Data.Send(m)
		}
	}
}
