package main

import (
	"fmt"
	"os"

	"github.com/RasmusLindroth/gaudio/pkg/pulseaudio"
)

func usageText() {
	fmt.Println("gtoggle <sink name> <sink name>  (without <>)")
	os.Exit(1)
}

func main() {
	if len(os.Args) < 3 {
		usageText()
	}

	sinks := pulseaudio.GetOutputSinks()

	if len(sinks) == 0 {
		fmt.Println("Could't find any sinks")
		os.Exit(1)
	}

	currentName := pulseaudio.GetCurrentOutput()
	fmt.Println(currentName)

	var sinksToUse []*pulseaudio.OutputSink
	for _, s := range sinks {
		if os.Args[1] == s.Name || os.Args[2] == s.Name {
			sinksToUse = append(sinksToUse, s)
		}
	}

	if len(sinksToUse) != 2 {
		fmt.Println("Couldn't find selected sinks")
		os.Exit(1)
	}

	if sinksToUse[0].Name == currentName {
		pulseaudio.MoveInputSinks(sinksToUse[1])
	} else {
		pulseaudio.MoveInputSinks(sinksToUse[0])
	}
}
