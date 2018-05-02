# go-window-io 

[![Build Status](https://travis-ci.org/urjitbhatia/go-window-io.svg?branch=master)](https://travis-ci.org/urjitbhatia/go-window-io)
[![GoDoc](https://godoc.org/github.com/urjitbhatia/go-window-io?status.svg)](https://godoc.org/github.com/urjitbhatia/go-window-io)
[![Go Report Card](https://goreportcard.com/badge/github.com/urjitbhatia/go-window-io)](https://goreportcard.com/report/github.com/urjitbhatia/go-window-io)

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


# Benchmarks:
```
• [MEASUREMENT]
Test window io
/Users/UrjitSingh/AdvancedApps/gocode/src/github.com/urjitbhatia/go-window-io/wio_test.go:19
  benchmark rolling window
  /Users/UrjitSingh/AdvancedApps/gocode/src/github.com/urjitbhatia/go-window-io/wio_test.go:199

  Ran 1000 samples:
  runtime:
    Fastest Time: 0.002s
    Slowest Time: 0.004s
    Average Time: 0.003s ± 0.000s
------------------------------

Ran 12 of 12 Specs in 2.566 seconds
SUCCESS! -- 12 Passed | 0 Failed | 0 Pending | 0 Skipped --- PASS: TestGoWindowIo (2.57s)
goos: darwin
goarch: amd64
pkg: github.com/urjitbhatia/go-window-io
BenchmarkRollingWindow/BufSize:_10240_wSize:_1-8         	50000000	        31.3 ns/op	       0 B/op	       0 allocs/op
BenchmarkRollingWindow/BufSize:_10240_wSize:_103-8       	50000000	        33.4 ns/op	       0 B/op	       0 allocs/op
BenchmarkRollingWindow/BufSize:_10240_wSize:_205-8       	50000000	        31.3 ns/op	       0 B/op	       0 allocs/op
BenchmarkRollingWindow/BufSize:_10240_wSize:_307-8       	50000000	        32.6 ns/op	       0 B/op	       0 allocs/op
BenchmarkRollingWindow/BufSize:_10240_wSize:_409-8       	50000000	        31.7 ns/op	       0 B/op	       0 allocs/op
BenchmarkRollingWindow/BufSize:_10240_wSize:_511-8       	50000000	        33.0 ns/op	       0 B/op	       0 allocs/op
BenchmarkRollingWindow/BufSize:_10240_wSize:_613-8       	50000000	        32.7 ns/op	       0 B/op	       0 allocs/op
BenchmarkRollingWindow/BufSize:_10240_wSize:_715-8       	50000000	        33.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkRollingWindow/BufSize:_10240_wSize:_817-8       	50000000	        30.6 ns/op	       0 B/op	       0 allocs/op
BenchmarkRollingWindow/BufSize:_10240_wSize:_919-8       	50000000	        33.2 ns/op	       0 B/op	       0 allocs/op
BenchmarkRollingWindow/BufSize:_10240_wSize:_1021-8      	50000000	        32.0 ns/op	       0 B/op	       0 allocs/op
BenchmarkRollingWindow/BufSize:_10240_wSize:_1123-8      	30000000	        33.4 ns/op	       0 B/op	       0 allocs/op
BenchmarkRollingWindow/BufSize:_10240_wSize:_1225-8      	50000000	        32.6 ns/op	       0 B/op	       0 allocs/op
BenchmarkRollingWindow/BufSize:_10240_wSize:_1327-8      	50000000	        32.9 ns/op	       0 B/op	       0 allocs/op
BenchmarkRollingWindow/BufSize:_10240_wSize:_1429-8      	50000000	        31.6 ns/op	       0 B/op	       0 allocs/op
BenchmarkRollingWindow/BufSize:_10240_wSize:_1531-8      	50000000	        33.2 ns/op	       0 B/op	       0 allocs/op
BenchmarkRollingWindow/BufSize:_10240_wSize:_1633-8      	50000000	        31.7 ns/op	       0 B/op	       0 allocs/op
BenchmarkRollingWindow/BufSize:_10240_wSize:_1735-8      	50000000	        31.9 ns/op	       0 B/op	       0 allocs/op
BenchmarkRollingWindow/BufSize:_10240_wSize:_1837-8      	50000000	        32.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkRollingWindow/BufSize:_10240_wSize:_1939-8      	50000000	        33.3 ns/op	       0 B/op	       0 allocs/op
BenchmarkRollingWindow/BufSize:_10240_wSize:_2041-8      	50000000	        32.1 ns/op	       0 B/op	       0 allocs/op
PASS
ok  	github.com/urjitbhatia/go-window-io	44.603s
```