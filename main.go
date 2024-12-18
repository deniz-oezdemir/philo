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
	personalTimeDie		time.Time
	mealsEaten			int
	satisfied			bool
	mu					sync.Mutex
}

var (
	config			*Config
	startTime		time.Time
	satisfiedPhilos	int
	dinnerEnd		bool
	philoDead		bool
)

func (p Philo) log(action string) {
	elapsed := time.Since(startTime).Milliseconds()
	fmt.Printf("%d %d %s\n", elapsed, p.id, action)
}

func (p *Philo) dine() {
	for {
		p.leftFork.Lock()
		p.log("has taken a fork")
		p.rightFork.Lock()
		p.log("has taken a fork")

		p.mu.Lock()
		p.personalTimeDie = time.Now().Add(time.Duration(config.timeDie) * time.Millisecond)
		p.mealsEaten += 1
		p.mu.Unlock()

		p.log("is eating")
		preciseSleep(int64(config.timeEat))

		p.rightFork.Unlock()
		p.leftFork.Unlock()

		p.log("is sleeping")
		preciseSleep(int64(config.timeSleep))

		p.log("is thinking")
		preciseSleep(int64(1)) // prevent philosophers from taking another philosphers fork who needs to eat more urgently
	}
}

func monitor(philos []*Philo) {
	satisfiedPhilos = 0
	for !dinnerEnd && !philoDead {
		for i := 0; i < config.numPhilos; i++ {

			// check whether philosopher is dead
			philos[i].mu.Lock()
			if philos[i].personalTimeDie.Before(time.Now()) {
				philoDead = true
				philos[i].log("died")
				philos[i].mu.Unlock()
				return
			} else {
				philos[i].mu.Unlock()
			}

			// check whether individual philosopher is satisfied and if so check whether all philosophers are satisfied
			philos[i].mu.Lock()
			if philos[i].mealsEaten == config.numMeals && !philos[i].satisfied {
				philos[i].satisfied = true
				philos[i].mu.Unlock()
				satisfiedPhilos += 1
				if satisfiedPhilos == config.numPhilos {
					dinnerEnd = true
					elapsed := time.Since(startTime).Milliseconds()
					fmt.Printf("%d all philosophers have eaten %d times\n", elapsed, config.numMeals)
				}
			} else {
				philos[i].mu.Unlock()
			}
		}
	}
}

func main() {
	var err error
	config, err = parseArgs()
	if err != nil {
		fmt.Println(err)
		return
	}

	forks := make([]*Fork, config.numPhilos)
	for i := 0; i < config.numPhilos; i++ {
		forks[i] = new(Fork)
	}

	dinnerEnd, philoDead = false, false
	startTime = time.Now()

	philos := make([]*Philo, config.numPhilos)
	for i := 0; i < config.numPhilos; i++ {
		philos[i] = &Philo{ id: i + 1,
			leftFork: forks[i],
			rightFork: forks[(i + 1) % config.numPhilos],
			personalTimeDie: startTime.Add(time.Duration(config.timeDie) * time.Millisecond),
			mealsEaten: 0,
			satisfied: false}
		go philos[i].dine()
		time.Sleep(time.Millisecond) // delay next philo goroutine shortly to mitigate each philo holding onto its own fork (deadlock)
	}

	monitor(philos)

	//select {} // only necessary for testing when we want to keep the main alive without a monitor as otherwise the philo goroutines exit with the main
}
