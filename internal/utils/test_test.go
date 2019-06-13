package utils_test

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/buyco/funicular/internal/utils"
)

var _ = Describe("Test", func() {

	It("should catch stdout", func() {
		stdout := CaptureStdout(func() { fmt.Sprint("foo:bar") })
		Expect(stdout).ToNot(ContainSubstring("foo:bar"))
	})
})
