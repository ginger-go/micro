package osstat

import (
	"log"

	"github.com/ginger-go/micro"
	"github.com/ginger-go/micro/plugins/auth"
	"github.com/ginger-go/micro/plugins/logger"
	"github.com/mackerelio/go-osstat/cpu"
	"github.com/mackerelio/go-osstat/memory"
)

// MonitorMemory monitor memory usage
// alertPercent is >= 0 and <= 1
func MonitorMemory(e *micro.Engine, alertFunc func()) {
	micro.Cron(e, "@every 10s", func() {
		m, err := memory.Get()
		if err != nil {
			log.Println("failed to get os memory info", err)
			return
		}
		if float64(m.Used)/float64(m.Total) > ALERT_PERCENTAGE_MEMORY {
			logger.Error(auth.GetSystemID(), "", "", "memory used percent alert: ", float64(m.Used)/float64(m.Total))
			if !DISABLE_ALERT_MEMORY {
				alertFunc()
			}
		}
	})
}

// MonitorCPU monitor cpu usage
// alertPercent is >= 0 and <= 1
func MonitorCPU(e *micro.Engine, alertFunc func()) {
	micro.Cron(e, "@every 10s", func() {
		c, err := cpu.Get()
		if err != nil {
			log.Println("failed to get os cpu info", err)
			return
		}
		if float64(c.User+c.System)/float64(c.Total) > ALERT_PERCENTAGE_CPU {
			logger.Error(auth.GetSystemID(), "", "", "cpu used percent alert: ", float64(c.User+c.System)/float64(c.Total))
			if !DISABLE_ALERT_CPU {
				alertFunc()
			}
		}
	})
}
