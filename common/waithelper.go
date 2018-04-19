package common

import (
	"log"
	"time"
)

func WaitUntil(condition func() bool, timeout time.Duration) {
	channel := make(chan bool)
	go waitUntilRoutine(condition, channel)
	select {
	case <-channel:
		return
	case <-time.After(timeout):
		log.Panic("Waiter has timed out")
		return
	}
}

func waitUntilRoutine(condition func() bool, channel chan bool) {
	for {
		if condition() {
			channel <- true
			return
		}
		<-time.After(time.Second)
	}
}
