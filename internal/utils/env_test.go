package utils_test

var _ = Describe("Env", func() {

	var err error

	Context("With allowed environment", func() {

		It("should fail to load non-existent file", func() {
			stdout := CaptureStdout(func() { err = LoadEnvFile("foo", "development") })
			Expect(stdout).To(BeEmpty())
			Expect(err).To(HaveOccurred())
		})

		It("should load env file", func() {
			stdout := CaptureStdout(func() { err = LoadEnvFile("../../.env-example", "development") })
			Expect(stdout).To(BeEmpty())
			Expect(err).ToNot(HaveOccurred())
		})
	})

	Context("With not allowed environment", func() {

		It("should fail to load non-existent file", func() {
			stdout := CaptureStdout(func() { err = LoadEnvFile("foo", "bar") })
			Expect(stdout).To(ContainSubstring("Environment file not loaded for the current env"))
			Expect(err).ToNot(HaveOccurred())
		})
	})
})
