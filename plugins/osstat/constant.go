package osstat

import "github.com/ginger-go/env"

var (
	ALERT_PERCENTAGE_CPU    = env.Float64("ALERT_PERCENTAGE_CPU", 0.85)
	DISABLE_ALERT_CPU       = env.Bool("DISABLE_ALERT_CPU", false)
	ALERT_PERCENTAGE_MEMORY = env.Float64("ALERT_PERCENTAGE_MEMORY", 0.85)
	DISABLE_ALERT_MEMORY    = env.Bool("DISABLE_ALERT_MEMORY", false)
)
