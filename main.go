package main

import (
	"fmt"
	"github.com/bnch/banchoreader/lib"
	"github.com/codegangsta/cli"
	"github.com/fatih/color"
	"io/ioutil"
	"os"
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
			Name:  "i",
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
	red := color.New(color.FgRed, color.Bold)

	yellow.Printf("Reading file '%s'... ", filename)
	packets, err := banchoreader.ReadPackets(contents)
	if err != nil {
		red.Printf("Could not read '%s': %s", filename, err)
		return
	}
	green.Print("done. ")
	fmt.Printf("Read %d packets.\n", len(packets))

	r := banchoreader.New()
	r.Colored = !color.NoColor
	r.Ignored = ignored
	r.DumpPackets(os.Stdout, packets)

	fmt.Println()
}
