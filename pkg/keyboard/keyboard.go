package keyboard

import (
	"io"
)

// Button represents a keyboard button pressed.
type Button = uint8

const (
	// UpArrow represents a keyboard up arrow button.
	UpArrow Button = iota
	// DownArrow represents a keyboard down arrow button.
	DownArrow
	// LeftArrow represents a keyboard left arrow button.
	LeftArrow
	// RightArrow represents a keyboard right arrow button.
	RightArrow
	// Esc represents a keyboard ESC button.
	Esc
)

// OnKeyDown listens to the given stream of pressed buttons.
func OnKeyDown(stream io.Reader, eventCh chan<- Button) {
	go func(stream io.Reader, eventCh chan<- Button) {
		buffer := make([]byte, 100)

		for {
			cnt, err := stream.Read(buffer)

			if err != nil {
				panic(err)
			}

			if cnt == 1 && buffer[0] == 0x1b {
				eventCh <- Esc
				continue
			}

			if buffer[0] == 0x1b && buffer[1] == '[' {
				switch buffer[2] {
				case 'A':
					eventCh <- UpArrow
				case 'B':
					eventCh <- DownArrow
				case 'C':
					eventCh <- RightArrow
				case 'D':
					eventCh <- LeftArrow
				}
			}
		}
	}(stream, eventCh)
}
