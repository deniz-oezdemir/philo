package main

import (
	"fmt"
	"os"
	"strconv"
)

// Config holds the configuration for the simulation from the command line arguments
type Config struct {
	numPhilos	int
	timeDie		int
	timeEat		int
	timeSleep	int
	numMeals	int
}

// parseArgs parses command-line arguments into a Config struct.
func parseArgs() (*Config, error) {
	if len(os.Args) != 5 && len(os.Args) != 6 {
		return nil, fmt.Errorf("Usage: philo number_of_philosophers time_to_die time_to_eat time_to_sleep [number_of_times_each_philosopher_must_eat]")
	}

	args := make([]int, len(os.Args)-1)
	for i, arg := range os.Args[1:] {
		val, err := strconv.Atoi(arg)
		if err != nil {
			return nil, fmt.Errorf("Invalid argument %s: %v", arg, err)
		}
		if val < 0 {
			return nil, fmt.Errorf("Argument %s must be non-negative", arg)
		}
		args[i] = val
	}

	config := &Config{
		numPhilos:	args[0],
		timeDie:	args[1],
		timeEat:	args[2],
		timeSleep:	args[3],
		numMeals:	-1,
	}

	if len(args) == 5 {
		config.numMeals = args[4]
	}

	return config, nil
}
