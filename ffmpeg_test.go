package ffmpeg_test

import (
	"fmt"
	"log"

	"github.com/northbright/ffmpeg"
)

func ExampleGetFormat() {
	f, err := ffmpeg.GetFormat("videos/cat-subtitled.mkv")
	if err != nil {
		log.Printf("GetFormat() error: %v", err)
		return
	}

	fmt.Printf("format: %v\n", f)

	// Output:
	// format: &{videos/cat-subtitled.mkv 4 3146526 9.649000}
}

func ExampleGetVideoStreams() {
	streams, err := ffmpeg.GetVideoStreams("videos/cat-subtitled.mkv")
	if err != nil {
		log.Printf("GetVideoStreams() error: %v", err)
		return
	}

	fmt.Printf("video streams: %v", streams)

	// Output:
	// video streams: [{{0 h264 video} 720 1280 30/1 {{1} 0 0} {{00:00:09.636000000} Core Media Video Lavc60.31.102 libx264}}]
}

func ExampleGetAudioStreams() {
	streams, err := ffmpeg.GetAudioStreams("videos/cat-subtitled.mkv")
	if err != nil {
		log.Printf("GetAudioStreams() error: %v", err)
		return
	}

	fmt.Printf("audio streams: %v", streams)

	// Output:
	// audio streams: [{{1 vorbis audio} 1 mono {{1} 0 0 0 0 0} {{00:00:09.649000000} Core Media Audio}}]
}

func ExampleGetSubtitleStreams() {
	streams, err := ffmpeg.GetSubtitleStreams("videos/cat-subtitled.mkv")
	if err != nil {
		log.Printf("GetSubtitleStreams() error: %v", err)
		return
	}

	fmt.Printf("subtitle streams: %v", streams)

	// Output:
	// subtitle streams: [{{2 subrip subtitle} {{1} 0} {{00:00:08.200000000} eng English}} {{3 subrip subtitle} {{0} 0} {{00:00:08.200000000} chi 中文}}]
}
