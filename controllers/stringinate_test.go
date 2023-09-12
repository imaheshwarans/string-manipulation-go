package controllers

import (
	"net/http"
	"net/http/httptest"
	"stringinator-go/constants"
	"stringinator-go/utils"
	"testing"

	"github.com/gorilla/mux"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func Test_allShortestStrings(t *testing.T) {

	strList = []string{"ihh", "hi", "a"}

	tests := []struct {
		name string
	}{
		{name: "valid case"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			allShortestStrings()
		})
	}
}

var _ = Describe("StringController", func() {
	var router *mux.Router
	var w *httptest.ResponseRecorder

	stringController := StringInate{}

	strList = []string{"ihh", "hi", "a"}

	BeforeEach(func() {
		router = mux.NewRouter()

	})

	Describe("GetStringInate", func() {
		Context("Get request with valid inputs", func() {
			It("Should return appropriate response", func() {
				router.Handle("/stats", utils.ErrorHandler(utils.JsonResponseHandler(stringController.Get))).Methods("GET")
				req, err := http.NewRequest(http.MethodGet, "/stats", nil)
				Expect(err).NotTo(HaveOccurred())
				req.Header.Set("Accept", constants.HTTPMediaTypeJson)
				req.Header.Set("Content-Type", constants.HTTPMediaTypeJson)
				w = httptest.NewRecorder()
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusOK))
			})
		})
	})
})
