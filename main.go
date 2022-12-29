package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type Scheduler struct {
	waitGroup    sync.WaitGroup
	workNotifyCh chan os.Signal
	workersCh    chan struct{}
	nWorkers     int
}

func main() {
	fmt.Println("Process ID:", os.Getpid())

	s := NewScheduler(5, 10)

	s.workerProcess()
	<-waitForExit()
	s.cleanup()
}

func NewScheduler(bufferSize, workers int) *Scheduler {
	wnCh := make(chan os.Signal, 1)
	signal.Notify(wnCh, syscall.SIGTERM)
	return &Scheduler{
		waitGroup:    sync.WaitGroup{},
		workNotifyCh: wnCh,
		workersCh:    make(chan struct{}, bufferSize),
		nWorkers:     workers,
	}
}

func waitForExit() chan struct{} {
	shutdownNotifyCh := make(chan struct{})
	shutdownCh := make(chan os.Signal, 1)
	signal.Notify(shutdownCh, os.Interrupt)
	go func() {
		defer close(shutdownNotifyCh)
		<-shutdownCh
	}()
	return shutdownNotifyCh
}

func (s *Scheduler) workerProcess() {
	go func() {
		for range s.workNotifyCh {
			s.workersCh <- struct{}{}
		}
	}()

	go func() {
		for i := 0; i < s.nWorkers; i++ {
			s.waitGroup.Add(1)
			go func(i int) {
				fmt.Println("launched worker process:", i)
				for {
					select {
					case _, open := <-s.workersCh:
						if !open {
							fmt.Println("shutting worker process;", i)
							s.waitGroup.Done()
							return
						}
					}
					fmt.Println("processing in worker process:", i)
					// DO WORK HERE
					time.Sleep(time.Second * 20)
				}
			}(i)
		}
	}()
}

func (s *Scheduler) cleanup() {
	fmt.Println("cleanuingup and shutting down")
	close(s.workNotifyCh)
	close(s.workersCh)
	s.waitGroup.Wait()
	fmt.Println("Goodbye!")
}
