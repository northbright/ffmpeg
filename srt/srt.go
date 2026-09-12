package srt

import (
	"context"
	"fmt"
	"os"
	"os/exec"

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
// output: output file. The input and output's container(format) should be the same.
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
// output: output file. The input and output's container(format) should be the same.
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
// output: output file. The input and output's container(format) should be the same.
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
