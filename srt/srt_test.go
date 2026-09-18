package srt_test

import (
	"context"
	"fmt"
	"log"

	"github.com/northbright/ffmpeg/srt"
)

func ExampleAddSoftSub() {
	subtitles := []srt.Subtitle{
		srt.Subtitle{"00:00:00,000", "00:00:05,090", "What's mimao playing?"},
		srt.Subtitle{"00:00:05,100", "00:00:09,500", "Does he really like it?"},
	}

	f1 := "../output/eng.srt"
	if err := srt.WriteFile(f1, subtitles); err != nil {
		log.Printf("srt.WriteFile(%s) error: %v", f1, err)
		return
	}

	subtitles = []srt.Subtitle{
		srt.Subtitle{"00:00:00,000", "00:00:05,090", "咪毛在玩啥？"},
		srt.Subtitle{"00:00:05,100", "00:00:09,500", "他真的喜欢玩这个？"},
	}

	f2 := "../output/chi.srt"
	if err := srt.WriteFile(f2, subtitles); err != nil {
		log.Printf("srt.WriteFile(%s) error: %v", f2, err)
		return
	}

	ctx := context.Background()

	// Output MKV.
	// Add English subtitle.
	input := "../videos/cat.mp4"
	output := "../output/cat-eng.mkv"
	isDefault := true

	out, err := srt.AddSoftSub(ctx, input, f1, "eng", "English", isDefault, output, true)
	if err != nil {
		log.Printf("srt.AddSoftSub() for English error: %v\noutput:\n%s", err, out)
		return
	}

	log.Printf("srt.AddSoftSub() for English OK. output:\n%s", out)

	// Add Chinese subtitle.
	input = "../output/cat-eng.mkv"
	output = "../output/cat-eng-chi.mkv"
	isDefault = false

	out, err = srt.AddSoftSub(ctx, input, f2, "chi", "Chinese", isDefault, output, true)
	if err != nil {
		log.Printf("srt.AddSoftSub() for Chinese error: %v\noutput:\n%s", err, out)
		return
	}

	log.Printf("srt.AddSoftSub() for Chinese OK. output:\n%s", out)

	// Output MP4.
	// Add English subtitle.
	input = "../videos/cat.mp4"
	output = "../output/cat-eng.mp4"
	isDefault = true

	out, err = srt.AddSoftSub(ctx, input, f1, "eng", "English", isDefault, output, true)
	if err != nil {
		log.Printf("srt.AddSoftSub() for English error: %v\noutput:\n%s", err, out)
		return
	}

	log.Printf("srt.AddSoftSub() for English OK. output:\n%s", out)

	// Add Chinese subtitle.
	input = "../output/cat-eng.mp4"
	output = "../output/cat-eng-chi.mp4"
	isDefault = false

	out, err = srt.AddSoftSub(ctx, input, f2, "chi", "Chinese", isDefault, output, true)
	if err != nil {
		log.Printf("srt.AddSoftSub() for Chinese error: %v\noutput:\n%s", err, out)
		return
	}

	log.Printf("srt.AddSoftSub() for Chinese OK. output:\n%s", out)

	// Set Chinese as the default subtitle stream.
	input = "../output/cat-eng-chi.mkv"
	output = "../output/cat-eng-chi-new.mkv"
	id := 1 // s:1 is Chinese subtitle stream.
	out, err = srt.SetDefaultSub(ctx, input, id, output, true)
	if err != nil {
		log.Printf("srt.SetDefaultSub(): set Chinese as default error: %v\noutput:\n%s", err, out)
		return
	}

	log.Printf("srt.SetDefaultSub(): set Chinese as default OK. output:\n%s", out)

	// Output:
}

func ExampleHardSub_VideoFilter() {
	subtitles := []srt.Subtitle{
		srt.Subtitle{"00:00:00,000", "00:00:05,090", "What's mimao playing?"},
		srt.Subtitle{"00:00:05,100", "00:00:09,500", "Does he really like it?"},
	}

	f := "../output/eng.srt"
	if err := srt.WriteFile(f, subtitles); err != nil {
		log.Printf("srt.WriteFile(%s) error: %v", f, err)
		return
	}

	hs := srt.NewHardSub(
		f,
		srt.FontName("Arial"),                  // default: "Arial".
		srt.FontSize(20),                       // default: 16.
		srt.PrimaryColour(0xFF, 0xD3, 0x80, 0), // default: white.
		srt.OutlineColour(0x2C, 0x48, 0x75, 0), // default: black.
		srt.BackColour(0x00, 0x20, 0x2E, 0),    // default: black.
		srt.Outline(3),                         // default: 2,
		srt.Shadow(2),                          // default: 0,
		srt.Bold(true),                         // default: false.
		srt.Italic(true),                       // default: false.
		srt.Alignment(1),                       // default: 2, range: 1 - 9(num keyboard layout).
		srt.MarginV(30),                        // default: 10.
	)

	// Get subtitles video filter used in ffmpeg.
	fmt.Printf("hard sub video filter: %s\n", hs.VideoFilter())

	// Output:
	// hard sub video filter: subtitles='../output/eng.srt':force_style='Alignment=1,BackColour=&H002e2000,Bold=1,FontName=Arial,FontSize=20,Italic=1,MarginV=30,Outline=3,OutlineColour=&H0075482c,PrimaryColour=&H0080d3ff,Shadow=2'
}

func ExampleAddHardSub() {
	subtitles := []srt.Subtitle{
		srt.Subtitle{"00:00:00,000", "00:00:05,090", "咪毛在玩啥？"},
		srt.Subtitle{"00:00:05,100", "00:00:09,500", "他真的喜欢玩这个？"},
	}

	srtFile := "../output/chi.srt"
	if err := srt.WriteFile(srtFile, subtitles); err != nil {
		log.Printf("srt.WriteFile(%s) error: %v", srtFile, err)
		return
	}

	input := "../videos/cat.mp4"
	output := "../output/cat-hard-sub-chi.mp4"
	out, err := srt.AddHardSub(
		context.Background(),
		input,
		srtFile,
		output,
		true,
		srt.FontName("Arial"),                  // default: "Arial".
		srt.FontSize(20),                       // default: 16.
		srt.PrimaryColour(0xFF, 0xD3, 0x80, 0), // default: white.
		srt.OutlineColour(0x2C, 0x48, 0x75, 0), // default: black.
		srt.BackColour(0x00, 0x20, 0x2E, 0),    // default: black.
		srt.Outline(3),                         // default: 2,
		srt.Shadow(2),                          // default: 0,
		srt.Bold(true),                         // default: false.
		srt.Italic(true),                       // default: false.
		srt.Alignment(1),                       // default: 2, range: 1 - 9(num keyboard layout).
		srt.MarginV(30),                        // default: 10.
	)

	if err != nil {
		log.Printf("srt.AddHardSub() for Chinese error: %v\noutput:\n%s", err, out)
		return
	}

	log.Printf("srt.AddHardSub() for Chinese OK. output:\n%s", out)

	// Output:
}
