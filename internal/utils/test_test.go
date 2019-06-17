package utils_test

import (
	. "github.com/buyco/funicular/internal/utils"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"fmt"
)

var _ = Describe("Test", func() {

	It("should catch stdout", func() {
		stdout := CaptureStdout(func() { fmt.Sprint("foo:bar") })
		Expect(stdout).ToNot(ContainSubstring("foo:bar"))
	})
})
