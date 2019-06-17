package utils_test

import (
	. "github.com/buyco/funicular/internal/utils"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
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
