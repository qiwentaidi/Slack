// 针对 amazonaws 存储桶数据泄露利用
package listbucket

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/url"
	"slack-wails/lib/clients"
)

// 最大分页
const MaxLimit = 100000000

// 使用二分法获取到亚马逊云的桶数量总数
// ?max-keys=1000&delimiter=1 可以通过delimiter参数进行分页

// CheckBucketCount uses a binary search to determine the total number of pages in the bucket
// func CheckBucketCount(bucketURL string) int {
// 	low, high := 1, MaxLimit
// 	totalPages := 0
// 	// 最终的内容hash
// 	prevHash := ""
// 	currentHash, _ := getBucketHash(bucketURL, high)

// 	for {
// 		mid := (low + high) / 2
// 		currentHash, success := getBucketHash(bucketURL, mid)
// 		if !success {
// 			break
// 		}

// 		if prevHash != "" && currentHash == prevHash {
// 			// No new data, adjust the high bound
// 			high = mid - 1
// 		} else {
// 			// New data found, adjust the low bound
// 			low = mid + 1
// 			totalPages = mid
// 			prevHash = currentHash
// 		}
// 	}

// 	return totalPages
// }

func getBucketHash(bucketURL string, delimiter int) (string, bool) {
	// Construct request URL
	params := url.Values{}
	params.Set("max-keys", "1000")
	params.Set("delimiter", fmt.Sprintf("%d", delimiter))

	reqURL := fmt.Sprintf("%s?%s", bucketURL, params.Encode())
	resp, body, err := clients.NewSimpleGetRequest(reqURL, clients.NewHttpClient(nil, true))
	if err != nil || resp == nil {
		return "", false
	}
	// Compute hash of the response
	hash := md5.Sum(body)
	return hex.EncodeToString(hash[:]), true
}
