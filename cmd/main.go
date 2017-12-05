package main

import (
	"../pkg/mib"
	"fmt"
	"bufio"
	"os"
)

const RFC1213MIB  = "RFC1213-MIB"

func main() {
	var t mib.Tree
	if mib.Parse(&t, RFC1213MIB) {
		fmt.Print(t)
		mainLoop(&t)
	}
}

func mainLoop(t *mib.Tree) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Print(r)
			mainLoop(t)
		}
	}()

	reader := bufio.NewReader(os.Stdin)
	for {
		if line, err := reader.ReadString('\n'); err == nil {
			fmt.Print(t.SubtreeString(mib.ShortOid(line)))
		}
	}
}
