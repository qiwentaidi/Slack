package urlutil

import (
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUTF8URLEncoding(t *testing.T) {
	exstring := '上'
	expected := `e4%b8%8a`
	val := getutf8hex(exstring)
	require.Equalf(t, val, expected, "failed to url encode utf char expected %v but got %v", expected, val)
}

func TestParamEncoding(t *testing.T) {
	testcases := []struct {
		Payload  string
		Expected string
	}{
		{"1+AND+(SELECT+*+FROM+(SELECT(SLEEP(12)))nQIP)", "1+AND+(SELECT+*+FROM+(SELECT(SLEEP(12)))nQIP)"},
		{"1 AND SELECT", "1+AND+SELECT"},
	}
	for _, v := range testcases {
		val := ParamEncode(v.Payload)
		require.Equalf(t, val, v.Expected, "failed to url encode payload expected %v got %v", v.Expected, val)
	}
}

func TestRawParam(t *testing.T) {
	p := NewParams()
	p.Add("sqli", "1+AND+(SELECT+*+FROM+(SELECT(SLEEP(12)))nQIP)")
	p.Add("xss", "<script>alert('XSS')</script>")
	p.Add("xssiwthspace", "<svg id=alert(1) onload=eval(id)>")
	p.Add("jsprotocol", "javascript://alert(1)")
	// Note keys are sorted
	expected := "jsprotocol=javascript://alert(1)&sqli=1+AND+(SELECT+*+FROM+(SELECT(SLEEP(12)))nQIP)&xss=<script>alert('XSS')</script>&xssiwthspace=<svg+id=alert(1)+onload=eval(id)>"
	require.Equalf(t, p.Encode(), expected, "failed to encode parameters expected %v but got %v", expected, p.Encode())
}

func TestParamIntegration(t *testing.T) {
	var routerErr error
	expected := "/params?jsprotocol=javascript://alert(1)&sqli=1+AND+(SELECT+*+FROM+(SELECT(SLEEP(12)))nQIP)&xss=<script>alert('XSS')</script>&xssiwthspace=<svg+id=alert(1)+onload=eval(id)>"

	http.HandleFunc("/params", func(w http.ResponseWriter, r *http.Request) {
		if r.RequestURI != expected {
			routerErr = fmt.Errorf("expected %v but got %v", expected, r.RequestURI)
		}
		w.WriteHeader(http.StatusOK)
	})
	//nolint:all
	go http.ListenAndServe(":9000", nil)

	p := NewParams()
	p.Add("sqli", "1+AND+(SELECT+*+FROM+(SELECT(SLEEP(12)))nQIP)")
	p.Add("xss", "<script>alert('XSS')</script>")
	p.Add("xssiwthspace", "<svg id=alert(1) onload=eval(id)>")
	p.Add("jsprotocol", "javascript://alert(1)")

	url, _ := url.Parse("http://localhost:9000/params")
	url.RawQuery = p.Encode()
	_, err := http.Get(url.String())
	require.Nil(t, err)
	require.Nil(t, routerErr)
}

func TestPercentEncoding(t *testing.T) {
	// From Burpsuite
	expected := "%74%65%73%74%20%26%23%20%28%29%20%70%65%72%63%65%6E%74%20%5B%5D%7C%2A%20%65%6E%63%6F%64%69%6E%67"
	payload := "test &# () percent []|* encoding"
	value := PercentEncoding(payload)
	require.Equalf(t, value, expected, "expected percentencoding to be %v but got %v", expected, value)
	decoded, err := url.QueryUnescape(value)
	require.Nil(t, err)
	require.Equal(t, payload, decoded)
}

func TestGetParams(t *testing.T) {
	values := url.Values{}
	values.Add("sqli", "1+AND+(SELECT+*+FROM+(SELECT(SLEEP(12)))nQIP)")
	values.Add("xss", "<script>alert('XSS')</script>")
	p := GetParams(values)
	require.NotNilf(t, p, "expected params but got nil")
	require.Equalf(t, p.Get("sqli"), values.Get("sqli"), "malformed or missing value for param sqli expected %v but got %v", values.Get("sqli"), p.Get("sqli"))
	require.Equalf(t, p.Get("xss"), values.Get("xss"), "malformed or missing value for param xss expected %v but got %v", values.Get("xss"), p.Get("xss"))
}

func TestURLEncode(t *testing.T) {
	example := "\r\n"
	got := URLEncodeWithEscapes(example)
	require.Equalf(t, "%0D%0A", got, "failed to url encode characters")

	// verify with stdlib
	for r := 0; r < 20; r++ {
		expected := url.QueryEscape(string(rune(r)))
		got := URLEncodeWithEscapes(string(rune(r)))
		require.Equalf(t, expected, got, "url encoding mismatch for non-printable char with ascii val:%v", r)
	}
}
