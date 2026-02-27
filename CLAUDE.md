# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Go-based video black border cropping tool that uses FFmpeg's `cropdetect` functionality to automatically detect and crop black borders from videos. The tool is designed to be self-contained with embedded FFmpeg binaries and supports cross-platform deployment (Windows, macOS, Linux).

## Architecture and Structure

The project follows a modular Go structure:
- `main.go`: Application entry point that initializes the command-line interface
- `internal/cmd/root.go`: Implements the CLI using Cobra with argument parsing and business logic orchestration
- `internal/ffmpeg/cropdetect.go`: Core FFmpeg functionality for detecting and cropping black borders
- `internal/ffmpeg/embed_*.go`: Platform-specific embedding of FFmpeg binaries using Go's embed.FS
- `internal/ffmpeg/binaries/`: Directory containing FFmpeg binaries for different platforms

The application works in two phases:
1. Detection phase: Uses FFmpeg's cropdetect filter to analyze a specified duration of video (default 10 seconds) to detect black borders
2. Cropping phase: Applies the detected crop parameters to create a new video file with borders removed

## Key Features

- Automatic black border detection using FFmpeg's cropdetect filter
- Self-contained deployment with embedded FFmpeg
- Cross-platform support (Windows, macOS, Linux)
- Real-time progress display with speed and time estimates
- Command-line interface with intuitive parameters

## Building and Running

### Development Setup
```bash
# Install dependencies
make deps

# Build for current platform
make build

# Or run directly with Go
go run main.go [input_file] [output_file]
```

### Cross-Platform Builds
```bash
# Build for specific platforms
make windows    # Windows executable
make macos      # macOS executable
make linux      # Linux executable
make build-all  # All platforms

# Build with CGO disabled (required for compatibility)
CGO_ENABLED=0 go build -o video_crop
```

### FFmpeg Requirements
The application can use either:
1. System-installed FFmpeg (if available in PATH)
2. Embedded FFmpeg binaries (automatically extracted when needed)

To download embedded binaries:
```bash
make download-ffmpeg
# Or use platform-specific scripts:
# ./download_ffmpeg.sh (Unix)
# powershell -ExecutionPolicy Bypass -File download_ffmpeg.ps1 (Windows)
```

## Common Commands

```bash
# Basic usage
video_crop input.mp4                    # Auto-generate output filename
video_crop input.mp4 output.mp4         # Specify output filename
video_crop --input input.mp4 --output output.mp4

# Advanced options
video_crop -i input.mp4 -o output.mp4 -y    # Overwrite existing output
video_crop -i input.mp4 -o output.mp4 -d 15 # Detect for 15 seconds
video_crop --help                           # Show help

# Available flags
--input, -i          Input video file
--output, -o         Output video file (default: [input]_cropped.[ext])
--overwrite, -y      Overwrite existing output file
--detect-time, -d    Duration in seconds for black border detection (default: 10)
```

## Important Considerations

1. **CGO must be disabled**: `CGO_ENABLED=0` must be set when building to ensure cross-platform compatibility and avoid issues like "missing LC_UUID load command" on macOS
2. **Large binary size**: The embedded FFmpeg increases the executable size to 50-80MB
3. **Detection duration**: Recommended detection time is 5-15 seconds; too short may miss borders, too long increases processing time
4. **FFmpeg priority**: The app prioritizes system FFmpeg when available, falling back to embedded binaries if needed
5. **Video quality**: Audio streams are copied without re-encoding (`-c:a copy`), preserving original quality

## Development Guidelines

When contributing to this project:
1. Maintain cross-platform compatibility by disabling CGO
2. Preserve the modular structure (CLI in cmd/, core functionality in ffmpeg/)
3. Follow Go coding conventions and maintain clear documentation
4. Test builds on multiple platforms before submitting changes
5. Keep embedded binaries updated in the internal/ffmpeg/binaries/ directory