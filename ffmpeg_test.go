package ffmpeg_test

import (
	"fmt"
	"log"
	"maps"
	"slices"

	"github.com/northbright/ffmpeg"
)

func ExampleGetFormat() {
	f, err := ffmpeg.GetFormat("videos/cat-subtitled.mkv")
	if err != nil {
		log.Printf("GetFormat() error: %v", err)
		return
	}

	fmt.Printf("format:\n")
	fmt.Printf("filename: %s\n", f.Filename)
	fmt.Printf("nb_streams: %d\n", f.StreamsNum)
	fmt.Printf("size: %s\n", f.Size)         // size is in string format.
	fmt.Printf("duration: %s\n", f.Duration) // duration may not exist(e.g. for MKV)

	if len(f.Tags) > 0 {
		// Sort keys.
		keys := slices.Sorted(maps.Keys(f.Tags))

		fmt.Printf("tags:\n")
		for _, k := range keys {
			fmt.Printf("%s: %s\n", k, f.Tags[k])
		}
	}

	// Output:
	// format:
	// filename: videos/cat-subtitled.mkv
	// nb_streams: 4
	// size: 3146998
	// duration: 9.659000
	// tags:
	// COMPATIBLE_BRANDS: isomiso2avc1mp41
	// ENCODER: Lavf60.16.100
	// MAJOR_BRAND: isom
	// MINOR_VERSION: 512
}

func ExampleGetVideoStreams() {
	streams, err := ffmpeg.GetVideoStreams("videos/cat-subtitled.mkv")
	if err != nil {
		log.Printf("GetVideoStreams() error: %v", err)
		return
	}

	fmt.Printf("video streams: %v", streams)

	// Output:
	// video streams: [{{0 h264 video} 720 1280 30/1  {1}}]
}

func ExampleGetAudioStreams() {
	streams, err := ffmpeg.GetAudioStreams("videos/cat-subtitled.mkv")
	if err != nil {
		log.Printf("GetAudioStreams() error: %v", err)
		return
	}

	fmt.Printf("audio streams: %v", streams)

	// Output:
	// audio streams: [{{1 aac audio} 1 mono  {1}}]
}

func ExampleGetSubtitleStreams() {
	streams, err := ffmpeg.GetSubtitleStreams("videos/cat-subtitled.mkv")
	if err != nil {
		log.Printf("GetSubtitleStreams() error: %v", err)
		return
	}

	fmt.Printf("subtitle streams: %v", streams)

	// Output:
	// subtitle streams: [{{2 subrip subtitle} {{1} 0} {eng English}} {{3 subrip subtitle} {{0} 0} {chi Chinese}}]
}
