package main

import (
	"C"
	"net"
	"fmt"
	"strings"
	"log"

	"github.com/zpatrick/go-config"
	"github.com/qnib/qcollect-ng/types"
	fTypes "github.com/qnib/qframe/types"
	"github.com/qnib/qframe/utils"
	"time"
)



// Run fetches everything from the Data channel and flushes it to opentsdb
func Run(qChan fTypes.QChan, cfg config.Config) {
	bg := qChan.Data.Join()
	inStr, err := cfg.StringOr("handler.opentsdb.inputs", "")
	inputs := strings.Split(inStr, ",")
	server, _ := cfg.StringOr("handler.opentsdb.server", "localhost")
	targetForm, _ := cfg.StringOr("handler.opentsdb.target", "")
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
		switch val.(type) {
		case qtypes.Metric:
			qm := val.(qtypes.Metric)
			if len(inputs) != 0 && ! qutils.IsInput(inputs, qm.Source) {
				continue
			}
			fmt.Fprintf(conn, qm.ConvertToOpenTSDBHandler(targetForm))
		}
	}
}

