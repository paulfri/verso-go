package http

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/versolabs/citra/tasks"
)

func (c *Controller) crawl(ctx *gin.Context) {
	task, err := tasks.NewFeedParseTask(2)

	if err != nil {
		fmt.Printf("could not create task: %v", err)
	}

	info, err := c.asynq.Enqueue(task)
	if err != nil {
		log.Fatalf("could not enqueue task: %v", err)
	}

	log.Printf("enqueued task: id=%s queue=%s", info.ID, info.Queue)
}
