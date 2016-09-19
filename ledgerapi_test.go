package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
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

	It("Should move funds", func() {
		data := url.Values{}
		data.Set("from", "bank1")
		data.Add("to", "bank2")
		data.Add("amount", "42.2")

		req, err := http.NewRequest("POST", "/funds/move", bytes.NewBufferString(data.Encode()))
		Expect(err).To(BeNil())
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)
		Expect(resp.Code).To(Equal(200))
	})

	Measure("Balance req speed", func(b Benchmarker) {
		runtime := b.Time("run", func() {
			_, err := NewBank("bank1", Clnt)
			Expect(err).Should(BeNil())
		})
		Ω(runtime.Seconds()).Should(BeNumerically("<", 1), "SomethingHard() shouldn't take too long.")

	}, 1000)
})
