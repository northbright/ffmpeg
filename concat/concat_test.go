package concat_test

import (
	"log"

	"github.com/northbright/ffmpeg/concat"
)

func ExampleConcatCommand() {

	datas := []concat.ClipData{
		concat.ClipData{
			File:     "../images/op.jpg",
			Start:    "",
			End:      "00:00:03",
			Subtitle: "Good Times with Maomi & Mimao",
			IsMuted:  true,
		},
		concat.ClipData{
			File:     "../videos/01.MP4",
			Start:    "",
			End:      "00:00:03",
			Subtitle: "Mido's tickling Mimao and he's enjoying...",
			IsMuted:  false,
		},
		concat.ClipData{
			File:     "../videos/02.MOV",
			Start:    "",
			End:      "",
			Subtitle: "Mimao's playing the toy.",
			IsMuted:  false,
		},
		concat.ClipData{
			File:     "../videos/03.MOV",
			Start:    "00:00:01",
			End:      "00:00:08",
			Subtitle: "It's hard to brush Maomi's teeth...",
			IsMuted:  false,
		},
		concat.ClipData{
			File:     "../videos/04.MOV",
			Start:    "00:00:02",
			End:      "",
			Subtitle: "What's this??!!\nCan I eat it?",
			IsMuted:  false,
		},
		concat.ClipData{
			File:     "../images/ed.jpg",
			Start:    "",
			End:      "00:00:03",
			Subtitle: "Mimao likes lying on father's bed...😂\nMusic by penguinmusic: Better Day",
			IsMuted:  true,
		},
	}

	w := 1920
	h := 1080
	fps := 30
	output := "../output/good-times-with-maomi-and-mimao.mp4"

	cmd, pre, _, err := concat.ConcatCommand(datas, output, w, h, fps, false)
	if err != nil {
		log.Printf("concat.ConcatCommand() error: %v", err)
		return
	}

	log.Printf("concat.ConcatCommand() OK.")
	log.Printf("command:\n%s\n", cmd.String())

	// Run pre func.
	if err := pre(); err != nil {
		log.Printf("pre error: %v", err)
		return
	}

	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("cmd.CombinedOutput() error: %v\nout:\n%s\n", err, string(out))
	}

	log.Printf("cmd.CombinedOutput() OK.\nout:\n%s\n", string(out))
	// Run post func.
	//post()

	// Output:
}
