package retryablehttp_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zan8in/retryablehttp"
)

// This test is just to make sure that the default client is initialized
// correctly.
func Test_DefaultHttpClient(t *testing.T) {
	require.NotNil(t, retryablehttp.DefaultHTTPClient)
	resp, err := retryablehttp.DefaultHTTPClient.Get("http://example.com")
	require.Nil(t, err)
	require.NotNil(t, resp)
}
