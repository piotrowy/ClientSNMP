package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"./snmp/ber"
	"./snmp/mib"
)

const (
	RFC1213MIB = "RFC1213-MIB"
	exit       = "exit"
)

func main() {
	if t, err := mib.Parse(RFC1213MIB); err == nil {
		fmt.Print(t)
		mainLoop(t)
	}
}

func mainLoop(t *mib.Tree) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
			mainLoop(t)
		}
	}()

	reader := bufio.NewReader(os.Stdin)
	for {
		if line, err := reader.ReadString('\n'); err == nil {
			if strings.TrimSpace(strings.ToLower(line)) == exit {
				break
			}

			splitted := strings.Split(line, " ")
			num, err := strconv.Atoi(splitted[0])

			if err != nil {
				panic(err)
			}

			fmt.Print(t.SubtreeString(mib.Oid{
				Number: num,
				Name:   splitted[1],
				Class:  splitted[2],
			}))
		}
	}
	ber.Encode("INTEGER", "-129")
}
