package main

import (
	"fmt"
	"strings"
)

const (
	indicatorRune  = '#'
	progressString = "\rProgress [%-10s] %3d%%"
)

var indicators [10]string

func init() {
	var buffer strings.Builder

	for i := 0; i < len(indicators); i++ {
		buffer.WriteRune(indicatorRune)
		indicators[i] = buffer.String()
	}
}

func printProgress(actual, total int64, last bool) {
	prc := int((float64(actual) / float64(total)) * 100)

	var prcLineIndicator string

	if prc >= 10 {
		idx := prc / len(indicators)
		prcLineIndicator = indicators[idx-1]
	}

	fmt.Printf(progressString, prcLineIndicator, prc)

	if last {
		fmt.Print("\n")
	}
}
