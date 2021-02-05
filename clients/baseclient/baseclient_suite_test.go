package baseclient_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestBaseclient(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Baseclient Suite")
}
