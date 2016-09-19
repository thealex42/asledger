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

	BeforeEach(func() {

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
		data.Set("from", "root")
		data.Add("to", "test")
		data.Add("amount", "42.2")

		req, err := http.NewRequest("POST", "/funds/move", bytes.NewBufferString(data.Encode()))
		Expect(err).To(BeNil())
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)
		Expect(resp.Code).To(Equal(200))
	})

	It("Should keep sequences", func() {

		modelRoot, err := NewBank("root", Clnt)
		modelBank1, err := NewBank("bank1", Clnt)

		modelRoot.addFunds(-100, Clnt)
		modelBank1.addFunds(100, Clnt)

		modelBank1, err = NewBank("bank1", Clnt)
		Expect(err).To(BeNil())
		Expect(modelBank1.Balance).To(Equal(100.0))

		modelRoot, err = NewBank("root", Clnt)
		Expect(err).To(BeNil())
		Expect(modelRoot.Seq).To(Equal(2))

		modelBank1, err = NewBank("bank1", Clnt)
		Expect(err).To(BeNil())
		modelBank2, err := NewBank("bank2", Clnt)
		Expect(err).To(BeNil())

		for i := 0; i < 100; i++ {
			modelBank1.addFunds(-1, Clnt)
			modelBank2.addFunds(1, Clnt)
		}

		modelBank1, err = NewBank("bank1", Clnt)
		Expect(err).To(BeNil())
		Expect(modelBank1.Balance).To(Equal(0.0))

		modelBank2, err = NewBank("bank2", Clnt)
		Expect(err).To(BeNil())
		Expect(modelBank2.Balance).To(Equal(100.0))

	})

	//	Measure("Balance req speed", func(b Benchmarker) {
	//		runtime := b.Time("run", func() {
	//			_, err := NewBank("bank1", Clnt)
	//			Expect(err).Should(BeNil())
	//		})
	//		Ω(runtime.Seconds()).Should(BeNumerically("<", 1), "SomethingHard() shouldn't take too long.")

	//	}, 1000)
})
