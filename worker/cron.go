package worker

import (
	"log"

	"github.com/robfig/cron/v3"
)

func (w Worker) ConfigureCron() {
	c := cron.New(cron.WithSeconds())

	c.AddFunc("@every 1m", func() {
		log.Println("cron!")
	})

	c.Run()
}
