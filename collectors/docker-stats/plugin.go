package main

import (
	"C"
	"log"
	"fmt"

	"github.com/docker/docker/client"
	"golang.org/x/net/context"
	"github.com/docker/docker/api/types"
	"github.com/fsouza/go-dockerclient"
	"github.com/pkg/errors"
	"github.com/zpatrick/go-config"

	"github.com/qnib/qcollect-ng/types"
	fTypes "github.com/qnib/qframe/types"
)

// ContainerListener spwan a goroutine per container and streams the stats into the metrics-channel
func ContainerListener(cli *docker.Client, outChan chan qtypes.Metric, id, name string) {
	errChannel := make(chan error, 1)
	statsChannel := make(chan *docker.Stats)

	opts := docker.StatsOptions{
		ID:     id,
		Stats:  statsChannel,
		Stream: true,
	}

	go func() {
		errChannel <- cli.Stats(opts)
	}()

	for {
		stats, ok := <-statsChannel
		if !ok {
			err := errors.New(fmt.Sprintf("Bad response getting stats for container: %s", id))
			log.Println(err.Error())
			return
		}
		dim := map[string]string{
			"container_id": id,
			"container_name": name,
		}
		val := float64(stats.CPUStats.CPUUsage.TotalUsage)
		qm := qtypes.NewExt("input", "docker-stats", "cpu-total", qtypes.Gauge, val, dim, stats.Read, false)
		outChan <-qm
	}
}

func ListenDispatcher(outChan chan qtypes.Metric, dockerHost string) {
	cntClient, err := docker.NewClient(dockerHost)
	if err != nil {
		log.Printf("[EE] Could not connect fsouza/go-dockerclient to '%s': %v", dockerHost, err)
		return
	}
	// Filter start/stop event of a container
	engineCli, err := client.NewClient(dockerHost, "v1.25", nil, nil)
	if err != nil {
		log.Printf("[EE] Could not connect docker/docker/client to '%s': %v", dockerHost, err)
		return
	}
	info, err := engineCli.Info(context.Background())
	if err != nil {
		log.Printf("[EE] Error during Info(): %v >err> %s", info, err)
		return
	} else {
		log.Printf("[II] Connected to '%s' w/ ServerVersion:'%s'", info.ID, info.ServerVersion)
	}
	msgs, errs := engineCli.Events(context.Background(), types.EventsOptions{})
	for {
		select {
		case dMsg := <-msgs:
			if dMsg.Type == "container" {
				switch dMsg.Action {
				case "start":
					log.Printf("[II] Container started ID:%s", dMsg.ID)
					go ContainerListener(cntClient, outChan, dMsg.ID, dMsg.Actor.Attributes["name"])
				case "die":
					log.Printf("[II] Container died ID:%s", dMsg.ID)
				default:
					//log.Printf("[DD] Unused Action: %s", dMsg.Action)
					continue
				}
			}
		case dErr := <-errs:
			if dErr != nil {
				log.Printf("[EE] %v", dErr)
			}
		}
	}
}

func Run(qChan fTypes.QChan, cfg config.Config) {
	log.Println("[II] Start docker-stats collector")
	dockerHost, _ := cfg.StringOr("collector.docker-stats.docker-host", "unix:///var/run/docker.sock")
	outChan := make(chan qtypes.Metric, 50)

	go ListenDispatcher(outChan, dockerHost)

	for {
		qm := <-outChan
		qChan.Data.Send(qm)
	}

}
