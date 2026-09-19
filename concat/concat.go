package concat

type Clip struct {
	Start           string
	End             string
	SubTitle        string
	IsMuted         bool
	FadeInDuration  int
	FadeOutDuration int
}

func NewVideoClip(file, start, end, subTitle string, isMuted bool, fadeInDuration, fadeOutDuration int) (Clip, error) {
	return nil, nil
}

type SingleImageClip struct {
	file string
	clip
}
