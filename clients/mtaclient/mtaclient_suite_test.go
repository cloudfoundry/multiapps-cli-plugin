package mtaclient_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestMtaClient(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "MtaClient Suite")
}
