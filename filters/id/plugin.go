package main

import (
	"C"
    "log"

	"github.com/zpatrick/go-config"
	"github.com/qnib/qcollect-ng/types"
	fTypes "github.com/qnib/qframe/types"
	"github.com/qnib/qframe/utils"
	"strings"
)

func Run(qChan fTypes.QChan, cfg config.Config) {
	log.Println("[II] Start 'id' filter")
	myId := qutils.GetGID()
	bg := qChan.Data.Join()
	for {
		val := bg.Recv()
		switch val.(type) {
		case qtypes.Metric:
			qm := val.(qtypes.Metric)
			if qm.SourceID == myId {
				continue
			}
			qm.Type = "filter"
			qm.Source = strings.Join(append(strings.Split(qm.Source, "->"), "id"), "->")
			qm.SourceID = myId
			qChan.Data.Send(qm)
		}
	}
}
