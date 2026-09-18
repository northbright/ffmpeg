package srt

import (
	"context"
	"fmt"
	"maps"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strconv"
	"strings"

	"github.com/northbright/ffmpeg"
)

// Subtitle represents a subtitle in a SRT file.
type Subtitle struct {
	Start string
	End   string
	Text  string
}

// WriteFile writes subtitles to a SRT file.
func WriteFile(srtFile string, subtitles []Subtitle) error {
	if len(subtitles) == 0 {
		return fmt.Errorf("empty subtitles")
	}

	content := ""
	for i, s := range subtitles {
		content += fmt.Sprintf("%d\n%s --> %s\n%s\n\n", i+1, s.Start, s.End, s.Text)
	}

	// Write subtitles in SRT file.
	if err := os.WriteFile(srtFile, []byte(content), 0644); err != nil {
		return err
	}

	return nil
}

// addSoftSubArgs returns the arguments of [os/exec.Cmd] to add a soft subtitle track to a video with ffmpeg.
// input: input video.
// srtFile: srt file.
// lang: three-letter [ISO 639-2 Code](e.g. "eng", "spa", "chi").
// title: a user-friendly name for the subtitle track selection menu(e.g. "English", "Spanish", "Chinese").
// isDefault: if set the subtitle stream as default.
// Video players show default subtitle stream automatically.
// output: output file. Supported containers: MKV, MP4, MOV.
// overwrite: if overwrite if output exists.
// [ISO 639-2 Code]: https://www.loc.gov/standards/iso639-2/php/code_list.php
func addSoftSubArgs(input, srtFile, lang, title string, isDefault bool, output string, overwrite bool) ([]string, error) {
	var args []string

	if overwrite {
		args = append(args, "-y")
	}

	args = append(args, "-i", input, "-i", srtFile)
	args = append(args, "-map", "0", "-map", "1")
	args = append(args, "-c", "copy")

	// Check if container is supported.
	ext := strings.ToLower(filepath.Ext(output))
	if ext != ".mkv" && ext != ".mp4" && ext != ".mov" {
		return nil, fmt.Errorf("unsupported container: %s", ext)
	}

	// Unlike MKV, the MP4 container does not natively support raw SRT text format.
	// To add subtitles to an MP4 container without re-encoding the video or audio,
	// convert the SRT subtitle stream to the MP4-compatible Timed Text format(mov_text).
	// See https://salivity.github.io/ffmpeg/article/how-to-add-srt-to-video-without-encoding-using-ffmpeg
	if ext == ".mp4" {
		args = append(args, "-c:s", "mov_text")
	}

	// Get subtitle stream count.
	streams, err := ffmpeg.GetSubtitleStreams(input)
	if err != nil {
		return nil, fmt.Errorf("ffmpeg.GetSubtitleStreams(%s) error: %v", input, err)
	}
	n := len(streams)

	args = append(args, fmt.Sprintf("-metadata:s:s:%d", n), fmt.Sprintf("language=%s", lang))
	args = append(args, fmt.Sprintf("-metadata:s:s:%d", n), fmt.Sprintf("title=%s", title))

	if isDefault {
		args = append(args, fmt.Sprintf("-disposition:s:%d", n), "default")
	}

	args = append(args, output)

	return args, nil
}

// AddSoftSubCommand returns the [os/exec.Cmd] to add a soft subtitle track to a video with ffmpeg.
// input: input video.
// srtFile: srt file.
// lang: three-letter [ISO 639-2 Code](e.g. "eng", "spa", "chi").
// title: a user-friendly name for the subtitle track selection menu(e.g. "English", "Spanish", "Chinese").
// isDefault: if set the subtitle stream as default.
// Video players show default subtitle stream automatically.
// output: output file. Supported containers: MKV, MP4, MOV.
// overwrite: if overwrite if output exists.
func AddSoftSubCommand(input, srtFile, lang, title string, isDefault bool, output string, overwrite bool) (*exec.Cmd, error) {
	args, err := addSoftSubArgs(input, srtFile, lang, title, isDefault, output, overwrite)
	if err != nil {
		return nil, err
	}

	return exec.Command("ffmpeg", args...), nil
}

// AddSoftSubCommandContext is the context version of [AddSoftSubCommand].
func AddSoftSubCommandContext(ctx context.Context, input, srtFile, lang, title string, isDefault bool, output string, overwrite bool) (*exec.Cmd, error) {
	args, err := addSoftSubArgs(input, srtFile, lang, title, isDefault, output, overwrite)
	if err != nil {
		return nil, err
	}

	return exec.CommandContext(ctx, "ffmpeg", args...), nil
}

// AddSoftSub adds a soft subtitle track to a video with ffmpeg.
// It returns the output from ffmpeg command.
// input: input video.
// srtFile: srt file.
// lang: three-letter [ISO 639-2 Code](e.g. "eng", "spa", "chi").
// title: a user-friendly name for the subtitle track selection menu(e.g. "English", "Spanish", "Chinese").
// isDefault: if set the subtitle stream as default.
// Video players show default subtitle stream automatically.
// output: output file. Supported containers: MKV, MP4, MOV.
// overwrite: if overwrite if output exists.
func AddSoftSub(ctx context.Context, input, srtFile, lang, title string, isDefault bool, output string, overwrite bool) (string, error) {
	var (
		err error
		cmd *exec.Cmd
	)

	if ctx == nil {
		cmd, err = AddSoftSubCommand(input, srtFile, lang, title, isDefault, output, overwrite)
	} else {
		cmd, err = AddSoftSubCommandContext(ctx, input, srtFile, lang, title, isDefault, output, overwrite)
	}

	if err != nil {
		return "", fmt.Errorf("Generate ffmpeg command to add soft subtitle error: %v", err)
	}

	out, err := cmd.CombinedOutput()
	return string(out), err
}

// HardSub represents the hard-coding subtitles.
type HardSub struct {
	srtFile string
	styles  map[string]string
}

// Style represents the ASS style.
type Style func(hs *HardSub)

// FontName returns the font name style.
// Default: "Arial".
func FontName(name string) Style {
	return func(hs *HardSub) {
		hs.styles["FontName"] = name
	}
}

// FontSize returns the font size style.
// Default: 16.
func FontSize(size int) Style {
	return func(hs *HardSub) {
		hs.styles["FontSize"] = strconv.Itoa(size)
	}
}

// PrimaryColour returns the primary color style.
// It'll generate the color in "&HAABBGGRR" format for ffmpeg automatically.
// Default: white.
func PrimaryColour(R, G, B, A uint8) Style {
	return func(hs *HardSub) {
		hs.styles["PrimaryColour"] = fmt.Sprintf("&H%02x%02x%02x%02x", A, B, G, R)
	}
}

// OutlineColour returns the outline color style.
// It'll generate the color in "&HAABBGGRR" format for ffmpeg automatically.
// Default: black.
func OutlineColour(R, G, B, A uint8) Style {
	return func(hs *HardSub) {
		hs.styles["OutlineColour"] = fmt.Sprintf("&H%02x%02x%02x%02x", A, B, G, R)
	}
}

// BackColour returns the back color style.
// It'll generate the color in "&HAABBGGRR" format for ffmpeg automatically.
// Default: black.
func BackColour(R, G, B, A uint8) Style {
	return func(hs *HardSub) {
		hs.styles["BackColour"] = fmt.Sprintf("&H%02x%02x%02x%02x", A, B, G, R)
	}
}

// Outline returns the outline thickness style.
// Default: 2.
func Outline(thickness uint8) Style {
	return func(hs *HardSub) {
		hs.styles["Outline"] = fmt.Sprintf("%d", thickness)
	}
}

// Shadow returns the shadow depth style.
// Default: 0.
func Shadow(depth uint8) Style {
	return func(hs *HardSub) {
		hs.styles["Shadow"] = fmt.Sprintf("%d", depth)
	}
}

// Bold returns the bold style.
// Default: false.
func Bold(bold bool) Style {
	return func(hs *HardSub) {
		if bold {
			hs.styles["Bold"] = "1"
		} else {
			hs.styles["Bold"] = "0"
		}
	}
}

// Italic returns the Italic style.
// Default: false.
func Italic(bold bool) Style {
	return func(hs *HardSub) {
		if bold {
			hs.styles["Italic"] = "1"
		} else {
			hs.styles["Italic"] = "0"
		}
	}
}

// Alignment return the subtitle alignment style.
// Range: 1 - 9(num keyboard layout).
// Default: 2.
func Alignment(align uint8) Style {
	return func(hs *HardSub) {
		hs.styles["Alignment"] = fmt.Sprintf("%d", align)
	}
}

// MarginV return the vertical margin of subtitle style.
// Default: 10.
func MarginV(margin uint8) Style {
	return func(hs *HardSub) {
		hs.styles["MarginV"] = fmt.Sprintf("%d", margin)
	}
}

// NewHardSub returns a new HardSub.
// srtFile: srt file.
// styles: slice of [Style].
func NewHardSub(srtFile string, styles ...Style) *HardSub {
	hs := &HardSub{srtFile: srtFile, styles: map[string]string{}}

	for _, style := range styles {
		style(hs)
	}

	return hs
}

// VideoFilter returns the video filter string used for the hard-coding subtitles.
func (hs *HardSub) VideoFilter() string {
	vf := fmt.Sprintf("subtitles='%s'", hs.srtFile)
	if len(hs.styles) == 0 {
		return vf
	}

	keys := slices.Sorted(maps.Keys(hs.styles))

	vf += ":force_style='"
	first := true
	for _, k := range keys {
		if !first {
			vf += ","
		}

		vf += fmt.Sprintf("%s=%s", k, hs.styles[k])
		first = false
	}
	vf += "'"

	return vf
}

/*
// AddHardSub adds hard coding subtitle to a video with ffmpeg.
// It returns the output from ffmpeg command.
// input: input video.
// srtFile: srt file.
// params:
// output: output file.
// overwrite: if overwrite if output exists.
func AddHardSub(ctx context.Context, input, srtFile string, params map[string]string, overwrite bool) (string, error) {
	var (
		err error
		cmd *exec.Cmd
	)

	if ctx == nil {
		cmd, err = AddSoftSubCommand(input, srtFile, lang, title, isDefault, output, overwrite)
	} else {
		cmd, err = AddSoftSubCommandContext(ctx, input, srtFile, lang, title, isDefault, output, overwrite)
	}

	if err != nil {
		return "", fmt.Errorf("Generate ffmpeg command to add soft subtitle error: %v", err)
	}

	out, err := cmd.CombinedOutput()
	return string(out), err
}
*/
