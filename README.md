# memory-dumper

## Overview

Read the memory of another process and save the data to a file.

After watching [Low Byte Productions' "Getting up in another processes memory"](https://www.youtube.com/watch?v=0ihChIaN8d0) video on how to dump the memory using a python script, I wanted to do the same thing with Go.


## Setup

You can compile the binary using the Makefile or the build command below:

```bash
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "-s -w" -trimpath -o memory-dumper
```

## Example

Dump the memory of a program called 'example'. You must run as root/sudo.

```bash
sudo ./memory-dumper -pid $(pidof example)
Wrote dump to: dump-1412302.bin
```

## References

- https://github.com/lowbyteproductions/memory-dumper
