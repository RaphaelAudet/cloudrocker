package focker_test

import (
	. "github.com/cloudcredo/cloudfocker/Godeps/_workspace/src/github.com/onsi/ginkgo"
	. "github.com/cloudcredo/cloudfocker/Godeps/_workspace/src/github.com/onsi/gomega"

	"testing"
)

func TestFocker(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Focker Suite")
}
