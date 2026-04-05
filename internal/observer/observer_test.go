package observer

import (
	"testing"
	"time"

	"github.com/pedrowilliamss/pingero/internal/testutil"
)

func TestObserverWithViewer(t *testing.T) {
	t.Run("should be able to pass a viewer to observer", func(t *testing.T) {
		viewer := make(chan UrlStatus, 1)
		register := createMockUrlRegister()
		httpRequest := createMockHttpRequest()

		httpRequest.MockResult(&PingResult{
			Success:    true,
			StatusCode: 200,
		})

		obs := Observer{
			url:           "any-url.com",
			register:      register,
			httpRequester: httpRequest,
			running:       false,
			viewer:        viewer,
		}

		go obs.Run(t.Context(), time.Second*10)

		<-httpRequest.NotifyCh
		<-register.Chan

		select {
		case value := <-viewer:
			testutil.ExpectEqual(t, value.Url, "any-url.com")
			testutil.ExpectEqual(t, value.Up, true)
			testutil.ExpectEqual(t, value.StatusCode, 200)
		case <-time.After(time.Second * 1):
			t.Error("viewer was not notified")
		}
	})
}

func createMockUrlRegister() *testutil.MockRegister[*UrlStatus] {
	return testutil.CreateMockRegister[*UrlStatus]()
}

func createMockHttpRequest() *testutil.MockHttpRequester[*PingResult] {
	return testutil.CreateMockHttpRequester[*PingResult]()
}
