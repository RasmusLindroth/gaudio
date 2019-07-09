package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/RasmusLindroth/gaudio/pkg/pulseaudio"
)

func usageText() {
	fmt.Println("gvol up <int> - no plus sign")
	fmt.Println("gvol down <int> - no minus sign")
	os.Exit(1)
}

func main() {
	if len(os.Args) < 3 {
		usageText()
	}

	if !(os.Args[1] == "up" || os.Args[1] == "down") {
		usageText()
	}

	vol, err := strconv.Atoi(os.Args[2])
	if err != nil {
		usageText()
	}

	sinks := pulseaudio.GetOutputSinks()

	if len(sinks) == 0 {
		fmt.Println("No sinks")
		os.Exit(0)
	}

	current := pulseaudio.GetCurrentOutput()

	var sinkToUse *pulseaudio.OutputSink
	for _, s := range sinks {
		if current == s.Name {
			sinkToUse = s
		}
	}

	if sinkToUse == nil {
		fmt.Println("No default sink")
		os.Exit(0)
	}

	switch os.Args[1] {
	case "up":
		sinkToUse.IncreaseVolume(vol)
	case "down":
		sinkToUse.DecreaseVolume(vol)
	}
}
