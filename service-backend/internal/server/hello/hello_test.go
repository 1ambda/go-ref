package hello_test

import (
	. "github.com/1ambda/go-ref/service-backend/internal/server/hello"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Name", func() {

	owner := "1ambda"
	invalid := "2ambda"
	anonymous := "gopher"

	Describe("GetName", func() {
		Context("owner", func() {
			It("should handle owner", func() {
				name, _ := GetName(owner)
				Expect(name).To(Equal("Kun"))
			})
		})

		Context("anonymous", func() {
			It("should return the same name", func() {
				name, _ := GetName(anonymous)
				Expect(name).To(Equal(anonymous))
			})
		})

		Context("invalid", func() {
			It("should throw error", func() {
				_, err := GetName(invalid)
				Expect(err).NotTo(Equal(nil))
			})
		})
	})
})
