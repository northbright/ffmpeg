package srt_test

import (
	"context"
	"log"

	"github.com/northbright/ffmpeg/srt"
)

func ExampleAddSoftSub() {
	subtitles := []srt.Subtitle{
		srt.Subtitle{"00:00:00,000", "00:00:05,090", "What's mimao's playing?"},
		srt.Subtitle{"00:00:05,100", "00:00:08,200", "Does he realy like it?"},
	}

	f1 := "../output/eng.srt"
	if err := srt.WriteFile(f1, subtitles); err != nil {
		log.Printf("srt.WriteFile(%s) error: %v", f1, err)
		return
	}

	subtitles = []srt.Subtitle{
		srt.Subtitle{"00:00:00,000", "00:00:05,090", "咪毛在玩啥？"},
		srt.Subtitle{"00:00:05,100", "00:00:08,200", "他真的喜欢玩这个？"},
	}

	f2 := "../output/chi.srt"
	if err := srt.WriteFile(f2, subtitles); err != nil {
		log.Printf("srt.WriteFile(%s) error: %v", f2, err)
		return
	}

	ctx := context.Background()

	// Add English subtitle.
	input := "../videos/cat.mp4"
	output := "../output/cat-eng.mkv"
	isDefault := true

	out, err := srt.AddSoftSub(ctx, input, f1, "eng", "English", isDefault, output, true)
	if err != nil {
		log.Printf("str.AddSoftSub() error: %v\noutput:\n%s", err, out)
		return
	}

	log.Printf("str.AddSoftSub() OK. output:\n%s", out)

	// Add Chinese subtitle.
	input = "../output/cat-eng.mkv"
	output = "../output/cat-eng-chi.mkv"
	isDefault = false

	out, err = srt.AddSoftSub(ctx, input, f2, "chi", "Chinese", isDefault, output, true)
	if err != nil {
		log.Printf("str.AddSoftSub() error: %v\noutput:\n%s", err, out)
		return
	}

	log.Printf("str.AddSoftSub() OK. output:\n%s", out)

	// Output:
}
