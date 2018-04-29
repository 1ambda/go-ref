package distributed_test

import (
	"github.com/1ambda/go-ref/service-location/internal/distributed"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Distributed", func() {

	BeforeEach(func() {

	})

	AfterEach(func() {

	})

	Describe("New()", func() {
		It("should return Connector", func() {
			connector, err := distributed.New()
			Expect(connector).NotTo(BeNil())
			Expect(err).To(BeNil())
		})
	})

})
