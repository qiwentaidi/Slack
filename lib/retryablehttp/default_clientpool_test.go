package retryablehttp_test

import (
	"testing"

	"github.com/zan8in/retryablehttp"
)

func TestGet2(t *testing.T) {
	for i := 0; i < 10; i++ {
		resp2, err := retryablehttp.RedirectClientPool.Get("http://127.0.0.1/redirect.php")
		t.Log(resp2, err)
	}
}
