package genetic

import (
	"math/rand"

	"workspaced/pkg/palette/api"
)

const (
	numSurvivors = 500
	numNewborns  = 49500
	mutationRate = 0.75
)

// initPopulation creates the initial random population with random colors
// Evolution will converge toward image colors through fitness function
func initPopulation(colors []api.LAB, count int, size int) []Individual {
	population := make([]Individual, size)

	for i := 0; i < size; i++ {
		individual := Individual{
			colors: make([]api.LAB, count),
		}

		// Generate completely random colors in LAB space
		for j := 0; j < count; j++ {
			individual.colors[j] = api.LAB{
				L: rand.Float64() * 100,        // Lightness: 0-100
				A: rand.Float64()*200 - 100,    // Green-Red: -100 to +100
				B: rand.Float64()*200 - 100,    // Blue-Yellow: -100 to +100
			}
		}

		population[i] = individual
	}

	return population
}

// crossover combines two parent individuals using alternating zip
// Based on Stylix Ai/Evolutionary.hs alternatingZip
func crossover(p1, p2 Individual) Individual {
	size := len(p1.colors)
	if len(p2.colors) < size {
		size = len(p2.colors)
	}

	offspring := Individual{
		colors: make([]api.LAB, size),
	}

	// Alternating zip: take from p1, then p2, alternating
	for i := 0; i < size; i++ {
		if i%2 == 0 {
			offspring.colors[i] = p1.colors[i]
		} else {
			offspring.colors[i] = p2.colors[i]
		}
	}

	return offspring
}

// mutate randomly replaces one color with probability 'rate'
// Based on Stylix Ai/Evolutionary.hs mutate function
func mutate(ind Individual, colors []api.LAB, rate float64) Individual {
	if len(colors) == 0 || rand.Float64() > rate {
		return ind
	}

	// Clone the individual
	mutated := Individual{
		colors: make([]api.LAB, len(ind.colors)),
	}
	copy(mutated.colors, ind.colors)

	// Replace one random color
	pos := rand.Intn(len(mutated.colors))
	mutated.colors[pos] = colors[rand.Intn(len(colors))]

	return mutated
}

// evolve creates next generation from survivors
// Based on Stylix Ai/Evolutionary.hs evolvePopulation
func evolve(survivors []scoredIndividual, imageColors []api.LAB) []Individual {
	newPopulation := make([]Individual, 0, numSurvivors+numNewborns)

	// Elitism: keep best individual unchanged
	if len(survivors) > 0 {
		newPopulation = append(newPopulation, survivors[0].individual)
	}

	// Generate offspring via crossover
	for i := 1; i < numSurvivors+numNewborns; i++ {
		// Select two random parents
		p1 := survivors[rand.Intn(len(survivors))].individual
		p2 := survivors[rand.Intn(len(survivors))].individual

		// Crossover
		offspring := crossover(p1, p2)

		// Mutate (skip first individual - elite)
		if i > 0 {
			offspring = mutate(offspring, imageColors, mutationRate)
		}

		newPopulation = append(newPopulation, offspring)
	}

	return newPopulation
}
