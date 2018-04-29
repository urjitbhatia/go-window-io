package wio

import (
	"bufio"
	"container/ring"
	"io"
	"log"
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

// NewStepping converts an io.Reader into a windowed io.Reader with a window
// size of wSize, and a step size of sSize
func NewStepping(r io.Reader, wSize, sSize int) WindowReader {
	return &wReader{
		reader: bufio.NewReader(r),
		wSize:  wSize,
		sSize:  sSize,
		ring:   ring.New(wSize)}
}

// NewRolling converts an io.Reader into a rolling windowed io.Reader
// The size of a Rolling Window is 1.
// This is same as calling: NewStepping(reader, w, 1)
func NewRolling(r io.Reader, w int) WindowReader {
	return NewStepping(r, w, 1)
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
				PrintRing(w.ring, "before unlink")
				w.ring = w.ring.Unlink(w.sSize - i + 1)
				PrintRing(w.ring, "after unlink")
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
	vals := []interface{}{}
	for i := 0; i < r.Len(); i++ {
		vals = append(vals, r.Value)
		r = r.Next()
	}
	log.Printf("%s :: Ring: %v", msg, vals)
}
