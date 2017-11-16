package name_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestName(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Name Suite")
}
