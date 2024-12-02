package input

import (
	"bytes"
	"os"
	"fmt"
)

// GenerateBinaryFile creates a binary file with commands to render the terminal screen.
func GenerateBinaryFile(filename string) error {
	var buffer bytes.Buffer

	// Add commands to the buffer
	// 0x1 - Screen setup: 80x24 characters, 256 colors
	buffer.WriteByte(0x1) // Command byte
	buffer.WriteByte(0x03) // Length byte
	buffer.Write([]byte{80, 24, 0x02}) // Data: width, height, color mode

	// 0x2 - Draw character: Place 'A' at (10, 5) with color 3
	buffer.WriteByte(0x2) // Command byte
	buffer.WriteByte(0x04) // Length byte
	buffer.Write([]byte{10, 5, 3, 'A'}) // Data: x, y, color, char

	// 0x3 - Draw line: From (20, 10) to (30, 15), color 4, using '-'
	buffer.WriteByte(0x3) // Command byte
	buffer.WriteByte(0x06) // Length byte
	buffer.Write([]byte{20, 10, 30, 15, 4, '-'}) // Data: x1, y1, x2, y2, color, char

	// 0x4 - Render text: "Hello" at (5, 3) with color 2
	buffer.WriteByte(0x4) // Command byte
	buffer.WriteByte(0x08) // Length byte
	buffer.Write([]byte{5, 3, 2})       // Data: x, y, color
	buffer.WriteString("Hello")         // Text

	// 0x7 - Clear screen
	buffer.WriteByte(0x7) // Command byte
	buffer.WriteByte(0x00) // Length byte (no data)

	// 0xFF - End of file
	buffer.WriteByte(0xFF) // Command byte
	buffer.WriteByte(0x00) // Length byte (no data)

	// Write the buffer to a file
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(buffer.Bytes())
	if err != nil {
		return err
	}

	fmt.Println("Binary file created:", filename)
	return nil
}
