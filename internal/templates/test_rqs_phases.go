package templates

import (
	"strconv"
	"sync"

	log "github.com/sirupsen/logrus"
	"my5G-RANTester/config"
	"my5G-RANTester/internal/control_test_engine/gnb"
	"my5G-RANTester/internal/monitoring"
)

type Phase struct {
	Rate     int // requests per second
	Duration int // seconds
}

func TestRqsLoop(numRqs int, interval int) int64 {
	wg := sync.WaitGroup{}

	monitor := monitoring.Monitor{
		RqsL: 0,
		RqsG: 0,
	}

	cfg, err := config.GetConfig()
	if err != nil {
		log.Fatal("Error in get configuration")
	}

	ranPort := 1000
	for y := 1; y <= interval; y++ {
		monitor.InitRqsLocal()

		for i := 1; i <= numRqs; i++ {
			cfg.GNodeB.PlmnList.GnbId = gnbIdGenerator(i)
			cfg.GNodeB.ControlIF.Port = ranPort
			go gnb.InitGnbForLoadSeconds(cfg, &wg, &monitor)
			wg.Add(1)
			ranPort++
		}

		wg.Wait()
		log.Warn("[TESTER][GNB] AMF Responses per Second:", monitor.GetRqsLocal())
		monitor.SetRqsGlobal(monitor.GetRqsLocal())
	}

	return monitor.GetRqsGlobal()
}

func TestRqsPhases(phases []Phase) int64 {
	totalResponses := int64(0)
	for _, p := range phases {
		log.Warn("Starting phase: ", p.Rate, " req/s for ", p.Duration, "s")
		totalResponses += TestRqsLoop(p.Rate, p.Duration)
	}
	return totalResponses
}

func gnbIdGenerator(i int) string {
	var base string
	switch {
	case i < 10:
		base = "00000"
	case i < 100:
		base = "0000"
	default:
		base = "000"
	}
	return base + strconv.Itoa(i)
}
