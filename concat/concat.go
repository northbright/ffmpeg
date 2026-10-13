package concat

import (
	"os/exec"
)

type Clip struct {
	id               int
	File             string `json:"file"`
	Start            string `json:"start"`
	End              string `json:"end"`
	SubTitle         string `json:"subtitle"`
	IsMuted          bool   `json:"is_muted,omitempty"`
	FadeInDuration   int    `json:"fade_in_duration,omitempty"`
	FadeOutDuration  int    `json:"fade_out_duration,omitempty"`
	AFadeInDuration  int    `json:"afade_in_duration,omitempty"`
	AFadeOutDuration int    `json:"afade_out_duration,omitempty"`
}

func ConcatCommand(clips []Clip, output string, w, h, fps int) (*exec.Cmd, error) {
	/*
		inputs := append(inputs, clips...)

		for id, input := range clips {
			clips[i].id = i
		}

		return nil, nil
	*/
}

func (c *Clip) FilterChain(w, h, fps int) (string, error) {
	/*
		str := ""
		// Video filter chain.
		// Set input.
		str += fmt.Sprintf("[%d:v]", c.id)

		// Scale.
		scale := fmt.Sprintf("scale=%d:%d:force_original_aspect_ratio=decrease", w, h)

		// Pad.
		pad := fmt.Sprintf("pad=%d:%d:(ow-iw)/2:(oh-ih)/2", w, h)

		// Setsar.
		setsar := "setsar=1:1"

		// Trim.
		if c.Start != "" || c.End != "" {
		}
	*/
}

/*
func NewSingleImageClip(file string, duration float64, subTitle string, options ...ClipOption) (*Clip, error) {
	return nil, nil
}
*/
