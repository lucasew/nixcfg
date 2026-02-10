package genetic

import (
	"math"

	"workspaced/pkg/palette/api"
)

// Individual represents a candidate palette solution
type Individual struct {
	colors  []api.LAB
	fitness float64
}

// calculateFitness evaluates how good a palette is
// Based on Stylix Stylix/Palette.hs fitness function
func calculateFitness(ind Individual, imageColors []api.LAB, polarity api.Polarity) float64 {
	var score float64

	// Primary scale similarity (base00-07 should be similar to each other)
	primarySim := maxDeltaE(ind.colors[:8])
	score -= primarySim / 10.0 // Penalize high maximum distance (want low spread)

	// Accent differentiation (base08-0F should be VERY different from each other)
	accentDiff := minDeltaE(ind.colors[8:16])
	score += accentDiff * 2.0 // Double weight - want distinct colors

	// CRITICAL: Penalize duplicate or near-duplicate accents very heavily
	duplicatePenalty := calculateDuplicatePenalty(ind.colors[8:16])
	score -= duplicatePenalty * 50.0 // Massive penalty for duplicates

	// Penalize accents that are too similar (under threshold)
	if accentDiff < 20.0 {
		score -= (20.0 - accentDiff) * 3.0 // Heavy penalty for similar accents
	}

	// Image similarity (colors should appear in the image)
	imageSim := calculateImageSimilarity(ind.colors, imageColors)
	score += imageSim * 5.0 // Reduced weight to allow more variety

	// Contrast requirements (critical for readability)
	contrastScore := calculateContrastScore(ind.colors)
	score += contrastScore * 5.0 // Strong weight on readability

	// Color diversity (hue variation in accents)
	hueDiversity := calculateHueDiversity(ind.colors[8:16])
	score += hueDiversity * 3.0 // Reward varied hues

	// Lightness scheme matching
	if polarity != api.PolarityAny {
		lightnessError := calculateLightnessError(ind.colors, polarity)
		score -= lightnessError
	}

	return score
}

// maxDeltaE finds the maximum perceptual distance between any two colors
func maxDeltaE(colors []api.LAB) float64 {
	if len(colors) < 2 {
		return 0
	}

	maxDist := 0.0
	for i := 0; i < len(colors); i++ {
		for j := i + 1; j < len(colors); j++ {
			dist := api.DeltaE(colors[i], colors[j])
			if dist > maxDist {
				maxDist = dist
			}
		}
	}
	return maxDist
}

// minDeltaE finds the minimum perceptual distance between any two colors
func minDeltaE(colors []api.LAB) float64 {
	if len(colors) < 2 {
		return 0
	}

	minDist := math.MaxFloat64
	for i := 0; i < len(colors); i++ {
		for j := i + 1; j < len(colors); j++ {
			dist := api.DeltaE(colors[i], colors[j])
			if dist < minDist {
				minDist = dist
			}
		}
	}
	return minDist
}

// calculateImageSimilarity measures how well the palette matches the image colors
func calculateImageSimilarity(paletteColors []api.LAB, imageColors []api.LAB) float64 {
	if len(imageColors) == 0 {
		return 0
	}

	// For each palette color, find closest image color
	totalDist := 0.0
	for _, palColor := range paletteColors {
		minDist := math.MaxFloat64
		for _, imgColor := range imageColors {
			dist := api.DeltaE(palColor, imgColor)
			if dist < minDist {
				minDist = dist
			}
		}
		totalDist += minDist
	}

	// Average distance (lower is better, so negate for score)
	avgDist := totalDist / float64(len(paletteColors))
	return -avgDist
}

// calculateLightnessError measures deviation from target lightness pattern
// Based on Stylix Stylix/Palette.hs lines 82-94
func calculateLightnessError(colors []api.LAB, polarity api.Polarity) float64 {
	var targetLightnesses []float64

	if polarity == api.PolarityDark {
		// Dark theme: background dark, foreground light
		targetLightnesses = []float64{10, 30, 45, 65, 75, 90, 95, 95}
	} else if polarity == api.PolarityLight {
		// Light theme: background light, foreground dark
		targetLightnesses = []float64{90, 70, 55, 35, 25, 10, 5, 5}
	}

	if len(targetLightnesses) == 0 {
		return 0
	}

	// Calculate error for base00-07 (primary scale)
	errorSum := 0.0
	for i := 0; i < 8 && i < len(colors); i++ {
		diff := colors[i].L - targetLightnesses[i]
		errorSum += math.Abs(diff)
	}

	return errorSum
}

// calculateContrastScore rewards palettes with good contrast for readability
func calculateContrastScore(colors []api.LAB) float64 {
	if len(colors) < 16 {
		return 0
	}

	var score float64

	// Base00 (background) vs Base05 (foreground) - must have high contrast
	bgFgContrast := math.Abs(colors[0].L - colors[5].L)
	if bgFgContrast >= 50 {
		score += 10.0 // Excellent contrast
	} else if bgFgContrast >= 40 {
		score += 5.0 // Good contrast
	} else {
		score -= (50 - bgFgContrast) // Penalize poor contrast
	}

	// Base07 (light background) vs Base02 (selection) - moderate contrast
	lightBgSelectionContrast := math.Abs(colors[7].L - colors[2].L)
	if lightBgSelectionContrast >= 15 && lightBgSelectionContrast <= 40 {
		score += 3.0 // Good selection visibility
	}

	// Accent colors (Base08-0F) vs background (Base00) - should be visible
	minAccentContrast := 100.0
	for i := 8; i < 16; i++ {
		contrast := math.Abs(colors[i].L - colors[0].L)
		if contrast < minAccentContrast {
			minAccentContrast = contrast
		}
	}
	if minAccentContrast >= 25 {
		score += 5.0 // All accents visible on background
	} else {
		score -= (25 - minAccentContrast) / 2 // Penalize invisible accents
	}

	return score
}

// calculateDuplicatePenalty heavily penalizes duplicate or near-duplicate colors
func calculateDuplicatePenalty(colors []api.LAB) float64 {
	penalty := 0.0

	for i := 0; i < len(colors); i++ {
		for j := i + 1; j < len(colors); j++ {
			diff := api.DeltaE(colors[i], colors[j])

			// If colors are identical or nearly identical
			if diff < 5.0 {
				penalty += 10.0 // Very high penalty per duplicate
			} else if diff < 10.0 {
				penalty += 5.0 // High penalty for very similar
			} else if diff < 15.0 {
				penalty += 2.0 // Moderate penalty for similar
			}
		}
	}

	return penalty
}

// calculateHueDiversity measures hue variation in accent colors
// Returns higher score for colors spread across color wheel
func calculateHueDiversity(colors []api.LAB) float64 {
	if len(colors) < 2 {
		return 0
	}

	// Calculate variance in A and B channels (chromaticity)
	var aSum, bSum float64
	for _, c := range colors {
		aSum += c.A
		bSum += c.B
	}
	aMean := aSum / float64(len(colors))
	bMean := bSum / float64(len(colors))

	var aVariance, bVariance float64
	for _, c := range colors {
		aVariance += math.Pow(c.A-aMean, 2)
		bVariance += math.Pow(c.B-bMean, 2)
	}
	aVariance /= float64(len(colors))
	bVariance /= float64(len(colors))

	// Higher variance = more color diversity
	diversity := math.Sqrt(aVariance + bVariance)

	// Normalize to reasonable range (0-20)
	return math.Min(diversity/5.0, 20.0)
}

// scoredIndividual pairs an individual with its fitness score
type scoredIndividual struct {
	individual Individual
	fitness    float64
}

// scorePop calculates fitness for entire population and sorts by fitness
func scorePop(population []Individual, imageColors []api.LAB, polarity api.Polarity) []scoredIndividual {
	scored := make([]scoredIndividual, len(population))

	for i, ind := range population {
		fitness := calculateFitness(ind, imageColors, polarity)
		scored[i] = scoredIndividual{
			individual: ind,
			fitness:    fitness,
		}
	}

	// Sort by fitness (descending - higher is better)
	for i := 0; i < len(scored); i++ {
		for j := i + 1; j < len(scored); j++ {
			if scored[j].fitness > scored[i].fitness {
				scored[i], scored[j] = scored[j], scored[i]
			}
		}
	}

	return scored
}
