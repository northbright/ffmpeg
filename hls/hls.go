package hls

import (
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
	"strconv"
)

const (
	// Default c:v for HLS segment.
	DefaultCV string = "libx264"
	// Default c:a for HLS segment.
	DefaultCA string = "aac"
)

// segmentArgs returns arguments slice for [os/exec.Cmd] to segment HLS.
// src: HLS video path.
// duration: Segment duration in seconds(e.g. 4).
// outDir: output dir to store "playlist.m3u8" and "%05d.ts".
func segmentArgs(src string, duration int, outDir string) []string {
	var args []string

	args = append(args, "-i", src)
	args = append(args, "-force_key_frames", fmt.Sprintf("expr:gte(t,n_forced*%d)", duration))
	args = append(args, "-c:v", DefaultCV)
	args = append(args, "-c:a", DefaultCA)
	args = append(args, "-f", "hls")
	args = append(args, "-hls_time", strconv.Itoa(duration))
	args = append(args, "-hls_segment_filename", filepath.Join(outDir, "%05d.ts"))
	args = append(args, filepath.Join(outDir, "playlist.m3u8"))

	return args
}

// SegmentCommand returns the [os/exec.Cmd] to segment HLS with ffmpeg.
// src: HLS video path.
// duration: Segment duration in seconds(e.g. 4).
// outDir: output dir to store "playlist.m3u8" and "%05d.ts".
func SegmentCommand(src string, duration int, outDir string) *exec.Cmd {
	args := segmentArgs(src, duration, outDir)
	return exec.Command("ffmpeg", args...)
}

// SegmentCommandContext is the context version of [SegmentCommand].
func SegmentCommandContext(ctx context.Context, src string, duration int, outDir string) *exec.Cmd {
	args := segmentArgs(src, duration, outDir)
	return exec.CommandContext(ctx, "ffmpeg", args...)
}

// SegmentCommand segments HLS with ffmpeg.
// ctx: context used to interrupt the process.
// src: HLS video path.
// duration: Segment duration in seconds(e.g. 4).
// outDir: output dir to store "playlist.m3u8" and "%05d.ts".
func Segment(ctx context.Context, src string, duration int, outDir string) (string, error) {
	var cmd *exec.Cmd

	if ctx == nil {
		cmd = SegmentCommand(src, duration, outDir)
	} else {
		cmd = SegmentCommandContext(ctx, src, duration, outDir)
	}

	output, err := cmd.CombinedOutput()
	return string(output), err
}

// concatArgs returns arguments slice for [os/exec.Cmd] to concat HLS(.ts files).
// tsFiles: .ts files.
// dst: output HLS.
// overwrite: if overwrite if output HLS exists.
func concatArgs(tsFiles []string, dst string, overwrite bool) []string {
	var args []string

	if overwrite {
		args = append(args, "-y")
	}

	args = append(args, "-i")

	concatArg := `concat:`
	l := len(tsFiles)

	for i, ts := range tsFiles {
		if i != l-1 {
			concatArg += fmt.Sprintf(`%s|`, ts)
		} else {
			concatArg += ts
		}
	}

	args = append(args, concatArg)
	// The `-bsf:a aac_adtstoasc` flag converts the ADTS-formatted AAC headers that MPEG-TS uses into the MPEG-4 format that MP4 containers expect. Skip it and the audio won't play.
	args = append(args, "-c", "copy", "-bsf:a", "aac_adtstoasc")
	args = append(args, dst)

	return args
}

// ConcatCommand returns the [os/exec.Cmd] to concat HLS with ffmpeg.
// tsFiles: .ts files.
// dst: output HLS.
// overwrite: if overwrite if output HLS exists.
func ConcatCommand(tsFiles []string, dst string, overwrite bool) *exec.Cmd {
	args := concatArgs(tsFiles, dst, overwrite)
	return exec.Command("ffmpeg", args...)
}

// ConcatCommandContext is the context version of [ConcatCommand].
func ConcatCommandContext(ctx context.Context, tsFiles []string, dst string, overwrite bool) *exec.Cmd {
	args := concatArgs(tsFiles, dst, overwrite)
	return exec.CommandContext(ctx, "ffmpeg", args...)
}

// Concat concats HLS with ffmpeg.
// ctx: context used to interrupt the process.
// tsFiles: .ts files.
// dst: output HLS.
// overwrite: if overwrite if output HLS exists.
func Concat(ctx context.Context, tsFiles []string, dst string, overwrite bool) (string, error) {
	var cmd *exec.Cmd

	if ctx == nil {
		cmd = ConcatCommand(tsFiles, dst, overwrite)
	} else {
		cmd = ConcatCommandContext(ctx, tsFiles, dst, overwrite)
	}

	output, err := cmd.CombinedOutput()
	return string(output), err
}
