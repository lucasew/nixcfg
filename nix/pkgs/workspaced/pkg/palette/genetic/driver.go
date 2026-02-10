package genetic

import (
	"context"
	"image"
	"log/slog"
	"math/rand"

	"workspaced/pkg/palette/api"
)

type Driver struct{}

func (d *Driver) Name() string {
	return "genetic"
}

func (d *Driver) Extract(ctx context.Context, img image.Image, opts api.Options) (*api.Palette, error) {
	// Use deterministic RNG for reproducibility (like Stylix)
	rand.Seed(42)

	// 1. Sample colors from image
	colors := api.SampleImage(img, opts.MaxSamples)
	if len(colors) == 0 {
		return nil, ctx.Err()
	}
	slog.Info("sampled colors from image", "unique_colors", len(colors))

	// 2. Convert to LAB color space
	labColors := make([]api.LAB, len(colors))
	for i, c := range colors {
		labColors[i] = api.RGBToLAB(c)
	}
	slog.Info("converted to LAB color space")

	// 3. Initialize population
	population := initPopulation(labColors, opts.ColorCount, numSurvivors+numNewborns)
	slog.Info("initialized population", "size", len(population), "colors_per_palette", opts.ColorCount)

	// 4. Evolution loop
	generation := 0
	var prevBestFitness float64
	maxGenerations := 100 // Safety limit

	for generation < maxGenerations {
		// Check for context cancellation
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		// Calculate fitness for all individuals
		scored := scorePop(population, labColors, opts.Polarity)

		// Check convergence (fitness stopped improving)
		bestFitness := scored[0].fitness
		slog.Info("generation completed",
			"generation", generation,
			"best_fitness", bestFitness,
			"population_size", len(population))

		if generation > 0 && bestFitness == prevBestFitness {
			slog.Info("converged - fitness unchanged", "generations", generation)
			break
		}
		prevBestFitness = bestFitness

		// Select survivors (top numSurvivors)
		survivors := scored
		if len(survivors) > numSurvivors {
			survivors = survivors[:numSurvivors]
		}

		// Generate offspring via crossover + mutation
		population = evolve(survivors, labColors)

		generation++
	}

	// 5. Get best individual
	scored := scorePop(population, labColors, opts.Polarity)
	best := scored[0].individual
	slog.Info("evolution complete", "final_fitness", scored[0].fitness, "total_generations", generation)

	// 6. Map to base16/base24 palette
	pal := mapToPalette(best, opts.ColorCount)
	slog.Info("palette generated successfully")
	return pal, nil
}

// mapToPalette converts an individual to a base16/base24 palette
func mapToPalette(ind Individual, colorCount int) *api.Palette {
	pal := &api.Palette{}

	// Convert LAB colors back to hex strings
	hexColors := make([]string, len(ind.colors))
	for i, lab := range ind.colors {
		rgb := api.LABToRGB(lab)
		hexColors[i] = api.ToHex(rgb)
	}

	// Base16 colors (always present)
	if len(hexColors) >= 16 {
		pal.Base00 = hexColors[0]
		pal.Base01 = hexColors[1]
		pal.Base02 = hexColors[2]
		pal.Base03 = hexColors[3]
		pal.Base04 = hexColors[4]
		pal.Base05 = hexColors[5]
		pal.Base06 = hexColors[6]
		pal.Base07 = hexColors[7]
		pal.Base08 = hexColors[8]
		pal.Base09 = hexColors[9]
		pal.Base0A = hexColors[10]
		pal.Base0B = hexColors[11]
		pal.Base0C = hexColors[12]
		pal.Base0D = hexColors[13]
		pal.Base0E = hexColors[14]
		pal.Base0F = hexColors[15]
	}

	// Base24 extras (optional)
	if colorCount >= 24 && len(hexColors) >= 24 {
		pal.Base10 = hexColors[16]
		pal.Base11 = hexColors[17]
		pal.Base12 = hexColors[18]
		pal.Base13 = hexColors[19]
		pal.Base14 = hexColors[20]
		pal.Base15 = hexColors[21]
		pal.Base16 = hexColors[22]
		pal.Base17 = hexColors[23]
	}

	return pal
}
