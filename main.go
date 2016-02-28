package main

import (
	"bytes"
	"fmt"
	"github.com/bnch/bancho/inbound"
	"github.com/bnch/bancho/pid"
	"io/ioutil"
	"os"
	"io"
	"encoding/binary"
	"github.com/fatih/color"
	"github.com/codegangsta/cli"
)

func main() {
	app := cli.NewApp()
	
	app.Name = "banchoreader"
	app.Version = "1.0.0"
	app.Action = mainCommand
	app.Usage = "read bancho packets in an elegant way"
	app.ArgsUsage = "[files ...]"
	app.Author = "Howl"
	app.Flags = []cli.Flag{
		cli.IntSliceFlag{
			Name: "i",
			Usage: "Packet IDs to ignore",
		},
	}
	
	app.Run(os.Args)
}

func mainCommand(c *cli.Context) {
	ignored := c.IntSlice("i")
	
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		// Some stuff is being dumped in stdin
		data, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			panic(err)
		}
		output("stdin", data, ignored)
	}

	files := c.Args()
	for _, filename := range files {
		data, err := ioutil.ReadFile(filename)
		if err != nil {
			fmt.Printf("Could not read %s: %s\n", filename, err)
			continue
		}
		output(filename, data, ignored)
	}
}

func output(filename string, contents []byte, ignored []int) {
	green := color.New(color.Bold, color.FgGreen)
	yellow := color.New(color.FgYellow)
	blue := color.New(color.FgBlue)
	
	yellow.Printf("Reading file '%s'... ", filename)
	packets := readPackets(contents)
	green.Print("done. ")
	fmt.Printf("Read %d packets.\n", len(packets))

	for _, packet := range packets {
		if intInSlice(int(packet.ID), ignored) {
			continue
		}
		blue.Printf("  %s (%d)", pid.String(packet.ID), packet.ID)
		switch len(packet.Content) {
		case 1:
			yellow.Printf(" (possible byte: %d)", packet.Content[0])
		case 4:
			var out int32
			binary.Read(bytes.NewReader(packet.Content), binary.LittleEndian, &out)
			yellow.Printf(" (possible int32: %d)", out)
		}
		fmt.Println()
		hexdump(packet.Content)
	}
	fmt.Println()
}

func readPackets(b []byte) (packets []inbound.BasePacket) {
	reader := bytes.NewReader(b)
	for {
		packet, err := inbound.GetPacket(reader)
		if err != io.EOF && err != nil {
			panic(err)
		}

		if !packet.Initialised {
			return
		}

		packets = append(packets, packet)
	}
}

func hexdump(s []byte) {
	reader := bytes.NewReader(s)
	for {
		bf := make([]byte, 16)
		read, _ := reader.Read(bf)

		if read != 0 {
			fmt.Printf("    %-16s | % x\n", safeString(bf)[:read], bf[:read])
		}

		if read < 16 {
			return
		}
	}
}

func safeString(s []byte) string {
	var ret string
	for _, b := range s {
		if b < 32 || b > 126 {
			ret += "."
		} else {
			ret += string(b)
		}
	}
	return ret
}

func intInSlice(a int, list []int) bool {
    for _, b := range list {
        if b == a {
            return true
        }
    }
    return false
}

