package slppclient_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestSlppclient(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Slppclient Suite")
}
