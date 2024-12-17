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
		preciseSleep(int64(p.config.timeEat))

		p.rightFork.Unlock()
		p.leftFork.Unlock()

		p.log("is sleeping")
		preciseSleep(int64(p.config.timeSleep))

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

	// TODO: if 1 philo: pick up fork, wait timeDie, die

	forks := make([]*Fork, config.numPhilos)
	for i := 0; i < config.numPhilos; i++ {
		forks[i] = new(Fork)
	}

	startTime = time.Now()

	// TODO: add monitor thread to check for end of dinner in case of dead philo or each philo ate numMeals

	philos := make([]*Philo, config.numPhilos)
	for i := 0; i < config.numPhilos; i++ {
		philos[i] = &Philo{ id: i + 1, leftFork: forks[i], rightFork: forks[(i + 1) % config.numPhilos], config: config}
		go philos[i].dine()
		time.Sleep(time.Nanosecond) // delay next philo goroutine shortly to mitigate each philo holding on to its own fork (deadlock)
	}

	select {} // keep main goroutine alive as otherwise philos' goroutines terminate with it
}
