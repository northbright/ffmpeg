package srt

import (
	"context"
	"fmt"
	"os"
	"os/exec"
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

// addSoftSubtitleArgs returns the arguments of [os/exec.Cmd] to add a soft subtitle track to a video with ffmpeg.
// input: input video.
// srtFile: srt file.
// lang: three-letter [ISO 639-2 Code](e.g. "eng", "spa", "chi").
// title: a user-friendly name for the subtitle track selection menu(e.g. "English", "Spanish", "Chinese").
// output: output file.
// overwrite: if overwrite if output exists.
// [ISO 639-2 Code]: https://www.loc.gov/standards/iso639-2/php/code_list.php
func addSoftSubtitleArgs(input, srtFile, lang, title, output string, overwrite bool) []string {
	var args []string

	if overwrite {
		args = append(args, "-y")
	}

	args = append(args, "-i", input, "-i", srtFile)
	args = append(args, "-c", "copy")
	args = append(args, "-metadata:s:s:0", fmt.Sprintf("language=%s", lang))
	args = append(args, "-metadata:s:s:0", fmt.Sprintf("title=%s", title))
	args = append(args, output)

	return args
}

// AddSoftSubtitleCommand returns the [os/exec.Cmd] to add a soft subtitle track to a video with ffmpeg.
// input: input video.
// srtFile: srt file.
// lang: three-letter [ISO 639-2 Code](e.g. "eng", "spa", "chi").
// title: a user-friendly name for the subtitle track selection menu(e.g. "English", "Spanish", "Chinese").
// output: output file.
// overwrite: if overwrite if output exists.
func AddSoftSubtitleCommand(input, srtFile, lang, title, output string, overwrite bool) *exec.Cmd {
	args := addSoftSubtitleArgs(input, srtFile, lang, title, output, overwrite)
	return exec.Command("ffmpeg", args...)
}

// AddSoftSubtitleCommandContext is the context version of [AddSoftSubtitleCommand].
func AddSoftSubtitleCommandContext(ctx context.Context, input, srtFile, lang, title, output string, overwrite bool) *exec.Cmd {
	args := addSoftSubtitleArgs(input, srtFile, lang, title, output, overwrite)
	return exec.Command("ffmpeg", args...)
}

// AddSoftSubtitle adds a soft subtitle track to a video with ffmpeg.
// It returns the output from ffmpeg command.
// input: input video.
// srtFile: srt file.
// lang: three-letter [ISO 639-2 Code](e.g. "eng", "spa", "chi").
// title: a user-friendly name for the subtitle track selection menu(e.g. "English", "Spanish", "Chinese").
// output: output file.
// overwrite: if overwrite if output exists.
func AddSoftSubtitle(ctx context.Context, input, srtFile, lang, title, output string, overwrite bool) (string, error) {
	var cmd *exec.Cmd

	if ctx == nil {
		cmd = AddSoftSubtitleCommand(input, srtFile, lang, title, output, overwrite)
	} else {
		cmd = AddSoftSubtitleCommandContext(ctx, input, srtFile, lang, title, output, overwrite)
	}

	out, err := cmd.CombinedOutput()
	return string(out), err
}
