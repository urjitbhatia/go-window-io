package wio_test

import (
	"bytes"
	"container/ring"
	"fmt"
	"io"
	"log"
	"math"
	"math/rand"
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	wio "github.com/urjitbhatia/go-window-io"
)

var _ = Describe("Test window io", func() {

	Context("Rolling window - roll over by 1", func() {
		It("provides a rolling window over data", func() {
			wSize := 3
			rollingBuf := []byte("\x00\x01\x02\x03\x04\x05\x06\x07\x08") // len 9
			rollingReader, err := wio.NewRolling(bytes.NewBuffer(rollingBuf), wSize)
			Expect(err).To(BeNil())

			// Number of rolling reads should be: bufLen - wSize + 1
			count := len(rollingBuf) - wSize + 1

			buf := [64]byte{}
			for i := 0; i < count; i++ {
				// For each rolling Hash, check it has the right value
				b := rollingBuf[i : i+wSize]
				n, err := rollingReader.Read(buf[:wSize])
				Expect(err).To(BeNil())
				Expect(n > 0).To(BeTrue())
				Expect(buf[:n]).To(Equal(b))
			}
			// Next one should err with EOF
			_, err = rollingReader.Read(buf[:wSize])
			Expect(err).To(Equal(io.EOF))
		})

		It("provides a rolling window over data - random window size (upto data len)", func() {
			rand.Seed(time.Now().Unix())
			rollingBuf := []byte("\x00\x01\x02\x03\x04\x05\x06\x07\x08\x00\x01\x02\x03\x04\x05\x06\x07\x08")
			wSize := rand.Intn(len(rollingBuf) + 1)
			rollingReader, err := wio.NewRolling(bytes.NewBuffer(rollingBuf), wSize)
			Expect(err).To(BeNil())

			// Number of rolling reads should be: bufLen - wSize + 1
			count := len(rollingBuf) - wSize + 1

			buf := [64]byte{}
			for i := 0; i < count; i++ {
				// For each rolling Hash, check it has the right value
				b := rollingBuf[i : i+wSize]
				n, err := rollingReader.Read(buf[:wSize])
				Expect(err).To(BeNil())
				Expect(n > 0).To(BeTrue())
				Expect(buf[:n]).To(Equal(b))
			}
			// Next one should err with EOF
			_, err = rollingReader.Read(buf[:wSize])
			Expect(err).To(Equal(io.EOF))
		})

		It("returns a short buffer error if read buf is smaller than window size", func() {
			rollingBuf := []byte("\x00\x01\x02\x03\x04\x05\x06\x07\x08\x00\x01\x02\x03\x04\x05\x06\x07\x08")
			wSize := 5
			rollingReader, err := wio.NewRolling(bytes.NewBuffer(rollingBuf), wSize)
			Expect(err).To(BeNil())

			buf := [5]byte{}
			n, err := rollingReader.Read(buf[:wSize-2])
			Expect(err).To(Equal(io.ErrShortBuffer))
			Expect(n).To(Equal(0))
		})
	})

	Context("Stepping window - roll over more than 1", func() {
		It("provides a rolling window over data", func() {
			wSize := 4
			sSize := 3
			rollingBuf := []byte("\x00\x01\x02\x03\x04\x05\x06\x07\x08\x09") // len 10
			rollingReader, _ := wio.NewStepping(bytes.NewBuffer(rollingBuf), wSize, sSize)

			// Number of rolling reads should be: math.Ceil((bufLen - wSize)/sSize + 1)
			count := int(math.Ceil(float64(len(rollingBuf)-wSize)/float64(sSize)) + 1.00)
			log.Printf("data len: %d, count : %d", len(rollingBuf), count)
			tail := 0
			head := tail + wSize
			buf := [64]byte{}
			for i := 0; i < int(count); i++ {
				// For each rolling Hash, check it has the right value
				b := rollingBuf[tail:head]
				n, err := rollingReader.Read(buf[:wSize])
				Expect(err).To(BeNil())
				Expect(n > 0).To(BeTrue())
				Expect(buf[:n]).To(Equal(b))
				tail += sSize
				head += sSize
				if head > len(rollingBuf) {
					head = len(rollingBuf)
				}
			}
			// Next one should err with EOF
			_, err := rollingReader.Read(buf[:wSize])
			Expect(err).To(Equal(io.EOF))
		})

		It("provides a rolling window over data - random window and step size", func() {
			rand.Seed(time.Now().Unix())

			// fill in a random buffer
			rollingBuf := make([]byte, rand.Intn(250)+10)
			_, err := rand.Read(rollingBuf)
			Expect(err).To(BeNil())

			wSize := rand.Intn(len(rollingBuf)) + 1
			sSize := rand.Intn(wSize) + 1

			rollingReader, _ := wio.NewStepping(bytes.NewBuffer(rollingBuf), wSize, sSize)

			// Number of rolling reads should be: math.Ceil((bufLen - wSize)/sSize + 1)
			count := int(math.Ceil(float64(len(rollingBuf)-wSize)/float64(sSize)) + 1.00)
			log.Printf("wSize: %d, sSize: %d, data len: %d, count : %d", wSize, sSize, len(rollingBuf), count)

			tail := 0
			head := tail + wSize
			buf := [300]byte{}
			for i := 0; i < int(count); i++ {
				// For each rolling Hash, check it has the right value
				b := rollingBuf[tail:head]
				n, err := rollingReader.Read(buf[:wSize])
				Expect(err).To(BeNil())
				Expect(n > 0).To(BeTrue())
				Expect(buf[:n]).To(Equal(b))
				tail += sSize
				head += sSize
				if head > len(rollingBuf) {
					head = len(rollingBuf)
				}
			}
			// Next one should err with EOF
			_, err = rollingReader.Read(buf[:wSize])
			Expect(err).To(Equal(io.EOF))
		})

		It("returns a short buffer error if read buf is smaller than window size", func() {
			rollingBuf := []byte("\x00\x01\x02\x03\x04\x05\x06\x07\x08\x00\x01\x02\x03\x04\x05\x06\x07\x08")
			wSize := 5
			sSize := 5
			rollingReader, _ := wio.NewStepping(bytes.NewBuffer(rollingBuf), wSize, sSize)

			buf := [5]byte{}
			n, err := rollingReader.Read(buf[:wSize-2])
			Expect(err).To(Equal(io.ErrShortBuffer))
			Expect(n).To(Equal(0))
		})
	})

	Context("Errors on bad window and step sizes", func() {
		It("Creates a window of size 0", func() {
			_, err := wio.NewRolling(bytes.NewBufferString("foo"), 0)
			Expect(err).To(Equal(wio.ErrZeroWindowSize))
		})

		It("Creates a step of size 0", func() {
			_, err := wio.NewStepping(bytes.NewBufferString("foo"), 2, 0)
			Expect(err).To(Equal(wio.ErrZeroStepSize))
		})
		It("Creates a step size larger than window size", func() {
			_, err := wio.NewStepping(bytes.NewBufferString("foo"), 2, 3)
			Expect(err).To(Equal(wio.ErrDisjointWindow))
		})
		It("Creates a step size equal to window size - should not fail", func() {
			_, err := wio.NewStepping(bytes.NewBufferString("foo"), 2, 2)
			Expect(err).To(BeNil())
		})
	})

	Context("Print ring", func() {
		It("prints the current state of container/ring", func() {
			r := ring.New(5)
			for i := 0; i < r.Len(); i++ {
				r.Value = i
				r = r.Next()
			}
			buf := bytes.Buffer{}
			log.SetOutput(&buf)
			wio.PrintRing(r, "test ring")
			Expect(buf.String()).To(ContainSubstring("test ring :: Ring: [0, 1, 2, 3, 4]"))
			log.SetOutput(GinkgoWriter)
		})
	})

	Measure("benchmark rolling window", func(b Benchmarker) {
		rand.Seed(time.Now().Unix())
		wSize := 3
		rollingBuf := [1024 * 10]byte{} //10 kb buffer
		_, err := rand.Read(rollingBuf[:])
		Expect(err).To(BeNil())
		rollingReader, err := wio.NewRolling(bytes.NewBuffer(rollingBuf[:]), wSize)
		Expect(err).To(BeNil())

		// Number of rolling reads should be: bufLen - wSize + 1
		buf := [64]byte{}

		runtime := b.Time("runtime", func() {
			for {
				// For each rolling Hash, check it has the right value
				n, err := rollingReader.Read(buf[:wSize])
				if err == io.EOF {
					break
				}
				Expect(n).To(Equal(wSize))
			}
		})

		Ω(runtime.Seconds()).Should(BeNumerically("<", 0.2), "SomethingHard() shouldn't take too long.")
	}, 1000)
})

func BenchmarkRollingWindow(b *testing.B) {
	// run the Fib function b.N times
	rand.Seed(time.Now().Unix())
	rollingBuf := [1024 * 10]byte{} //10 kb buffer
	_, err := rand.Read(rollingBuf[:])
	if err != nil {
		b.Error(err)
	}
	b.ResetTimer()
	for wSize := 1; wSize < len(rollingBuf)/10; wSize += len(rollingBuf) / 50 {
		b.StopTimer()
		rollingReader, err := wio.NewRolling(bytes.NewBuffer(rollingBuf[:]), wSize)
		if err != nil {
			b.Error(err)
		}
		buf := make([]byte, wSize)
		b.StartTimer()
		b.Run(fmt.Sprintf("BufSize: %d wSize: %d", len(rollingBuf), wSize), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				for {
					// For each rolling Hash, check it has the right value
					n, err := rollingReader.Read(buf[:wSize])
					if err == io.EOF {
						break
					}
					if n != wSize {
						b.Fatal("expected byte read count: ", wSize)
					}
				}
			}
		})
	}
}
