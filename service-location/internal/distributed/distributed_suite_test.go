package distributed_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var T *testing.T = nil

func TestDistributed(t *testing.T) {
	T = t

	RegisterFailHandler(Fail)
	RunSpecs(t, "Distributed Suite")
}
