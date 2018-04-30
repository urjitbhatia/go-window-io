package wio

import (
	"bufio"
	"container/ring"
	"errors"
	"fmt"
	"io"
	"log"
	"strings"
)

// WindowReader wraps io.Reader and reads in a window of data
// The window size and the window rolling step size can be configured
type WindowReader interface {
	io.Reader
}

type wReader struct {
	firstWindowDone bool          // if true, this is not the first window
	reader          *bufio.Reader // The input io.Reader
	sSize           int           // The number of bytes to move window by
	wSize           int           // The size of the window
	ring            *ring.Ring    // Internal ring of data
}

// ErrDisjointWindow - Step size has to be >= 1 && <= window size
var ErrDisjointWindow = errors.New("Step size cannot be larger than window size")

// ErrZeroWindowSize - Window size has to be >= 1
var ErrZeroWindowSize = errors.New("Window size cannot be zero. Should be 1 or more")

// ErrZeroStepSize - Step size has to be >= 1
var ErrZeroStepSize = errors.New("Step size cannot be zero. Should be 1 or more")

// NewStepping converts an io.Reader into a windowed io.Reader with a window
// size of wSize, and a step size of sSize
func NewStepping(r io.Reader, wSize, sSize int) (WindowReader, error) {
	if sSize == 0 {
		return nil, ErrZeroStepSize
	}
	if sSize > wSize {
		return nil, ErrDisjointWindow
	}
	return &wReader{
			reader: bufio.NewReader(r),
			wSize:  wSize,
			sSize:  sSize,
			ring:   ring.New(wSize)},
		nil
}

// NewRolling converts an io.Reader into a rolling windowed io.Reader
// The size of a Rolling Window is 1.
// This is same as calling: NewStepping(reader, w, 1)
func NewRolling(r io.Reader, wSize int) (WindowReader, error) {
	if wSize == 0 {
		return nil, ErrZeroWindowSize
	}
	return NewStepping(r, wSize, 1)
}

// Read the next "window" into the given buffer
func (w *wReader) Read(buf []byte) (int, error) {

	// If buf can't hold the window, throw short buffer err
	if len(buf) < w.wSize {
		return 0, io.ErrShortBuffer
	}

	// For a regular read (not the very first window, read up to step sizes)
	readCount := w.sSize
	if !w.firstWindowDone {
		// For the very first read, read up to the entire window size to fill up the window
		readCount = w.wSize
		w.firstWindowDone = true
	}

	// Read data from the input
	for i := 0; i < readCount; i++ {
		b, err := w.reader.ReadByte()
		if err != nil {
			if err != io.EOF {
				// Something is wrong, return err
				return 0, err
			}
			// EOF - check if we managed to read something in our window
			if i > 0 {
				// We consumed some last remaining data, put it in the ring
				// and unlink the things we consumed
				w.ring = w.ring.Prev()
				w.ring.Unlink(w.sSize - i)
				w.ring = w.ring.Next()
				break
			}
			return 0, err
		}
		w.ring.Value = b
		w.ring = w.ring.Next()
	}
	bufIdx := 0
	w.ring.Do(func(val interface{}) {
		buf[bufIdx] = val.(byte)
		bufIdx++
	})
	return bufIdx, nil
}

// PrintRing is a helpful debug function that prints a ring
// with the given message
func PrintRing(r *ring.Ring, msg string) {
	s := []string{}
	for i := 0; i < r.Len(); i++ {
		s = append(s, fmt.Sprintf("%v", r.Value))
		r = r.Next()
	}

	log.Printf("%s :: Ring: [%s]", msg, strings.Join(s, ", "))
}
