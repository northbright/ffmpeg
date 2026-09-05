package srt_test

import (
	"log"

	"github.com/northbright/ffmpeg/srt"
)

func ExampleWriteFile() {
	subtitles := []srt.SubTitle{
		srt.SubTitle{"00:00:00,000", "00:00:05,090", "What's mimao's playing?"},
		srt.SubTitle{"00:00:05,100", "00:00:08,200", "Does he realy like it?"},
	}

	f := "output/eng.srt"
	if err := srt.WriteFile(f, subtitles); err != nil {
		log.Printf("srt.WriteFile(%s) error: %v", f, err)
		return
	}

	subtitles = []srt.SubTitle{
		srt.SubTitle{"00:00:00,000", "00:00:05,090", "咪毛在玩啥？"},
		srt.SubTitle{"00:00:05,100", "00:00:08,200", "他真的喜欢玩这个？"},
	}

	f := "output/chi.srt"
	if err := srt.WriteFile(f, subtitles); err != nil {
		log.Printf("srt.WriteFile(%s) error: %v", f, err)
		return
	}

	// Output:
}
