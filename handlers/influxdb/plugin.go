package main

import (
	"C"
	"log"
	"fmt"
	"strings"

	"github.com/zpatrick/go-config"
	"github.com/influxdata/influxdb/client/v2"

	"github.com/qnib/qcollect-ng/types"
	fTypes "github.com/qnib/qframe/types"
	"github.com/qnib/qframe/utils"

)

func createBatchPoitns(cfg config.Config) (bp client.BatchPoints, err error) {
	database, _ := cfg.StringOr("handler.influxdb.database", "qcollect")
	precision, _ := cfg.StringOr("handler.influxdb.precision", "s")
	bp, err = client.NewBatchPoints(client.BatchPointsConfig{
		Database:  database,
		Precision: precision,
	})
	return
}
func createConnection(cfg config.Config) (con client.Client, err error) {
	server, _ := cfg.StringOr("handler.influxdb.server", "localhost")
	port, _ := cfg.StringOr("handler.influxdb.port", "8086")
	username, _ := cfg.StringOr("handler.influxdb.username", "root")
	password, _ := cfg.StringOr("handler.influxdb.password", "root")
	// Create DB connection
	addr := fmt.Sprintf("http://%s:%s", server, port)
	con, err = client.NewHTTPClient(client.HTTPConfig{
		Addr:     addr,
		Username: username,
		Password: password,
	})
	if err != nil {
		log.Printf("[EE] Failed to establish connection to '%s': %v", addr, err)
	} else {
		log.Printf("[DD] Establish connection to '%s'", addr)
	}
	return
}

// Run fetches everything from the Data channel and flushes it to stdout
func Run(qChan fTypes.QChan, cfg config.Config) {
	bg := qChan.Data.Join()
	inStr, err := cfg.String("handler.influxdb.inputs")
	if err != nil {
		inStr = ""
	}
	inputs := strings.Split(inStr, ",")

	c, err := createConnection(cfg)
	if err != nil {
		return
	}
	// Create a new point batch
	bp, err := createBatchPoitns(cfg)
	if err != nil {
		return
	}
	for {
		val := bg.Recv()
		switch val.(type) {
		case qtypes.Metric:
			qm := val.(qtypes.Metric)
			if len(inputs) != 0 && ! qutils.IsInput(inputs, qm.Source) {
				//fmt.Printf("%s %-7s sType:%-6s sName:%-10s[%d] DROPED : %s\n", qm.TimeString(), qm.LogString(), qm.Type, qm.Source, qm.SourceID, qm.Msg)
				continue
			}
			// Create a point and add to batch
			fields := map[string]interface{}{
				"value": qm.Value,
			}

			pt, err := client.NewPoint(qm.Name, qm.Dimensions, fields, qm.Time)
			if err != nil {
				log.Printf("[EE] Failed to create new point: %v", err)
			}
			bp.AddPoint(pt)

			// Write the batch
			if err := c.Write(bp); err != nil {
				log.Printf("[EE] Failed to write batch: %v", err)
			}
		}
	}
}

