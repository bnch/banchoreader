package banchoreader

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/bnch/bancho/inbound"
	"github.com/bnch/bancho/pid"
	"github.com/fatih/color"
)

// Dumper is a container of settings of banchoreader, and is used to call the methods of banchoreader.
type Dumper struct {
	Ignored         []int
	Colored         bool
	IndentationSize int
}

// New returns a new empty Dumper.
func New() Dumper {
	colors()
	return Dumper{}
}

// Dump prints out to a certain file the human readable version of a/some bancho packet(s) passed as a byte slice.
func (d Dumper) Dump(file io.Writer, rawPackets []byte) error {
	packets, err := ReadPackets(rawPackets)
	if err != nil {
		return err
	}

	return d.DumpPackets(file, packets)
}

// DumpPackets prints out to a certain file the human readable version of some inbound.BasePacket passed as a slice.
func (d Dumper) DumpPackets(file io.Writer, packets []inbound.BasePacket) error {
	for _, packet := range packets {
		err := d.DumpPacket(file, packet)
		if err != nil {
			return err
		}
	}
	return nil
}

// DumpPacket prints out to a certain file the human readable version of a inbound.BasePacket.
func (d Dumper) DumpPacket(file io.Writer, packet inbound.BasePacket) error {
	if intInSlice(int(packet.ID), d.Ignored) {
		return nil
	}
	d.p(file, blue, "%s%s (%d)", d.indent(1), pid.String(packet.ID), packet.ID)
	switch len(packet.Content) {
	case 1:
		d.p(file, yellow, " (possible byte: %d)", packet.Content[0])
	case 4:
		var out int32
		binary.Read(bytes.NewReader(packet.Content), binary.LittleEndian, &out)
		d.p(file, yellow, " (possible int32: %d)", out)
	}
	d.p(file, def, "\n")
	d.hexdump(file, packet.Content)
	return nil
}

func (d Dumper) hexdump(file io.Writer, s []byte) {
	reader := bytes.NewReader(s)
	for {
		bf := make([]byte, 16)
		read, _ := reader.Read(bf)

		if read != 0 {
			d.p(file, def, "%s%-16s | % x\n", d.indent(2), safeString(bf)[:read], bf[:read])
		}

		if read < 16 {
			return
		}
	}
}
func (d Dumper) indent(multiply int) string {
	base := d.IndentationSize
	if base == 0 {
		base = 2
	}
	spaces := base * multiply

	var ret string
	for i := 0; i < spaces; i++ {
		ret += " "
	}
	return ret
}

// ReadPackets decodes a byte slice into a slice of inbound.BasePackets.
func ReadPackets(b []byte) ([]inbound.BasePacket, error) {
	var packets []inbound.BasePacket
	reader := bytes.NewReader(b)
	for {
		packet, err := inbound.GetPacket(reader)
		if err != io.EOF && err != nil {
			return packets, err
		}

		if !packet.Initialised {
			return packets, nil
		}

		packets = append(packets, packet)
	}
}
func intInSlice(a int, list []int) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
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

// MustDump does the same exact thing as Dump, but panics if an error happens.
func (d Dumper) MustDump(file io.Writer, rawPackets []byte) {
	err := d.Dump(file, rawPackets)
	if err != nil {
		panic(err)
	}
}

func (d Dumper) p(file io.Writer, color *color.Color, format string, data ...interface{}) {
	if d.Colored {
		fmt.Fprint(file, color.SprintfFunc()(format, data...))
		return
	}
	fmt.Fprintf(file, format, data...)
}
