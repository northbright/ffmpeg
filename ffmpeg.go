package ffmpeg

import (
	"encoding/json"
	"fmt"
	"os/exec"
)

// Format represents ffprobe's format.
type Format struct {
	Filename   string `json:"filename"`
	StreamsNum int    `json:"nb_streams"`
	Size       string `json:"size"`
	Duration   string `json:"duration"`
}

// Stream represents ffprobe's stream.
type Stream struct {
	Index     int    `json:"index"`
	CodecName string `json:"codec_name"`
	CodecType string `json:"codec_type"`
}

// VideoStream represents ffprobe's video stream.
type VideoStream struct {
	Stream
	Width      int    `json:"width"`
	Height     int    `json:"height"`
	RFrameRate string `json:"r_frame_rate"`
}

// Tags represents ffprobe's tags.
type Tags struct {
	Language string `json:"language"`
}

// AudioStream represents ffprobe's audio stream.
type AudioStream struct {
	Stream
	ChannelsNum    int    `json:"channels"`
	ChannelsLayout string `json:"channel_layout"`
	Tags           `json:"tags"`
}

// SubtitleStream represents ffprobe's subtitle stream.
type SubtitleStream struct {
	Stream
	Tags `json:"tags"`
}

// GetFormat returns the format by running ffprobe on the media file.
func GetFormat(file string) (*Format, error) {
	var args []string

	args = append(args, "-v", "error", "-show_entries", "format", "-of", "json", file)
	cmd := exec.Command("ffprobe", args...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("ffprobe on %s error: %v", file, err)
	}

	type _data struct {
		Format `json:"format"`
	}

	var d _data

	if err := json.Unmarshal(output, &d); err != nil {
		return nil, fmt.Errorf("json.Unmarshal on ffprobe output(%s) error: %s", string(output), err)
	}

	return &(d.Format), nil
}

// GetVideoStreams returns the video streams by running ffprobe on the media file.
func GetVideoStreams(file string) ([]VideoStream, error) {
	var args []string

	args = append(args, "-v", "error", "-select_streams", "v", "-show_entries", "stream", "-of", "json", file)
	cmd := exec.Command("ffprobe", args...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("ffprobe on %s error: %v", file, err)
	}

	type _data struct {
		Streams []VideoStream `json:"streams"`
	}

	var d _data

	if err := json.Unmarshal(output, &d); err != nil {
		return nil, fmt.Errorf("json.Unmarshal on ffprobe output(%s) error: %s", string(output), err)
	}

	return d.Streams, nil
}

// GetAudioStreams returns the audio streams by running ffprobe on the media file.
func GetAudioStreams(file string) ([]AudioStream, error) {
	var args []string

	args = append(args, "-v", "error", "-select_streams", "a", "-show_entries", "stream", "-of", "json", file)
	cmd := exec.Command("ffprobe", args...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("ffprobe on %s error: %v", file, err)
	}

	type _data struct {
		Streams []AudioStream `json:"streams"`
	}

	var d _data

	if err := json.Unmarshal(output, &d); err != nil {
		return nil, fmt.Errorf("json.Unmarshal on ffprobe output(%s) error: %s", string(output), err)
	}

	return d.Streams, nil
}

// GetSubtitleStreams returns the subtitle streams by running ffprobe on the media file.
func GetSubtitleStreams(file string) ([]SubtitleStream, error) {
	var args []string

	args = append(args, "-v", "error", "-select_streams", "s", "-show_entries", "stream", "-of", "json", file)
	cmd := exec.Command("ffprobe", args...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("ffprobe on %s error: %v", file, err)
	}

	type _data struct {
		Streams []SubtitleStream `json:"streams"`
	}

	var d _data

	if err := json.Unmarshal(output, &d); err != nil {
		return nil, fmt.Errorf("json.Unmarshal on ffprobe output(%s) error: %s", string(output), err)
	}

	return d.Streams, nil
}
