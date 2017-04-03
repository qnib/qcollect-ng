package opentsdb

import (
	"C"
	"net"
	"fmt"
	"strings"
	"sort"
	"log"

	"github.com/zpatrick/go-config"
	"github.com/qnib/qcollect-ng/types"
	fTypes "github.com/qnib/qframe/types"
	"github.com/qnib/qframe/utils"
	"time"
)

func convertToOpenTSDBHandler(incomingMetric qtypes.Metric) (datapoint string) {
	//orders dimensions so datapoint keeps consistent name
	var keys []string
	dimensions := incomingMetric.GetDimensions(make(map[string]string))
	for k := range dimensions {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	datapoint = fmt.Sprintf("put %s", incomingMetric.Name)
	var dims []string
	for _, key := range keys {
		dims = append(dims, fmt.Sprintf("%s=%s", key, dimensions[key]))
	}
	if len(dimensions) == 0 {
		datapoint = fmt.Sprintf("%s %d %f\n", datapoint, incomingMetric.GetTime().Unix(), incomingMetric.Value)
	} else {
		datapoint = fmt.Sprintf("%s %s %d %f\n", datapoint, strings.Join(dims[:], ","), incomingMetric.GetTime().Unix(), incomingMetric.Value)
	}
	return datapoint
}


// Run fetches everything from the Data channel and flushes it to opentsdb
func Run(qChan fTypes.QChan, cfg config.Config) {
	bg := qChan.Data.Join()
	inStr, err := cfg.StringOr("handler.opentsdb.inputs", "")
	inputs := strings.Split(inStr, ",")
	server, _ := cfg.StringOr("handler.opentsdb.server", "localhost")
	port, _ := cfg.StringOr("handler.opentsdb.port", "4242")
	timeout, _ := cfg.IntOr("handler.opentsdb.timeout", 0)
	addr := fmt.Sprintf("%s:%s", server, port)
	conn, err := net.DialTimeout("tcp", addr, time.Duration(timeout) * time.Second)
	if err != nil {
		log.Printf("[EE] Failed to connect %s", addr)
		return
	} else {
		log.Printf("[II] Sucessfully connected to %s", addr)
	}
	for {
		val := bg.Recv()
		qm := val.(qtypes.Metric)
		if len(inputs) != 0 && ! qutils.IsInput(inputs, qm.Source) {
			continue
		}
		fmt.Fprintf(conn, convertToOpenTSDBHandler(qm))
	}
}

