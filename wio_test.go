package wio_test

import (
	"bytes"
	"io"
	"log"
	"math"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	wio "github.com/urjitbhatia/go-window-io"
)

var _ = Describe("Test window io", func() {
	Context("Rolling window - roll over by 1", func() {
		It("provides a rolling window over data", func() {
			wSize := 7
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

	})
})
