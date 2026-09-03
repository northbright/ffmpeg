package concat

type Clip interface {
	OriginalDuration() int
	Start() string
	End() string
	SubTitle() string
	IsMuted() bool
	FadeInDuration() int
	FadeOutDuration() int
	FiterChain() string
}

type clip struct {
	originalDuration int
	start            string
	end              string
	subTitle         string
	isMuted          bool
	fadeInDuration   int
	fadeOutDuration  int
}

type VideoClip struct {
	file string
	clip
}

func NewVideoClip(file, start, end, subTitle string, isMuted bool, fadeInDuration, fadeOutDuration int) (Clip, error) {
	return nil, nil
}

type SingleImageClip struct {
	file string
	clip
}
