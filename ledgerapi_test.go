package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Jeffail/gabs"
	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestLedger(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Ledger API test")
}

var _ = Describe("ledgerapi", func() {

	var r *gin.Engine

	BeforeSuite(func() {
		r = buildRouter()
	})

	It("Should have connections", func() {
		Expect(Clnt).ToNot(BeNil())
	})

	It("Should return some balance", func() {
		req, err := http.NewRequest("GET", "/balance/1", nil)
		Expect(err).To(BeNil())

		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)
		Expect(resp.Code).To(Equal(200))

		jsonParsed, err := gabs.ParseJSON([]byte(resp.Body.String()))
		Expect(jsonParsed.Exists("balance")).To(BeTrue())
	})
})
