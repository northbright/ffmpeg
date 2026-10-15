package ffmpeg

import "fmt"

// FilterChain represents the filterchain of ffmpeg.
type FilterChain struct {
	inputs  []any
	outputs []string
	filters []string
}

// filterChainOutputData stores the filterchan and the output ID to generate output label as another filterchain's input.
type filterChainOutputData struct {
	fc *FilterChain
	id int
}

// NewFilterChain returns a filterchain by pre-defined outputs(labels) in the "[OUTPUT_LABEL]" format.
// e.g. concatFilterChain := ffcmd.NewFilterChain("[out_v]", "[out_a]")
func NewFilterChain(outputs ...string) *FilterChain {
	return &FilterChain{inputs: []any{}, outputs: outputs, filters: []string{}}
}

// AddInput adds raw string(label) as the input.
func (fc *FilterChain) AddInput(input string) {
	fc.inputs = append(fc.inputs, input)
}

// AddInputByID adds input by specifying input ID, stream type and ID in the stream(e.g. "[0:v:0]", [0:a:1]").
// inputID: 0-based input ID.
// streamType: "v" for video, "a" for audio, "s" for subtitles, "d" for data. See ffmpeg's doc for more.
// streamID: index of the stream for the stream type.
func (fc *FilterChain) AddInputByID(inputID int, streamType string, streamID int) {
	fc.inputs = append(fc.inputs, fmt.Sprintf("[%d:%s:%d]", inputID, streamType, streamID))
}

// AddInputByOutput adds another filterchain's output as input.
// It's useful when another filterchain's output is generated dynamically.
func (fc *FilterChain) AddInputByOutput(fcOut *FilterChain, outputID int) {
	fc.inputs = append(fc.inputs, &filterChainOutputData{fcOut, outputID})
}

// Input returns the input string by 0-based index.
func (fc *FilterChain) Input(id int) string {
	if id < 0 || id >= len(fc.inputs) {
		return ""
	}

	switch vv := fc.inputs[id].(type) {
	case string:
		return vv
	case *filterChainOutputData:
		fc := vv.fc
		id := vv.id
		return fc.Output(id)
	default:
		return ""
	}
}

// Inputs returns all inputs.
func (fc *FilterChain) Inputs() []string {
	var inputs []string

	for _, in := range fc.inputs {
		switch vv := in.(type) {
		case string:
			inputs = append(inputs, vv)
		case *filterChainOutputData:
			fc := vv.fc
			id := vv.id
			input := fc.Output(id)
			inputs = append(inputs, input)
		default:
		}
	}
	return inputs
}

// Output returns the output by 0-based index.
func (fc *FilterChain) Output(id int) string {
	if len(fc.filters) == 0 {
		return fc.Input(id)
	} else {
		if id < 0 || id >= len(fc.outputs) {
			return ""
		}
		return fc.outputs[id]
	}
}

// Outputs returns all outputs.
func (fc *FilterChain) Outputs() []string {
	if len(fc.filters) == 0 {
		return fc.Inputs()
	} else {
		return fc.outputs
	}
}

// Chain chains filter and returns a filterchain to chain next filter(e.g. fc.Chain("fps=30").Chain("scale=1280:720"))
func (fc *FilterChain) Chain(filter string) *FilterChain {
	if filter != "" {
		fc.filters = append(fc.filters, filter)
	}
	return fc
}

// String returns the filterchain's string for ffmpeg command.
func (fc *FilterChain) String() string {
	l := len(fc.filters)

	if l == 0 {
		// No filter in the chain, just return empty string as do nothing in the chain.
		return ""
	}

	str := ""
	for _, input := range fc.Inputs() {
		str += input
	}

	for i, filter := range fc.filters {
		str += filter
		if i < l-1 {
			str += ","
		}
	}

	for _, output := range fc.Outputs() {
		str += output
	}
	return str
}
