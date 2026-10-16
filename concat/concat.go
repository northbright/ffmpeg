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

type ClipData struct {
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

type Clip struct {
	inputID     int
	isFromImage bool
	duration    *timestamp.Timestamp
	ClipData
}

func NewClip(inputID int, d ClipData) (*Clip, error) {
	var (
		err        error
		start, end *timestamp.Timestamp
	)

	c := &Clip{inputID: inputID, ClipData: d}

	// Check if clip is from image.
	if c.isFromImage, err = ffmpeg.IsImage(c.File); err != nil {
		return nil, fmt.Errorf("check if file: %s is image error: %v", c.File, err)
	}

	// Get clip duration.
	// Get start and end timestamps.
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
	if c.duration, err = end.Sub(start); err != nil {
		return nil, fmt.Errorf("end.Sub(start) error: %v", err)
	}

	return c, nil
}

// TrimParams returns parameters(start / end) for "trim" / "atrim" filter.
// It returns empty string when there's no need to trim(e.g. start == end).
func TrimParams(start, end string) (string, error) {
	if start == end {
		return "", nil
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

func (c *Clip) Duration() *timestamp.Timestamp {
	return c.duration
}

func (c *Clip) WriteTempSRTFile() error {
	f := c.TempSRTFile()

	d := c.Duration()

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

func (c *Clip) InputArgs() []string {
	var args []string

	if c.isFromImage {
		d := int(c.Duration().Second())
		args = append(args, "-loop", "1", "-t", fmt.Sprintf("%d", d), "-i", c.File)
	}

	args = append(args, "-i", c.File)
	return args
}

func (c *Clip) VideoFilterChain(w, h, fps int, styles ...srt.Style) (*ffmpeg.FilterChain, error) {
	// Create the video filer chain.
	vOutLabel := fmt.Sprintf("[v%d]", c.inputID)
	v := ffmpeg.NewFilterChain(vOutLabel)

	// Add input.
	// "X:v:0" which X is the input ID.
	v.AddInputByID(c.inputID, "v", 0)

	// fps for clip from image only.
	if c.isFromImage {
		// fps.
		FPS := fmt.Sprintf("fps=%d", fps)
		v.Chain(FPS)
	}

	// scale.
	scale := fmt.Sprintf("scale=%d:%d:force_original_aspect_ratio=decrease", w, h)

	// pad.
	pad := fmt.Sprintf("pad=%d:%d:(ow-iw)/2:(oh-ih)/2", w, h)

	// setsar.
	setsar := "setsar=1:1"

	v.Chain(scale).Chain(pad).Chain(setsar)

	// format for clip from image only.
	if c.isFromImage {
		// format.
		format := "format=yuv420p"
		v.Chain(format)
	} else {
		// trim for video clip only.
		// Get trim parameters.
		param, err := TrimParams(c.Start, c.End)
		if err != nil {
			return nil, fmt.Errorf("get trim params error: %v", err)
		}

		// trim if parameters are not empty.
		if param != "" {
			trim := "trim=" + param
			v.Chain(trim)

			// setpts.
			setpts := "setpts=PTS-STARTPTS"
			v.Chain(setpts)
		}
	}

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
	// Create the audio filter chain.
	aOutLabel := fmt.Sprintf("[a%d]", c.inputID)
	a := ffmpeg.NewFilterChain(aOutLabel)

	// aevalsrc to generate silence for clip from image or muted clip only.
	if c.isFromImage || c.IsMuted {
		// aevalsrc.
		d := int(c.Duration().Second())
		aevalsrc := fmt.Sprintf("aevalsrc=0:c=stereo:s=48000:d=%d", d)
		a.Chain(aevalsrc)
	} else {
		// Add input.
		// "X:a:0" which X is the input ID.
		a.AddInputByID(c.inputID, "a", 0)

		// atrim.
		// Get atrim parameters.
		param, err := TrimParams(c.Start, c.End)
		if err != nil {
			return nil, fmt.Errorf("get trim params error: %v", err)
		}

		// atrim if parameters are not empty.
		if param != "" {
			atrim := "atrim=" + param
			a.Chain(atrim)

			// asetpts.
			asetpts := "asetpts=PTS-STARTPTS"
			a.Chain(asetpts)
		}
	}

	return a, nil
}

func concatCommandArgs(datas []ClipData, output string, w, h, fps int, isRankMode bool, styles ...srt.Style) ([]string, error) {
	var (
		args    []string
		clips   []*Clip
		inputs  []string
		vChains []*ffmpeg.FilterChain
		aChains []*ffmpeg.FilterChain
	)

	// Create a filter graph.
	fg := ffmpeg.NewFilterGraph()

	// Create a filter chain used for concat filter.
	fc := ffmpeg.NewFilterChain("[outv]", "[outa]")

	for i, data := range datas {
		c, err := NewClip(i, data)
		if err != nil {
			return nil, fmt.Errorf("new clip: %v error: %v", data, err)
		}
		clips = append(clips, c)

		// Get ffmpeg input.
		inputs = append(inputs, c.InputArgs()...)

		// Get video filter chain.
		v, err := c.VideoFilterChain(w, h, fps, styles...)
		if err != nil {
			return nil, fmt.Errorf("get video filter chain of %v error: %v", c.ClipData, err)
		}

		vChains = append(vChains, v)
		// Add video chain's output as concat filter chain's input.
		fc.AddInput(v.Output(0))

		// Get audio filter chain.
		a, err := c.AudioFilterChain()
		if err != nil {
			return nil, fmt.Errorf("get audio filter chain of %v error: %v", c.ClipData, err)
		}

		aChains = append(aChains, a)
		// Add audio chain's output as concat filter chain's input.
		fc.AddInput(a.Output(0))
	}

	// Add ffmpeg inputs.
	args = append(args, inputs...)

	// Add concat filter.
	concat := fmt.Sprintf("concat=n=%d:v=1:a=1", len(clips))
	fc.Chain(concat)

	// Add filter graph's args.
	fg.Append(vChains...)
	fg.Append(aChains...)
	fg.Append(fc)

	args = append(args, fg.Args()...)

	// Add outputs.
	// Get filter graph's outputs(last filter chain: concat filter chain's outputs)
	outputs := fg.Outputs()

	for _, output := range outputs {
		args = append(args, "-map", output)
	}

	// Video encoder.
	args = append(args, "-c:v", "libx264")

	// Audio encoder.
	args = append(args, "-c:a", "aac")

	// Output.
	args = append(args, output)

	return args, nil
}

func ConcatCommand(datas []ClipData, output string, w, h, fps int, isRankMode bool, styles ...srt.Style) (*exec.Cmd, error) {
	args, err := concatCommandArgs(datas, output, w, h, fps, isRankMode, styles...)
	if err != nil {
		return nil, err
	}

	return exec.Command("ffmpeg", args...), nil
}
