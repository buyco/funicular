package sync_test

import (
	"github.com/buyco/funicular/pkg/sync"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Pool", func() {

	var pool *sync.Pool

	BeforeEach(func() {
		pool = sync.NewPool(1, nil)
	})

	Context("From constructor function", func() {

		It("creates a valid instance", func() {
			Expect(pool).To(BeAssignableToTypeOf(&sync.Pool{}))
		})

		It("contains zero clients", func() {
			cli := pool.Get()
			Expect(cli).To(BeNil())
		})

		It("gets capacity", func() {
			capacity := pool.GetCapacity()
			Expect(capacity).To(Equal(uint(1)))
		})

		It("puts in pool", func() {
			err := pool.Put("foo")
			Expect(err).ToNot(HaveOccurred())
		})
	})

	It("sets factory", func() {
		pool.SetFactory(func() interface{} { return "bar" })
		Expect(pool.Get()).To(Equal("bar"))
	})
})
