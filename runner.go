package main

import (
	"context"
	crand "crypto/rand"
	"fmt"
	"log"
	"math"
	"math/big"
	"math/rand"
	"os"
	"time"
)

type JobRunner struct {
	conf   *JobConfig
	logger *log.Logger
}

func newRunner(conf *JobConfig) JobRunner {
	seed, _ := crand.Int(crand.Reader, big.NewInt(math.MaxInt64))
	rand.Seed(seed.Int64())

	logger := log.New(os.Stdout, fmt.Sprintf("[JOB-%d]", rand.Intn(255)), log.Ldate|log.Ltime)

	return JobRunner{
		conf:   conf,
		logger: logger,
	}
}

func (j *JobRunner) Start() {
	j.logger.Println("[WORK] start")

	done := make(chan struct{})

	go j.executeScript(done)

	select {
	case <-done:
		j.logger.Println("[WORK] end")
	}
}

func (j *JobRunner) executeScript(done chan struct{}) {
	defer close(done)

	j.logger.Printf("[BEFORE] start")
	for _, c := range j.conf.BeforeScript {
		b, err := c.Run()
		if err != nil {
			j.logger.Printf("[BEFORE] error: %+v", string(b))
			return
		}
	}

	ctx := context.Background()

	j.logger.Printf("[JOB] start")
	for _, job := range j.conf.Jobs {
		job_ctx, _ := context.WithTimeout(ctx, j.conf.JobSetting.Timeout)

		start := time.Now()
		b, err := job.Run(job_ctx)

		end := time.Now()
		j.logger.Printf("[JOB] time: %+v", end.Sub(start))

		if err != nil {
			j.logger.Printf("[JOB] error: %+v", string(b))
			j.logger.Printf("[JOB] error: %+v", err)
		} else {
			j.logger.Printf("[JOB] output: %+v", string(b))
		}
	}

	j.logger.Printf("[AFTER] start")
	for _, c := range j.conf.AfterScript {
		b, err := c.Run()
		if err != nil {
			j.logger.Printf("[AFTER] error: %+v", string(b))
			return
		}
	}
}
