package database

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestDatabase(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Database Suite")
}

var _ = Describe("Database", func() {
	Context("when trying to initialize the databas", func() {
		err := InitDatabase()
		It("it does not fail", func() {
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
