package mongodb

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestMongoDB(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "MongoDB Suite")
}
