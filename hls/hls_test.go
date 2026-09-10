package hls_test

import (
	"context"
	"log"
	"path/filepath"

	"github.com/northbright/ffmpeg/hls"
)

func ExampleConcat() {

	// Segment.
	src := filepath.Join("../videos", "cat.mp4")
	duration := 4
	outDir := filepath.Join("../output")

	output, err := hls.Segment(context.Background(), src, duration, outDir)
	if err != nil {
		log.Printf("Segment() error: %v", err)
	} else {
		log.Printf("Segment() OK")
	}

	log.Printf("output:\n%s", output)

	// Concat.
	s := []string{"00000.ts", "00001.ts", "00002.ts"}
	var tsFiles []string

	for _, f := range s {
		tsFiles = append(tsFiles, filepath.Join("../output", f))
	}

	dst := filepath.Join("../output", "cat-concat.mp4")

	if output, err = hls.Concat(context.Background(), tsFiles, dst, true); err != nil {
		log.Printf("Concat() error: %v", err)
	} else {
		log.Printf("Concat() OK")
	}

	log.Printf("output:\n%s", output)

	// Output:
}
