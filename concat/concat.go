package concat

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/northbright/ffmpeg"
	"github.com/northbright/ffmpeg/srt"
	"github.com/northbright/timestamp"
)

type Clip struct {
	inputID          int
	File             string `json:"file"`
	Start            string `json:"start"`
	End              string `json:"end"`
	Subtitle         string `json:"subtitle"`
	IsMuted          bool   `json:"is_muted,omitempty"`
	FadeInDuration   int    `json:"fade_in_duration,omitempty"`
	FadeOutDuration  int    `json:"fade_out_duration,omitempty"`
	AFadeInDuration  int    `json:"afade_in_duration,omitempty"`
	AFadeOutDuration int    `json:"afade_out_duration,omitempty"`
}

func ConcatCommand(clips []Clip, output string, w, h, fps int, isRankMode bool) (*exec.Cmd, error) {
	/*
		inputs := append(inputs, clips...)

		for id, input := range clips {
			clips[i].inputID = i
		}

		return nil, nil
	*/
	return nil, nil

}

// TrimParams returns parameters(start / end) for "trim" / "atrim" filter.
func TrimParams(start, end string) (string, error) {
	if start == end {
		return "", fmt.Errorf("start equals to end")
	}

	str := ""
	if start != "" {
		tsStart, err := timestamp.New(start)
		if err != nil {
			return "", fmt.Errorf("create start timestamp for %s error: %v", start, err)
		}
		str += fmt.Sprintf("start=%s:", tsStart.SecondStr())
	}

	if end != "" {
		tsEnd, err := timestamp.New(end)
		if err != nil {
			return "", fmt.Errorf("create end timestamp for %s error: %v", end, err)
		}
		str += fmt.Sprintf("end=%s:", tsEnd.SecondStr())
	} else {
		str = strings.TrimSuffix(str, ":")
	}

	return str, nil
}

func (c *Clip) TempSRTFile() string {
	f := strings.Replace(filepath.Base(c.File), filepath.Ext(c.File), ".srt", -1)
	f = filepath.Join(os.TempDir(), f)

	return f
}

func (c *Clip) Duration() (*timestamp.Timestamp, error) {
	// Get start and end timestamps.
	var (
		err        error
		start, end *timestamp.Timestamp
	)

	if c.Start == "" {
		c.Start = "00:00:00"
	}
	if start, err = timestamp.New(c.Start); err != nil {
		return nil, fmt.Errorf("new timestamp: %s error: %v", c.Start, err)
	}

	if c.End == "" {
		mf, err := ffmpeg.Format(c.File)
		if err != nil {
			return nil, fmt.Errorf("get format of %s error: %v", c.File, err)
		}

		// Use duration as end timestamp.
		if end, err = timestamp.New(mf.Duration); err != nil {
			return nil, fmt.Errorf("new timestamp: %s error: %v", mf.Duration, err)
		}
	} else {
		if end, err = timestamp.New(c.End); err != nil {
			return nil, fmt.Errorf("new timestamp: %s error: %v", c.End, err)
		}
	}

	// Get clip duration.
	d, err := end.Sub(start)
	if err != nil {
		return nil, fmt.Errorf("end.Sub(start) error: %v", err)
	}

	return d, nil
}

func (c *Clip) WriteTempSRTFile(duration *timestamp.Timestamp) error {
	f := c.TempSRTFile()

	d, err := c.Duration()
	if err != nil {
		return fmt.Errorf("get duration of %s error: %v", c.File, err)
	}

	subs := []srt.Subtitle{
		srt.Subtitle{
			Start: "00:00:00,000",
			End:   d.StringForSRT(),
		},
	}

	if err := srt.WriteFile(f, subs); err != nil {
		return fmt.Errorf("write temp srt file: %s error: %v", f, err)
	}

	return nil
}

func (c *Clip) RemoveTempSRTFile() error {
	f := c.TempSRTFile()

	if err := os.Remove(f); err != nil {
		return fmt.Errorf("remove temp srt file: %s error: %v", f, err)
	}

	return nil
}

func (c *Clip) VideoFilterChain(w, h, fps int, styles ...srt.Style) (*ffmpeg.FilterChain, error) {
	// Create the video filer chain.
	vOutLabel := fmt.Sprintf("[clip_%d_v]", c.inputID)
	v := ffmpeg.NewFilterChain(vOutLabel)

	// Add input.
	// "X:v:0" which X is the input ID.
	v.AddInputByID(c.inputID, "v", 0)

	// Get clip duration.
	d, err := c.Duration()
	if err != nil {
		return nil, fmt.Errorf("get duration of %s error: %v", c.File, err)
	}

	// Check if input file is image.
	isImage, err := ffmpeg.IsImage(c.File)
	if err != nil {
		return nil, fmt.Errorf("check if file: %s is image error: %v", c.File, err)
	}

	if isImage {
		// fps.
		FPS := fmt.Sprintf("fps=%d", fps)

		// loop.
		loop := fmt.Sprintf("loop=loop=%d:size=1", int(d.Second())*fps)

		// format.
		format := "format=pix_fmts=yuv420p"

		v.Chain(FPS).Chain(loop).Chain(format)
	}

	// scale.
	scale := fmt.Sprintf("scale=%d:%d:force_original_aspect_ratio=decrease", w, h)

	// pad.
	pad := fmt.Sprintf("pad=%d:%d:(ow-iw)/2:(oh-ih)/2", w, h)

	// setsar.
	setsar := "setsar=1:1"

	v.Chain(scale).Chain(pad).Chain(setsar)

	// trim.
	param, err := TrimParams(c.Start, c.End)
	if err != nil {
		return nil, fmt.Errorf("get trim params error: %v", err)
	}

	trim := "trim=" + param
	v.Chain(trim)

	// setpts.
	setpts := "setpts=PTS-STARTPTS"
	v.Chain(setpts)

	// subtitles.
	if c.Subtitle != "" {
		f := c.TempSRTFile()
		hs := srt.NewHardSub(f, styles...)
		sub := hs.VideoFilter()

		v.Chain(sub)
	}

	return v, nil
}

func (c *Clip) AudioFilterChain() (*ffmpeg.FilterChain, error) {
	if c.IsMuted {
		return nil, nil
	}

	// Create the audio filter chain.
	aOutLabel := fmt.Sprintf("[clip_%d_a]", c.inputID)
	a := ffmpeg.NewFilterChain(aOutLabel)

	// Add input.
	// "X:a:0" which X is the input ID.
	a.AddInputByID(c.inputID, "a", 0)

	// atrim.
	param, err := TrimParams(c.Start, c.End)
	if err != nil {
		return nil, fmt.Errorf("get trim params error: %v", err)
	}

	atrim := "atrim=" + param
	a.Chain(atrim)

	// asetpts.
	asetpts := "asetpts=PTS-STARTPTS"
	a.Chain(asetpts)

	return a, nil
}
