package output

import (
	"bytes"
	"errors"
)

var (
	ErrLineBufferOverflow = errors.New("line buffer overflow")

	ErrAlreadyFinished      = errors.New("already finished")
	ErrNotFoundCommand      = errors.New("command not found")
	ErrNotExecutePermission = errors.New("not execute permission")
	ErrInvalidArgs          = errors.New("Invalid argument to exit")
	ErrProcessTimeout       = errors.New("throw process timeout")
	ErrProcessCancel        = errors.New("active cancel process")

	DefaultExitCode = 2
)

type OutputStream struct {
	streamChan chan string
	bufSize    int
	buf        []byte
	lastChar   int
}

// NewOutputStream creates a new streaming output on the given channel.
func NewOutputStream(streamChan chan string) *OutputStream {
	out := &OutputStream{
		streamChan: streamChan,
		bufSize:    16384,
		buf:        make([]byte, 16384),
		lastChar:   0,
	}
	return out
}

// Write makes OutputStream implement the io.Writer interface.
func (rw *OutputStream) Write(p []byte) (n int, err error) {
	n = len(p) // end of buffer
	firstChar := 0

	for {
		newlineOffset := bytes.IndexByte(p[firstChar:], '\n')
		if newlineOffset < 0 {
			break // no newline in stream, next line incomplete
		}

		// End of line offset is start (nextLine) + newline offset. Like bufio.Scanner,
		// we allow \r\n but strip the \r too by decrementing the offset for that byte.
		lastChar := firstChar + newlineOffset // "line\n"
		if newlineOffset > 0 && p[newlineOffset-1] == '\r' {
			lastChar -= 1 // "line\r\n"
		}

		// Send the line, prepend line buffer if set
		var line string
		if rw.lastChar > 0 {
			line = string(rw.buf[0:rw.lastChar])
			rw.lastChar = 0 // reset buffer
		}
		line += string(p[firstChar:lastChar])
		rw.streamChan <- line // blocks if chan full

		// Next line offset is the first byte (+1) after the newline (i)
		firstChar += newlineOffset + 1
	}

	if firstChar < n {
		remain := len(p[firstChar:])
		bufFree := len(rw.buf[rw.lastChar:])
		if remain > bufFree {
			var line string
			if rw.lastChar > 0 {
				line = string(rw.buf[0:rw.lastChar])
			}
			line += string(p[firstChar:])
			err = ErrLineBufferOverflow
			n = firstChar
			return // implicit
		}
		copy(rw.buf[rw.lastChar:], p[firstChar:])
		rw.lastChar += remain
	}

	return // implicit
}

func (rw *OutputStream) Lines() <-chan string {
	return rw.streamChan
}

func (rw *OutputStream) SetLineBufferSize(n int) {
	rw.bufSize = n
	rw.buf = make([]byte, rw.bufSize)
}
