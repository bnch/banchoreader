#  banchoreader

A bancho packet reader for your terminal, written in Go.

![](https://y.zxq.co/wyzzsw.png)

Features:

- Outputs bancho packets with their components separated
- Converts Packet ID to human-readable string
- Hex editor-like raw byte output
- Automatic detection of bytes and int32s.
- **RAINBOWS**
- Lighting fast

## Installation

### The Go Way™

Assuming you have `$GOPATH/bin` in your `$PATH`, then it's just

```
go get github.com/bnch/banchoreader
```

### The n00b way

Grab a [build artifact](http://zxq.co:60291/view/bnch/banchoreader) and put it in your PATH somewhere.

## Usage

Can be either used passing the content through stdin...

```
cat sniffdump.txt | banchoreader
```

... or making it directly read some files!

```
banchoreader dump1.txt dump2.txt dump3.txt
```

