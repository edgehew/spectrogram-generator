package main

import (
	"flag"
	"log"
	"os"
)

func main() {
	// Define command-line flags
	url := flag.String("url", "", "YouTube video URL")
	output := flag.String("output", "spectrogram.png", "Output path for the spectrogram image")
	windowSize := flag.Int("window", 1024, "FFT window size")
	hopSize := flag.Int("hop", 512, "Hop size for the spectrogram")
	flag.Parse()

	// Check if the URL is provided
	if *url == "" {
		log.Fatal("You must provide a YouTube video URL using the -url flag.")
	}

	// Download audio from the YouTube video
	audioPath := "/tmp/audio"
	err := DownloadAudio(*url, audioPath)
	if err != nil {
		log.Fatalf("Failed to get WAV audio File: %v", err)
	}
	filename := audioPath + ".wav"
	defer os.Remove(filename) // Clean up the audio file after generating the spectrogram

	// Generate the spectrogram
	err = GenerateSpectrogram(filename, *output, *windowSize, *hopSize)
	if err != nil {
		log.Fatalf("Failed to generate spectrogram: %v", err)
	}

	log.Println("Spectrogram generated successfully and saved to", *output)
}
