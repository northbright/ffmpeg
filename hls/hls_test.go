package hls_test

import (
	"context"
	"log"
	"path/filepath"

	"github.com/northbright/ffmpeg/hls"
)

func ExampleSegment() {
	src := filepath.Join("../videos", "cat.MOV")
	duration := 4
	outDir := filepath.Join("../videos", "output")

	output, err := hls.Segment(context.Background(), src, duration, outDir)
	if err != nil {
		log.Printf("Segment() error: %v", err)
	} else {
		log.Printf("Segment() OK")
	}

	log.Printf("output:\n%s", output)

	// Output:
}

func ExampleConcat() {
	s := []string{"00000.ts", "00001.ts", "00002.ts"}
	var tsFiles []string

	for _, f := range s {
		tsFiles = append(tsFiles, filepath.Join("../videos", "../output", f))
	}

	dst := filepath.Join("../videos", "../output", "cat-concat.MOV")

	output, err := hls.Concat(context.Background(), tsFiles, dst, true)
	if err != nil {
		log.Printf("Concat() error: %v", err)
	} else {
		log.Printf("Concat() OK")
	}

	log.Printf("output:\n%s", output)

	// Output:
}
