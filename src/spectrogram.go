package main

import (
	"fmt"
	"image/color"
	"math"
	"os"

	"github.com/fogleman/gg"
	"github.com/mjibson/go-dsp/fft"
	"github.com/youpy/go-wav"
)

func GenerateSpectrogram(inputPath, outputPath string, windowSize, hopSize int) error {
	// Open WAV file
	file, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("failed to open WAV file: %v", err)
	}
	defer file.Close()

	// Create WAV reader
	reader := wav.NewReader(file)
	format, err := reader.Format()
	if err != nil {
		return fmt.Errorf("failed to read WAV format: %v", err)
	}

	// Read samples in chunks
	samples := make([]float64, 0)
	bufferSize := uint32(1024) // Number of samples to read at once

	for {
		buffer, err := reader.ReadSamples(bufferSize)
		if err != nil && err.Error() != "EOF" {
			return fmt.Errorf("failed to read samples: %v", err)
		}
		for i := range len(buffer) {
			// Use first channel (mono or left channel for stereo)
			samples = append(samples, float64(buffer[i].Values[0])/math.MaxInt16)
		}
		if uint32(len(buffer)) < bufferSize {
			break // End of file
		}
	}

	// Spectrogram parameters
	//sampleRate := int(format.SampleRate)
	numChannels := int(format.NumChannels)
	width := len(samples)/(hopSize*numChannels) + 1
	height := windowSize / 2 // Nyquist frequency
	spectrogram := make([][]float64, width)

	// Compute FFT for each window
	for i := 0; i < width; i++ {
		start := i * hopSize
		end := start + windowSize
		if end > len(samples) {
			end = len(samples)
		}
		if start >= len(samples) {
			break
		}

		// Extract window and apply Hann windowing
		window := make([]float64, windowSize)
		for j := 0; j < windowSize && start+j < len(samples); j++ {
			hann := 0.5 * (1 - math.Cos(2*math.Pi*float64(j)/float64(windowSize-1)))
			window[j] = samples[start+j] * hann
		}

		// Compute FFT
		fftResult := fft.FFTReal(window)
		spectrogram[i] = make([]float64, height)
		for j := 0; j < height; j++ {
			// Convert to magnitude (dB)
			mag := math.Sqrt(real(fftResult[j])*real(fftResult[j]) + imag(fftResult[j])*imag(fftResult[j]))
			if mag > 0 {
				spectrogram[i][j] = 20 * math.Log10(mag)
			} else {
				spectrogram[i][j] = -100
			}
		}
	}

	// Normalize spectrogram for visualization
	maxDB := -100.0
	minDB := 0.0
	for i := 0; i < len(spectrogram); i++ {
		for j := 0; j < len(spectrogram[i]); j++ {
			if spectrogram[i][j] > maxDB {
				maxDB = spectrogram[i][j]
			}
			if spectrogram[i][j] < minDB {
				minDB = spectrogram[i][j]
			}
		}
	}

	// Create colored image
	dc := gg.NewContext(width, height)
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			normalized := (spectrogram[x][y] - minDB) / (maxDB - minDB)
			// Map to color gradient (blue -> green -> yellow -> red)
			r, g, b := colorGradient(normalized)
			dc.SetColor(color.RGBA{R: uint8(r * 255), G: uint8(g * 255), B: uint8(b * 255), A: 255})
			dc.SetPixel(x, height-y-1)
		}
	}

	// Save spectrogram as PNG
	if err := dc.SavePNG(outputPath); err != nil {
		return fmt.Errorf("failed to save PNG: %v", err)
	}

	return nil
}

// colorGradient maps a normalized value (0 to 1) to a color gradient (blue -> green -> yellow -> red).
func colorGradient(value float64) (r, g, b float64) {
	if value < 0.25 {
		// Blue to cyan
		r = 0
		g = value * 4
		b = 1
	} else if value < 0.5 {
		// Cyan to green
		r = 0
		g = 1
		b = 1 - (value-0.25)*4
	} else if value < 0.75 {
		// Green to yellow
		r = (value - 0.5) * 4
		g = 1
		b = 0
	} else {
		// Yellow to red
		r = 1
		g = 1 - (value-0.75)*4
		b = 0
	}
	return
}
