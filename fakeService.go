package main

import (
	"errors"
	"math/rand"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

func ExecuteFakeTask(wg *sync.WaitGroup) {
	defer wg.Done()
	time.Sleep(time.Second)
	log.Info("fake task 1 done")
}

func ExecuteFakeTaskWithReturn(wg *sync.WaitGroup, channel chan any) {
	defer wg.Done()
	rand.Seed(time.Now().UnixNano())
	if rand.Intn(2) > 0 {
		log.Info("fake task 2 error")
		channel <- errors.New("I lost my way")
	} else {
		log.Info("fake task 2 error")
		channel <- 42
	}
}
