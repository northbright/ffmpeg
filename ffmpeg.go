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

// Disposition represents the disposition data of a stream.
type Disposition struct {
	Default int `json:"default"`
}

// VideoDisposition represents the disposition data of a video stream.
type VideoDisposition struct {
	Disposition
	AttachedPic int `json:"attached_pic"`
	StillImage  int `json:"still_image"`
}

// AudioDisposition represents the disposition data of an audio stream.
type AudioDisposition struct {
	Disposition
	Dub      int `json:"dub"`
	Original int `json:"original"`
	Comment  int `json:"comment"`
	Lyrics   int `json:"lyrics"`
	Karaoke  int `json:"karaoke"`
}

// SubtitleDisposition represents the disposition data of a subtitle stream.
type SubtitleDisposition struct {
	Disposition
	Forced int `json:"forced"`
}

// Tags represents the tags of a stream.
type Tags struct {
	Duration string `json:"DURATION"`
}

// VideoTags represents the tags of a video stream.
type VideoTags struct {
	Tags
	HandlerName string `json:"HANDLER_NAME"`
	Encoder     string `json:"ENCODER"`
}

// AudioTags represents the tags of an audio stream.
type AudioTags struct {
	Tags
	HandlerName string `json:"HANDLER_NAME"`
}

// SubtitleTags represents the tags of a subtitle stream.
type SubtitleTags struct {
	Tags
	Language string `json:"language"`
	Title    string `json:"title"`
}

// VideoStream represents ffprobe's video stream.
type VideoStream struct {
	Stream
	Width            int    `json:"width"`
	Height           int    `json:"height"`
	RFrameRate       string `json:"r_frame_rate"`
	VideoDisposition `json:"disposition"`
	VideoTags        `json:"tags"`
}

// AudioStream represents ffprobe's audio stream.
type AudioStream struct {
	Stream
	ChannelsNum      int    `json:"channels"`
	ChannelsLayout   string `json:"channel_layout"`
	AudioDisposition `json:"disposition"`
	AudioTags        `json:"tags"`
}

// SubtitleStream represents ffprobe's subtitle stream.
type SubtitleStream struct {
	Stream
	SubtitleDisposition `json:"disposition"`
	SubtitleTags        `json:"tags"`
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
