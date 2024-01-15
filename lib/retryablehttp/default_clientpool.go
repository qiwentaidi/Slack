package retryablehttp

var (
	DefaultClientPool  *Client
	RedirectClientPool *Client
)

var options *PoolOptions

func init() {
	options = &DefaultPoolOptions
	InitClientPool(options)

	DefaultClientPool, _ = GetPool(options)

	options.EnableRedirect(FollowAllRedirect)
	RedirectClientPool, _ = GetPool(options)
}
