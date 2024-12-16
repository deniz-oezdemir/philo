package main

import (
	"fmt"
	"sync"
	"time"
)

type Fork struct{ sync.Mutex }

type Philo struct {
	id					int
	leftFork, rightFork	*Fork
	config				*Config
}

var startTime time.Time

func (p Philo) log(action string) {
	elapsed := time.Since(startTime).Milliseconds()
	fmt.Printf("%d %d %s\n", elapsed, p.id, action)
}

func (p Philo) dine() {
	for {
		p.leftFork.Lock()
		p.log("has taken a fork")
		p.rightFork.Lock()
		p.log("has taken a fork")

		p.log("is eating")
		time.Sleep(time.Millisecond * time.Duration(p.config.timeEat))

		p.rightFork.Unlock()
		p.leftFork.Unlock()

		p.log("is sleeping")
		time.Sleep(time.Millisecond * time.Duration(p.config.timeSleep))

		p.log("is thinking")
	}
}

func main() {
	config, err := parseArgs()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Config: %+v\n", *config)

	//TODO: if 1 philo: pick up fork, wait timeDie, die

	forks := make([]*Fork, config.numPhilos)
	for i := 0; i < config.numPhilos; i++ {
		forks[i] = new(Fork)
	}

	startTime = time.Now()

	philos := make([]*Philo, config.numPhilos)
	for i := 0; i < config.numPhilos; i++ {
		philos[i] = &Philo{ id: i + 1, leftFork: forks[i], rightFork: forks[(i + 1) % config.numPhilos], config: config}
		go philos[i].dine()
		time.Sleep(time.Microsecond * time.Duration(1))
	}

	select {} // keep main goroutine alive as otherwise philos' goroutines terminate with main goroutine
}
