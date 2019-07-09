package pulseaudio

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/RasmusLindroth/gaudio/pkg/command"
)

func GetCurrentOutput() string {
	out, err := command.RunWithOutput("pactl", "info")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	for _, l := range strings.Split(out, "\n") {
		parts := splitByDelimiter(l, ":")

		if len(parts) < 2 {
			continue
		}

		if parts[0] == "Default Sink" {
			return parts[1]
		}
	}

	return ""
}

type OutputSink struct {
	Index    string
	Name     string
	State    string
	Linked   int
	Priority int
	Volume   int
}

func GetOutputSinks() []*OutputSink {
	out, err := command.RunWithOutput("pacmd", "list-sinks")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	return parseOutputSinks(out)
}

func parseOutputSink(lines []string) *OutputSink {
	sink := &OutputSink{}

	for _, l := range lines {
		colon := splitByDelimiter(l, ":")
		if len(colon) > 1 {
			switch colon[0] {
			case "index":
				sink.Index = colon[1]
			case "* index":
				sink.Index = colon[1]
			case "name":
				sink.Name = strings.Trim(colon[1], "<>")
			case "state":
				sink.State = colon[1]
			case "priority":
				priority := colon[1]
				if i, err := strconv.Atoi(priority); err == nil {
					sink.Priority = i
				}
			case "linked by":
				linked := colon[1]
				if i, err := strconv.Atoi(linked); err == nil {
					sink.Linked = i
				}
			case "volume":
				if len(colon) < 3 {
					continue
				}
				volParts := splitByDelimiter(colon[2], "/")
				vol, err := strconv.Atoi(strings.Trim(volParts[1], "%"))
				if err != nil {
					continue
				}
				sink.Volume = vol
			}
		}
	}

	return sink
}
func parseOutputSinks(raw string) []*OutputSink {
	var sinks []*OutputSink

	var lines []string
	for _, l := range strings.Split(raw, "\n") {
		lines = append(lines, strings.TrimSpace(l))
	}

	inSink := false
	var currSink []string
	for _, l := range lines {
		parts := strings.Split(l, ":")
		if len(parts) > 1 && parts[0] == "index" || parts[0] == "* index" {
			if inSink {
				sinks = append(sinks, parseOutputSink(currSink))
			}

			inSink = true
			currSink = []string{}
		}

		if inSink {
			currSink = append(currSink, l)
		}
	}

	if inSink {
		sinks = append(sinks, parseOutputSink(currSink))
	}

	return sinks
}

func GetInputSinks() []*InputSink {
	out, err := command.RunWithOutput("pacmd", "list-sink-inputs")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	return parseInputSinks(out)
}

func SetDefaultSink(sink *OutputSink) {
	command.RunCmd("pacmd", "set-default-sink", sink.Index)
}

func MoveInputSinks(sink *OutputSink) {
	inputs := GetInputSinks()

	SetDefaultSink(sink)

	for _, i := range inputs {
		MoveInputSink(i, sink.Index)
	}
}

func MoveInputSink(sink *InputSink, dest string) {
	command.RunCmd("pacmd", "move-sink-input", sink.Index, dest)
}

func (sink *OutputSink) IncreaseVolume(vol int) {
	ChangeVolume(sink.Index, sink.Volume+vol)
}

func (sink *OutputSink) DecreaseVolume(vol int) {
	ChangeVolume(sink.Index, sink.Volume-vol)
}

func ChangeVolume(dest string, vol int) {
	if vol < 0 {
		vol = 0
	} else if vol > 100 {
		vol = 100
	}
	volume := fmt.Sprintf("%d%%", vol)
	command.RunCmd("pactl", "set-sink-volume", dest, volume)
}

type InputSink struct {
	Index   string
	Sink    string
	SinkID  string
	AppName string
}

func splitByDelimiter(l string, del string) []string {
	parts := strings.Split(l, del)
	var r []string

	for _, p := range parts {
		r = append(r, strings.TrimSpace(p))
	}

	return r
}

func parseInputSink(lines []string) *InputSink {
	sink := &InputSink{}

	for _, l := range lines {
		colon := splitByDelimiter(l, ":")
		if len(colon) > 1 {
			switch colon[0] {
			case "index":
				sink.Index = colon[1]
			case "sink":
				sink.Sink = colon[1]

				sinkParts := strings.Split(sink.Sink, " ")
				if len(sinkParts) > 0 {
					sink.SinkID = sinkParts[0]
				}
			}
		}

		equal := splitByDelimiter(l, "=")
		if len(equal) > 1 {
			switch equal[0] {
			case "application.name":
				sink.AppName = strings.Trim(equal[1], "\"")
			}
		}
	}

	return sink
}

func parseInputSinks(raw string) []*InputSink {
	var sinks []*InputSink

	var lines []string
	for _, l := range strings.Split(raw, "\n") {
		lines = append(lines, strings.TrimSpace(l))
	}

	inSink := false
	var currSink []string
	for _, l := range lines {
		parts := strings.Split(l, ":")
		if len(parts) > 1 && parts[0] == "index" {
			if inSink {
				sinks = append(sinks, parseInputSink(currSink))
			}

			inSink = true
			currSink = []string{}
		}

		if inSink {
			currSink = append(currSink, l)
		}
	}

	if inSink {
		sinks = append(sinks, parseInputSink(currSink))
	}

	return sinks
}
