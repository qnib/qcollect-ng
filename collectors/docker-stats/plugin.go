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
	"github.com/qnib/qcollect-ng/utils"
	fTypes "github.com/qnib/qframe/types"
	"sync"
)

// ContainerListener spwan a goroutine per container and streams the stats into the metrics-channel
func ContainerListener(cli *docker.Client, qChan fTypes.QChan, id, name string) {
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
		pre := qutils.TransformFsouzaToDocker(stats.PreCPUStats)
		cur := qutils.TransformFsouzaToDocker(stats.CPUStats)
		cpuStats := qutils.DiffCPUStats(pre, cur)
		qChan.Data.Send(qtypes.NewExt("input", "docker-stats", "cpu.system.ms", qtypes.Gauge, float64(cpuStats.SystemUsage/10000000), dim, stats.Read, false))
		qChan.Data.Send(qtypes.NewExt("input", "docker-stats", "cpu.usage.ms", qtypes.Gauge, float64(cpuStats.CPUUsage.TotalUsage/10000000), dim, stats.Read, false))
		qChan.Data.Send(qtypes.NewExt("input", "docker-stats", "memory.usage.bytes", qtypes.Gauge, float64(stats.MemoryStats.Usage), dim, stats.Read, false))
		qChan.Data.Send(qtypes.NewExt("input", "docker-stats", "memory.limit.bytes", qtypes.Gauge, float64(stats.MemoryStats.Limit), dim, stats.Read, false))
		qChan.Data.Send(qtypes.NewExt("input", "docker-stats", "pid.current.count", qtypes.Gauge, float64(stats.PidsStats.Current), dim, stats.Read, false))
		qChan.Data.Send(qtypes.NewExt("input", "docker-stats", "net.rx.bytes", qtypes.Gauge, float64(stats.Network.RxBytes), dim, stats.Read, false))
		qChan.Data.Send(qtypes.NewExt("input", "docker-stats", "net.tx.bytes", qtypes.Gauge, float64(stats.Network.TxBytes), dim, stats.Read, false))
	}
}

func ListenDispatcher(qChan fTypes.QChan, dockerHost string) {
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
					go ContainerListener(cntClient, qChan, dMsg.ID, dMsg.Actor.Attributes["name"])
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
	var wg sync.WaitGroup
	wg.Add(1)
	go ListenDispatcher(qChan, dockerHost)

	wg.Wait()

}
