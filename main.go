package main

import (
	"bytes"
	"fmt"
	"os"
	"encoding/binary"

	"github.com/nsf/termbox-go"
	"github.com/O-tero/terminal-screen/input"
	
)

const (
	CommandSetup        = 0x1
	CommandDrawChar     = 0x2
	CommandDrawLine     = 0x3
	CommandRenderText   = 0x4
	CommandCursorMove   = 0x5
	CommandDrawAtCursor = 0x6
	CommandClearScreen  = 0x7
	CommandEOF          = 0xFF
)

type Screen struct {
	width      int
	height     int
	colorMode  int
	buffer     [][]rune
	colorIndex [][]termbox.Attribute
}

func main() {
	// Generate the binary file
	err := input.GenerateBinaryFile("input.bin")
	if err != nil {
		fmt.Println("Error generating binary file:", err)
		return
	}

	// Ensure the input file exists
	if len(os.Args) < 2 {
		fmt.Println("Usage: main <input file>")
		return
	}

	inputFile := os.Args[1]
	fmt.Printf("Processing binary file: %s\n", inputFile)

	// Step 2: Initialize termbox
	if err := termbox.Init(); err != nil {
		fmt.Println("Error initializing terminal:", err)
		return
	}
	defer termbox.Close()

	screen := Screen{}
	data, err := os.ReadFile(inputFile)
	if err != nil {
		fmt.Println("Error reading input file:", err)
		return
	}
	reader := bytes.NewReader(data)
	var setupDone bool

	for reader.Len() > 0 {
		// Read command and length
		var command, length byte
		if err := readByte(reader, &command); err != nil {
			break
		}
		if command != CommandEOF {
			if err := readByte(reader, &length); err != nil {
				break
			}
		}

		// Process commands
		switch command {
		case CommandSetup:
			if length < 3 {
				fmt.Println("Invalid setup command")
				return
			}
			screen.Setup(reader)
			setupDone = true
		case CommandDrawChar:
			if !setupDone {
				fmt.Println("Screen not setup before draw character")
				continue
			}
			screen.DrawChar(reader)
		case CommandClearScreen:
			if setupDone {
				termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
			}
		case CommandEOF:
			termbox.Sync()
			return
		default:
			fmt.Println("Unknown command:", command)
		}

		// Render after each command
		screen.Render()
	}
}

// Helper functions to process commands
func (s *Screen) Setup(reader *bytes.Reader) {
	var width, height, colorMode byte
	readByte(reader, &width)
	readByte(reader, &height)
	readByte(reader, &colorMode)

	s.width = int(width)
	s.height = int(height)
	s.colorMode = int(colorMode)
	s.buffer = make([][]rune, s.height)
	s.colorIndex = make([][]termbox.Attribute, s.height)
	for i := range s.buffer {
		s.buffer[i] = make([]rune, s.width)
		s.colorIndex[i] = make([]termbox.Attribute, s.width)
	}
}

func (s *Screen) DrawChar(reader *bytes.Reader) {
	var x, y, color, char byte
	readByte(reader, &x)
	readByte(reader, &y)
	readByte(reader, &color)
	readByte(reader, &char)

	if int(x) < s.width && int(y) < s.height {
		s.buffer[int(y)][int(x)] = rune(char)
		s.colorIndex[int(y)][int(x)] = termbox.Attribute(color)
	}
}

func (s *Screen) Render() {
	for y, row := range s.buffer {
		for x, char := range row {
			color := termbox.ColorDefault
			if x < len(s.colorIndex[y]) {
				color = s.colorIndex[y][x]
			}
			termbox.SetCell(x, y, char, color, termbox.ColorDefault)
		}
	}
	termbox.Flush()
}

func readByte(reader *bytes.Reader, dest *byte) error {
	if err := binary.Read(reader, binary.LittleEndian, dest); err != nil {
		return err
	}
	return nil
}
