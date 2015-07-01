package models_test

import (
	"github.com/mrap/guestimator/db"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"gopkg.in/pg.v3"

	"testing"
)

var DB *pg.DB

func TestModels(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Models Suite")
}

var _ = BeforeSuite(func() {
	DB = pg.Connect(db.Options("test"))
})

var _ = AfterSuite(func() {
	err := DB.Close()
	Expect(err).NotTo(HaveOccurred())
})
