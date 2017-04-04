package main

import (
	"C"
	"log"
	"regexp"

	"github.com/zpatrick/go-config"
	"github.com/qnib/qcollect-ng/types"
	"github.com/qnib/qcollect-ng/utils"
	fTypes "github.com/qnib/qframe/types"
	"github.com/hpcloud/tail"

	"strconv"
	"time"
	"fmt"
	"github.com/qnib/qwatch/inputs"
)


var (
   rx = map[string]string{
	   "graphite": `(?P<metric>[a-z\-\.\_0-9]+)\s+(?P<value>[0-9\.]+)\s+(?P<time>\d+)`,
   }

)

func Run(qChan fTypes.QChan, cfg config.Config) {
	log.Println("[II] Start file collector")
	fPath, err := cfg.String("collector.file.path")
	if err != nil {
		log.Println("[EE] No file path for collector.file.path set")
		return
	}
	fileReopen, err := cfg.BoolOr("collector.file.reopen", true)
	t, err := tail.TailFile(fPath, tail.Config{Follow: true, ReOpen: fileReopen})
	if err != nil {
		log.Printf("[WW] File collector failed to open %s: %s", fPath, err)
	}
	mForm, _ := cfg.StringOr("collector.file.format", "graphite")
	regX := regexp.MustCompile(rx[mForm])
	dim := make(map[string]string)
	for line := range t.Lines {
		m := qinput.GetParams(regX, line.Text)
		if len(m) == 0  {
			continue
		}
		val, err := strconv.ParseFloat(m["value"], 64)
		t, tErr := strconv.Atoi(m["time"])
		if err == nil && tErr == nil {
			qm := qtypes.NewExt("input", "file", m["metric"], qtypes.Gauge, val, dim, time.Unix(int64(t),0), false)
			qChan.Data.Send(qm)
		} else {
			fmt.Printf("err:%s, tErr:%s\n", err, tErr)
		}
	}
}
