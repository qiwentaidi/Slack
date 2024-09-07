package crack

import (
	"encoding/base64"
	"io"
	"net/http"
	"net/url"
)

const defaultApiURL = "http://localhost:8000/ocr"

func DdddOcr(imageData []byte) (string, error) {
	base64Image := base64.StdEncoding.EncodeToString(imageData)
	data := url.Values{}
	data.Set("image", base64Image)
	data.Set("probability", "false")
	data.Set("png_fix", "false")

	resp, err := http.PostForm(defaultApiURL, data)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}
