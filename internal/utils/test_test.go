package utils_test

import (
	"fmt"
)

var _ = Describe("Test", func() {

	It("should catch stdout", func() {
		stdout := CaptureStdout(func() { fmt.Sprint("foo:bar") })
		Expect(stdout).ToNot(ContainSubstring("foo:bar"))
	})
})
