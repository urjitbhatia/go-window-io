package wio_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestGoWindowIo(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "GoWindowIo Suite")
}
