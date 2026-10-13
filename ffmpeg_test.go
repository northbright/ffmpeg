package ffmpeg_test

import (
	"fmt"
	"log"
	"maps"
	"slices"

	"github.com/northbright/ffmpeg"
)

func ExampleFormat() {
	f, err := ffmpeg.Format("videos/cat-subtitled.mkv")
	if err != nil {
		log.Printf("Format() error: %v", err)
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

func ExampleVideoStreams() {
	streams, err := ffmpeg.VideoStreams("videos/cat-subtitled.mkv")
	if err != nil {
		log.Printf("VideoStreams() error: %v", err)
		return
	}

	fmt.Printf("video streams: %v", streams)

	// Output:
	// video streams: [{{0 h264 video} 720 1280 30/1  {1}}]
}

func ExampleAudioStreams() {
	streams, err := ffmpeg.AudioStreams("videos/cat-subtitled.mkv")
	if err != nil {
		log.Printf("AudioStreams() error: %v", err)
		return
	}

	fmt.Printf("audio streams: %v", streams)

	// Output:
	// audio streams: [{{1 aac audio} 1 mono  {1}}]
}

func ExampleSubtitleStreams() {
	streams, err := ffmpeg.SubtitleStreams("videos/cat-subtitled.mkv")
	if err != nil {
		log.Printf("SubtitleStreams() error: %v", err)
		return
	}

	fmt.Printf("subtitle streams: %v", streams)

	// Output:
	// subtitle streams: [{{2 subrip subtitle} {{1} 0} {eng English}} {{3 subrip subtitle} {{0} 0} {chi Chinese}}]
}

func ExampleIsImage() {
	files := []string{
		"images/cat.jpg",
		"images/cat.png",
		"images/cat.bmp",
		"images/cat.tiff",
		"images/cat.webp",
		"images/cat.jxl", // JPEGXL
		"images/cat.avif",
		"images/cat.heic",
		"videos/cat.mkv",
		"videos/cat.mp4",
	}

	for _, f := range files {
		b, err := ffmpeg.IsImage(f)
		if err != nil {
			log.Printf("IsImage(%s) error: %v", f, err)
			return
		}
		fmt.Printf("%s --> %v\n", f, b)
	}

	// Output:
	// images/cat.jpg --> true
	// images/cat.png --> true
	// images/cat.bmp --> true
	// images/cat.tiff --> true
	// images/cat.webp --> true
	// images/cat.jxl --> true
	// images/cat.avif --> true
	// images/cat.heic --> true
	// videos/cat.mkv --> false
	// videos/cat.mp4 --> false
}
