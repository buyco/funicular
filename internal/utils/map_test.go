package utils_test

var _ = Describe("Map", func() {

	var (
		testMap = map[string]interface{}{"foo": 0}
	)

	It("should copy given map", func() {
		copyMap := CopyMap(testMap)
		Expect(copyMap["foo"]).To(BeZero())
	})
})
