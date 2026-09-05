package srt

import (
	"fmt"
	"os"
	"os/exec"
)

type SubTitle struct {
	Start string
	End   string
	Text  string
}

func WriteFile(srtFile string, subtitles []Subtitle) error {
	if len(subtitles) == 0 {
		return "", fmt.Errorf("empty subtitles")
	}

	content := ""
	for i, s := range subtitles {
		content = append(content, fmt.Sprintf("%d\n%s --> %s\n%s\n\n", i+1, s.Start, s.End, s.Text))
	}

	// Write subtitles in SRT file.
	if err := os.WriteFile(srtFile, []byte(content), 0644); err != nil {
		return err
	}

	return nil
}

func AddSoftSubtitlesCommand(videoFile string, subtitles map[string]string) (*exec.Cmd, error) {
	return nil, nil
}
