package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
)

// convertToWAV converts an audio file to WAV format using FFmpeg
func convertToWAV(inputPath, outputPath string) error {
	// Check if input file exists
	if _, err := os.Stat(inputPath); os.IsNotExist(err) {
		return fmt.Errorf("input file does not exist: %s", inputPath)
	}

	// Create FFmpeg command to convert M4A to WAV
	cmd := exec.Command("ffmpeg", "-i", inputPath, "-acodec", "pcm_s16le", "-ar", "44100", "-ac", "2", outputPath)

	// Capture any error output
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("error creating stderr pipe: %v", err)
	}

	// Start the command
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("error starting FFmpeg: %v", err)
	}

	// Read error output
	ffmpegErrorLog, _ := io.ReadAll(stderr)
	if len(ffmpegErrorLog) > 0 {
		fmt.Printf("FFmpeg warnings/errors: %s\n", ffmpegErrorLog)
	}

	// Wait for the command to complete
	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("FFmpeg conversion failed: %v", err)
	}

	// Verify output file was created
	if !fileExists(outputPath) {
		return fmt.Errorf("output file was not created: %s", outputPath)
	}

	return nil
}
