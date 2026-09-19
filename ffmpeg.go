package ffmpeg

import (
	"encoding/json"
	"fmt"
	"os/exec"
)

// Format represents ffprobe's format.
type Format struct {
	Filename   string            `json:"filename"`
	StreamsNum int               `json:"nb_streams"`
	Size       string            `json:"size"`
	Duration   string            `json:"duration,omitempty"`
	Tags       map[string]string `json:"tags,omitempty"`
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

// SubtitleDisposition represents the disposition data of a subtitle stream.
type SubtitleDisposition struct {
	Disposition
	Forced int `json:"forced"`
}

// SubtitleTags represents the tags of a subtitle stream.
type SubtitleTags struct {
	Language string `json:"language"`
	Title    string `json:"title"`
}

// VideoStream represents ffprobe's video stream.
type VideoStream struct {
	Stream
	Width       int    `json:"width"`
	Height      int    `json:"height"`
	RFrameRate  string `json:"r_frame_rate"`
	Duration    string `json:"duration,omitempty"` // not exist in MKV.
	Disposition `json:"disposition"`
}

// AudioStream represents ffprobe's audio stream.
type AudioStream struct {
	Stream
	ChannelsNum    int    `json:"channels"`
	ChannelsLayout string `json:"channel_layout"`
	Duration       string `json:"duration,omitempty"` // not exist in MKV.
	Disposition    `json:"disposition"`
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

// MajorBrand returns "major_brand" by running ffprobe.
// "major_brand" may not exist in format tags.
func MajorBrand(file string) (string, error) {
	// Get format.
	f, err := GetFormat(file)
	if err != nil {
		return "", fmt.Errorf("get format info error: %v", err)
	}

	if len(f.Tags) == 0 {
		return "", nil
	}

	majorBrand, ok := f.Tags["major_brand"]
	if !ok {
		return "", nil
	}

	return majorBrand, nil
}

// IsImage returns if a file is image by running ffprobe.
func IsImage(file string) (bool, error) {
	// Get v:0 stream.
	streams, err := GetVideoStreams(file)
	if err != nil {
		return false, fmt.Errorf("get video streams error: %v", err)
	}

	if len(streams) == 0 {
		return false, fmt.Errorf("no video stream")
	}

	switch streams[0].CodecName {
	case "mjpeg", "png", "bmp", "tiff", "webp", "jpegxl":
		return true, nil

	case "av1":
		majorBrand, err := MajorBrand(file)
		if err != nil {
			return false, fmt.Errorf("get major_brand error: %v", err)
		}
		return majorBrand == "avif", nil

	case "hevc":
		majorBrand, err := MajorBrand(file)
		if err != nil {
			return false, fmt.Errorf("get major_brand error: %v", err)
		}
		return majorBrand == "heic" || majorBrand == "heif" || majorBrand == "mif1", nil

	default:
		return false, nil
	}
}
