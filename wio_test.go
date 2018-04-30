package wio_test

import (
	"bytes"
	"container/ring"
	"io"
	"log"
	"math"
	"math/rand"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	wio "github.com/urjitbhatia/go-window-io"
)

var _ = Describe("Test window io", func() {
	Context("Rolling window - roll over by 1", func() {
		It("provides a rolling window over data", func() {
			wSize := 3
			rollingBuf := []byte("\x00\x01\x02\x03\x04\x05\x06\x07\x08") // len 9
			rollingReader := wio.NewRolling(bytes.NewBuffer(rollingBuf), wSize)

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
			_, err := rollingReader.Read(buf[:wSize])
			Expect(err).To(Equal(io.EOF))
		})

		It("provides a rolling window over data - random window size (upto data len)", func() {
			rollingBuf := []byte("\x00\x01\x02\x03\x04\x05\x06\x07\x08\x00\x01\x02\x03\x04\x05\x06\x07\x08")
			wSize := rand.Intn(len(rollingBuf) + 1)
			rollingReader := wio.NewRolling(bytes.NewBuffer(rollingBuf), wSize)

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
			_, err := rollingReader.Read(buf[:wSize])
			Expect(err).To(Equal(io.EOF))
		})

		It("returns a short buffer error if read buf is smaller than window size", func() {
			rollingBuf := []byte("\x00\x01\x02\x03\x04\x05\x06\x07\x08\x00\x01\x02\x03\x04\x05\x06\x07\x08")
			wSize := 5
			rollingReader := wio.NewRolling(bytes.NewBuffer(rollingBuf), wSize)

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
			rollingReader := wio.NewStepping(bytes.NewBuffer(rollingBuf), wSize, sSize)

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

		It("returns a short buffer error if read buf is smaller than window size", func() {
			rollingBuf := []byte("\x00\x01\x02\x03\x04\x05\x06\x07\x08\x00\x01\x02\x03\x04\x05\x06\x07\x08")
			wSize := 5
			sSize := 5
			rollingReader := wio.NewStepping(bytes.NewBuffer(rollingBuf), wSize, sSize)

			buf := [5]byte{}
			n, err := rollingReader.Read(buf[:wSize-2])
			Expect(err).To(Equal(io.ErrShortBuffer))
			Expect(n).To(Equal(0))
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
})
