package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/kkdai/youtube/v2"
)

func DownloadAudio(url, audioPath string) error {
	client := youtube.Client{}
	video, err := client.GetVideo(url)
	if err != nil {
		return fmt.Errorf("failed to get YouTube video: %v", err)
	}
	log.Println("Video title:", video.Title)

	if video.Duration == 0 {
		return fmt.Errorf("video duration is zero, cannot proceed with download")
	}
	log.Println("Video duration:", video.Duration)
	if video.Duration.Minutes() > 10 {
		return fmt.Errorf("video duration is longer than 10 minutes, please use a shorter video")
	}

	// Find an audio-only format (e.g., mpeg, m4a)
	var format *youtube.Format
	var fileExt string
	for _, f := range video.Formats {
		if strings.Contains(f.MimeType, "audio") { // mp3
			log.Println("Available format:", f.MimeType, f.AudioSampleRate, f.Bitrate)
			format = &f
			fileExt = format.MimeType[strings.Index(format.MimeType, "/")+1:] // e.g., "mp3", "m4a"
			break
		}
	}
	if format == nil {
		return fmt.Errorf("no suitable audio format found for video: %s", video.Title)
	}

	log.Println("Downloading audio from YouTube video:", url)
	stream, _, err := client.GetStream(video, format)
	if err != nil {
		return fmt.Errorf("failed to get audio stream: %v", err)
	}
	log.Println("Selected audio format:", format.MimeType)

	tempFilename := audioPath + "." + fileExt
	file, err := os.Create(tempFilename)
	if err != nil {
		return fmt.Errorf("failed to create audio file: %v", err)
	}
	defer file.Close()
	defer os.Remove(tempFilename) // Clean up the audio file after converting to WAV file

	log.Println("Saving audio to:", tempFilename)
	_, err = file.ReadFrom(stream)
	if err != nil {
		return fmt.Errorf("failed to save audio file: %v", err)
	}

	tempWavPath := audioPath + ".wav"
	err = convertToWAV(tempFilename, tempWavPath)
	if err != nil {
		return fmt.Errorf("failed to convert audio to WAV format: %v", err)
	}

	if !fileExists(tempWavPath) {
		return fmt.Errorf("failed to convert audio to WAV format")
	}
	return err
}
