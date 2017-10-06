package restclient_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestRestclient(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Restclient Suite")
}
