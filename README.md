# go-window-io 

[![Build Status](https://travis-ci.org/urjitbhatia/go-window-io.svg?branch=master)](https://travis-ci.org/urjitbhatia/go-window-io)
[![GoDoc](https://godoc.org/github.com/urjitbhatia/go-window-io?status.svg)](https://godoc.org/github.com/urjitbhatia/go-window-io)

A Sliding window over arbitrary io.Readers in go.

`wio` or `windowed-io` wraps an `io.Reader` and implements `io.Read` but reads data chunked over a *Sliding Window*.

`wio` provides two controls: `window size` and the `step size`.

For a canonical `Rolling window`, the `step size` is 1.

### Examples

```golang
wSize := 3 // window size of 3

rollingBuf := []byte("\x00\x01\x02\x03\x04\x05\x06\x07\x08")
reader := bytes.NewBuffer(rollingBuf)
rollingReader := wio.NewRolling(reader, wSize)

// For a window that rolls by more than 1:
// sSize := 2 // window of 3 slide by 2
// steppingReader := wio.NewRolling(reader, wSize, sSize)

for {
  // buf is some buffer at-least wSize in size
  n, err := rollingReader.Read(buf)
  if err != nil {
    if err == io.EOF {
      // handle EOF
    }
    if err == io.ErrShortBuffer {
      // len of `buf` was shorter than wSize - make sure the given buf has correct size
    }
  }

  // n can be less than wSize for the last window
  consumeData(buf[:n])
}
```
