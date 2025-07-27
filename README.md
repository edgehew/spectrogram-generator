# YouTube Spectrogram Generator

This project is a Go application that allows users to generate an audio spectrogram from a YouTube video URL. The application downloads the audio from the specified YouTube video URL and creates a spectrogram image in PNG format.

This project relies on ffmepg to convert downloaded audio files to WAV format.

## Project Structure

```
youtube-spectrogram
├── src
│   ├── downloader.go    # Functions for downloading audio from YouTube
│   ├── main.go          # Entry point of the application
│   ├── spectrogram.go   # Functions for generating the audio spectrogram
│   ├── utils.go         # Utility functions for checking if a file exists
│   └── wavConverter.go  # Uses ffmpeg to convert an audio file to wav format
└── README.md            # Project documentation
```

## Setup Instructions

1. **Clone the repository:**
   ```
   git clone https://github.com/yourusername/youtube-spectrogram.git
   cd youtube-spectrogram
   ```

2. **Install Go:** Make sure you have Go installed on your machine. You can download it from [golang.org](https://golang.org/dl/).

3. **Install ffmpeg** Make sure you have ffmpeg installed
   ```
   # Example:
   sudo apt install ffmpeg
   ```

4. **Install dependencies:** Navigate to the `src` directory and run:
   ```
   go mod tidy
   ```

5. **Build the application:** Run the following command
   ```
   go build src/*.go
   ```

## Usage

To run the application, use the following command:

```
go run src/*.go --help
// or after building the application
./downloader --help
```

```
Usage:
  -hop int
        Hop size for the spectrogram (default 512)
  -output string
        Output path for the spectrogram image (default "spectrogram.png")
  -url string
        YouTube video URL
  -window int
        FFT window size (default 1024)
```

Replace `<YouTube-Video-URL>` with the actual URL of the YouTube video you want to process.

## Example

```
go run src/*.go -url "https://www.youtube.com/watch?v=dQw4w9WgXcQ" -output rickroll.png
// After building
./downloader -url "https://www.youtube.com/watch?v=dQw4w9WgXcQ" -output rickroll.png
```

This command will download the audio from the specified YouTube video and generate a spectrogram image named `rickroll.png` in the current directory.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.