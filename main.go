package main

import (
	"fmt"
	"sort"

	"github.com/RasmusLindroth/gaudio/pkg/pulseaudio"
)

func main() {
	sinks := pulseaudio.GetOutputSinks()
	sort.Sort(pulseaudio.SortByUsage(sinks))
	for _, s := range sinks {
		fmt.Println(s)
	}
}
