package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"syscall"
)

const (
	maps_info_regex string = "^([a-f0-9]+)\\-([a-f0-9]+)\\s(...)"
)

type memoryRegion struct {
	startaddress int64
	endaddress   int64
	permissions  string
}

func readMapsInfo(pid int) ([]memoryRegion, error) {
	var maps []memoryRegion

	filepath := fmt.Sprintf("/proc/%d/maps", pid)
	maps_file, err := os.ReadFile(filepath)
	if err != nil {
		return maps, err
	}

	maps_memory := strings.Split(string(maps_file), "\n")

	re, err := regexp.Compile(maps_info_regex)
	if err != nil {
		return maps, err
	}

	for _, m := range maps_memory {
		parts := re.FindStringSubmatch(m)
		if len(parts) == 4 {
			start, err := strconv.ParseInt(parts[1], 16, 64)
			if err != nil {
				continue
			}

			end, err := strconv.ParseInt(parts[2], 16, 64)
			if err != nil {
				continue
			}

			maps = append(maps, memoryRegion{
				startaddress: start,
				endaddress:   end,
				permissions:  parts[3],
			})
		}
	}

	return maps, nil
}

func main() {
	pid := flag.Int("pid", 0, "PID of the program to dump the memory from")
	flag.Parse()

	if *pid == 0 {
		fmt.Println("You must provide a PID")
		os.Exit(0)
	}

	maps, err := readMapsInfo(*pid)
	if err != nil {
		fmt.Printf("Error reading maps info: %v\n", err)
		os.Exit(1)
	}

	err = syscall.PtraceAttach(*pid)
	if err != nil {
		fmt.Printf("Error attaching to process: %v\n", err)
		os.Exit(1)
	}
	defer syscall.PtraceDetach(*pid)

	memFile := fmt.Sprintf("/proc/%d/mem", *pid)
	mem, err := os.Open(memFile)
	if err != nil {
		fmt.Printf("Error opening proc mem: %v\n", err)
		os.Exit(1)
	}
	defer mem.Close()

	outfilename := fmt.Sprintf("dump-%d.bin", *pid)
	outfile, err := os.Create(outfilename)
	if err != nil {
		fmt.Printf("Error creating dump file: %v\n", err)
		os.Exit(1)
	}
	defer outfile.Close()

	for _, memregion := range maps {
		if strings.Contains(memregion.permissions, "r") && strings.Contains(memregion.permissions, "w") {
			bloblength := memregion.endaddress - memregion.startaddress
			blob := make([]byte, bloblength)
			_, err := mem.ReadAt(blob, int64(memregion.startaddress))
			if err != nil {
				fmt.Println("Error reading memory")
				continue
			}
			outfile.Write(blob)
			clear(blob)
		}
	}
	fmt.Printf("Wrote dump to: %s\n", outfilename)
}
