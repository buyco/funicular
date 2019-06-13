package utils_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/buyco/funicular/internal/utils"
)

var _ = Describe("Error", func() {

	It("should create error with variables args", func() {
		err := ErrorPrintf("foo %s", "bar")
		Expect(err.Error()).To(MatchRegexp("foo bar"))
	})

	It("should create error with string", func() {
		err := ErrorPrint("foo")
		Expect(err.Error()).To(MatchRegexp("foo"))
	})
})
