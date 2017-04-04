package main
/*
import (
	"C"
	"log"
	"regexp"

	"github.com/zpatrick/go-config"
	"github.com/qnib/qcollect-ng/types"
	"github.com/qnib/qcollect-ng/utils"
	fTypes "github.com/qnib/qframe/types"
 	"runtime"
	"strconv"
	"time"
	"fmt"
)


func extractMetrics() (metrics []qtypes.Metric) {
	counters := map[string]float64{
		"NumGoroutine": float64(runtime.NumGoroutine()),
		"TotalAlloc":   float64(m.TotalAlloc),
		"Lookups":      float64(m.Lookups),
		"Mallocs":      float64(m.Mallocs),
		"Frees":        float64(m.Frees),
		"PauseTotalNs": float64(m.PauseTotalNs),
		"NumGC":        float64(m.NumGC),
	}

	gauges := map[string]float64{
		"Alloc":        float64(m.Alloc),
		"Sys":          float64(m.Sys),
		"HeapAlloc":    float64(m.HeapAlloc),
		"HeapSys":      float64(m.HeapSys),
		"HeapIdle":     float64(m.HeapIdle),
		"HeapInuse":    float64(m.HeapInuse),
		"HeapReleased": float64(m.HeapReleased),
		"HeapObjects":  float64(m.HeapObjects),
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
	}

	return
}

func Run(qChan fTypes.QChan, cfg config.Config) {
	log.Println("[II] Start internal collector")

	stats := new(runtime.MemStats)
	runtime.ReadMemStats(stats)

	for {
		//qm := qtypes.NewExt("input", "file", m["metric"], qtypes.Gauge, val, dim, time.Unix(int64(t),0), false)
		//qChan.Data.Send(qm)
		fmt.Printf("huhu\n")
	}
}
*/
