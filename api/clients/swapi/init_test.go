package swapi

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestSwapi(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Swapi Suite")
}
