package vektor_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "github.com/clarktrimble/stam/vektor"
)

func TestServer(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Vektor Suite")
}

var _ = Describe("Vektor", func() {
	var (
		vk *Vektor
	)

	BeforeEach(func() {
		vk = New(2)
		vk.Set(0, 0, 11)
		vk.Set(9, 9, 11)
	})

	Describe("Setting and getting values", func() {
		var (
			val  float64
			i, j int
		)

		JustBeforeEach(func() {
			val = vk.Get(i, j)
		})

		When("all goes well", func() {

			It("gets the value that was set", func() {
				Expect(val).To(Equal(11.0))
			})
		})

		When("grid coords are out of range", func() {
			BeforeEach(func() {
				i, j = 9, 9
			})

			It("gets zero and set is noop", func() {
				Expect(val).To(Equal(0.0))
			})
		})
	})

	Describe("Swapping values", func() {
		var (
			other *Vektor
		)

		JustBeforeEach(func() {
			vk.Swap(other)
		})

		When("all goes well", func() {
			BeforeEach(func() {
				other = New(2)
				other.Set(1, 1, 22)
			})

			It("swaps the values", func() {
				Expect(vk.Get(0, 0)).To(Equal(0.0))
				Expect(vk.Get(1, 1)).To(Equal(22.0))
				Expect(other.Get(0, 0)).To(Equal(11.0))
				Expect(other.Get(1, 1)).To(Equal(0.0))
			})
		})
	})
})
