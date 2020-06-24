package utils

import (
	"gopkg.in/robfig/cron.v2"
)

// ExecuteCronJob excutes the given function, given the schedule provided
func ExecuteCronJob(schedule string, fnRun func()) {
	c := cron.New()
	c.AddFunc(schedule, fnRun)
	c.Start()
}
