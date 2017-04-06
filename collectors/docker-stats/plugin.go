package main

import (
	"C"
	"log"
	"fmt"
	"sync"
	"regexp"

	"github.com/docker/docker/client"
	"golang.org/x/net/context"
	"github.com/docker/docker/api/types"
	"github.com/fsouza/go-dockerclient"
	"github.com/pkg/errors"
	"github.com/zpatrick/go-config"

	"github.com/qnib/qcollect-ng/utils"
	"github.com/qnib/qcollect-ng/types"
	fTypes "github.com/qnib/qframe/types"
)


var SuperMap = make(map[string]ContainerSupervisor)

// struct to keep info and channels to goroutine
// -> get heartbeats so that we know it's still alive
// -> allow for gracefully shutdown of the supervisor
type ContainerSupervisor struct {
	CntID 	string 			 // ContainerID
	CntName string			 // sanatized name of container
	Com 	chan interface{} // Channel to communicate with goroutine
	cli 	*docker.Client
	qChan 	fTypes.QChan
	IgnoreRegex *regexp.Regexp
}

func NewCntSuper(cli *docker.Client, qChan fTypes.QChan, cntID, cntName string, ignReg *regexp.Regexp) ContainerSupervisor {
	return ContainerSupervisor{
		CntID: cntID,
		CntName: cntName,
		Com: make(chan interface{}),
		cli: cli,
		qChan: qChan,
		IgnoreRegex: ignReg,
	}
}

// Run spwan a goroutine per container and streams the stats into the metrics-channel
func (cs ContainerSupervisor) Run() {
	log.Printf("[II] Start listener for already running '%s' [%s]", cs.CntName, cs.CntID)
	errChannel := make(chan error, 1)
	statsChannel := make(chan *docker.Stats)

	opts := docker.StatsOptions{
		ID:     cs.CntID,
		Stats:  statsChannel,
		Stream: true,
	}

	go func() {
		errChannel <- cs.cli.Stats(opts)
	}()

	for {
		select {
		case msg := <-cs.Com:
			switch msg {
			case "died":
				log.Printf("[DD] Container [%s]->'%s' died -> BYE!", cs.CntID, cs.CntName)
				return
			default:
				log.Printf("[DD] Container [%s]->'%s' got message from cs.Com: %v\n", cs.CntID, cs.CntName, msg)
			}
		case stats, ok := <-statsChannel:
			if !ok {
				err := errors.New(fmt.Sprintf("Bad response getting stats for container: %s", cs.CntID))
				log.Println(err.Error())
				return
			}

			dim := map[string]string{
				"container_id":   cs.CntID,
				"container_name": cs.CntName,
				"service_name":   "none",
				"task_slot":      "none",
				"task_id":        "none",
			}
			task, err := qutils.ContainerNameExtractService([]string{cs.CntName})
			if err == nil {
				dim["task_id"] = task.TaskID
				dim["task_slot"] = task.Slot
				dim["service_name"] = task.Name
			}

			/*
			//////// Panic when this is active: https://github.com/qnib/qcollect-ng/issues/14
			pre := qutils.TransformFsouzaToDocker(stats.PreCPUStats)
			cur := qutils.TransformFsouzaToDocker(stats.CPUStats)
			cpuStats := qutils.DiffCPUStats(pre, cur)
			cs.qChan.Data.Send(qtypes.NewExt("input", "docker-stats", "cpu.system.ms", qtypes.Gauge, float64(cpuStats.SystemUsage/10000000), dim, stats.Read, false))
			cs.qChan.Data.Send(qtypes.NewExt("input", "docker-stats", "cpu.usage.ms", qtypes.Gauge, float64(cpuStats.CPUUsage.TotalUsage/10000000), dim, stats.Read, false))
			*/
			cs.qChan.Data.Send(qtypes.NewExt("input", "docker-stats", "memory.usage.bytes", qtypes.Gauge, float64(stats.MemoryStats.Usage), dim, stats.Read, false))
			cs.qChan.Data.Send(qtypes.NewExt("input", "docker-stats", "memory.limit.bytes", qtypes.Gauge, float64(stats.MemoryStats.Limit), dim, stats.Read, false))
			cs.qChan.Data.Send(qtypes.NewExt("input", "docker-stats", "pid.current.count", qtypes.Gauge, float64(stats.PidsStats.Current), dim, stats.Read, false))
			cs.qChan.Data.Send(qtypes.NewExt("input", "docker-stats", "net.rx.bytes", qtypes.Gauge, float64(stats.Network.RxBytes), dim, stats.Read, false))
			cs.qChan.Data.Send(qtypes.NewExt("input", "docker-stats", "net.tx.bytes", qtypes.Gauge, float64(stats.Network.TxBytes), dim, stats.Read, false))
		}
	}
}

func ListenDispatcher(qChan fTypes.QChan, ignReg *regexp.Regexp, dockerHost string, wg sync.WaitGroup) {
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
	// Initialize already running containers
	cnts, err := engineCli.ContainerList(context.Background(), types.ContainerListOptions{})
	for _, cnt := range cnts {
		SuperMap[cnt.ID] = NewCntSuper(cntClient, qChan, cnt.ID, qutils.SanatizeContainerName(cnt.Names), ignReg)
		go SuperMap[cnt.ID].Run()
	}

	msgs, errs := engineCli.Events(context.Background(), types.EventsOptions{})
	for {
		select {
		case dMsg := <-msgs:
			if dMsg.Type == "container" {
				switch dMsg.Action {
				case "start":
					log.Printf("[II] Container [%s]->'%s' started", dMsg.ID, dMsg.Actor.Attributes["name"])
					SuperMap[dMsg.ID] = NewCntSuper(cntClient, qChan, dMsg.ID, dMsg.Actor.Attributes["name"], ignReg)
					go SuperMap[dMsg.ID].Run()
				case "die", "resize", "attach":
					//log.Printf("[II] Container [%s]->'%s' died", dMsg.ID, dMsg.Actor.Attributes["name"])
					SuperMap[dMsg.ID].Com <- dMsg.Action
				default:
					log.Printf("[DD] Container [%s]->'%s' unused Action: %s", dMsg.ID, dMsg.Actor.Attributes["name"], dMsg.Action)
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
	log.Println("[II] Start docker-stats collector %s")
	dockerHost, _ := cfg.StringOr("collector.docker-stats.docker-host", "unix:///var/run/docker.sock")
	var wg sync.WaitGroup
	wg.Add(1)
	ignReg, _ := cfg.StringOr("collector.docker-stats.ignore-regex", "^[a-z]+_[a-z]+$")
	go ListenDispatcher(qChan, regexp.MustCompile(ignReg), dockerHost, wg)
	wg.Wait()

}
