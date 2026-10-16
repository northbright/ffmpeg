package ffmpeg

type FilterGraph struct {
	chains []*FilterChain
}

func NewFilterGraph() *FilterGraph {
	var chains []*FilterChain
	return &FilterGraph{chains: chains}
}

func (fg *FilterGraph) Append(chains ...*FilterChain) {
	fg.chains = append(fg.chains, chains...)
}

func (fg *FilterGraph) Outputs() []string {
	l := len(fg.chains)

	if l == 0 {
		return nil
	}

	// Use last chain's output as filter graph's output.
	chain := fg.chains[l-1]

	return chain.Outputs()
}

func (fg *FilterGraph) Args() []string {
	var args []string

	l := len(fg.chains)
	if l == 0 {
		return nil
	}

	args = append(args, "-filter_complex")

	for i, fc := range fg.chains {
		str := fc.String()
		if i != l-1 {
			str += ";\n"
		}
		args = append(args, str)
	}

	return args
}
