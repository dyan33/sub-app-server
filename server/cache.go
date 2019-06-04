package server

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

const CACHEDIR = "_cache"

const (
	STREAM = "application/octet-stream"
	FONT   = "application/font-woff"
)

var caches []string

func init() {

	caches = []string{
		"text/css",
		"image/png",
		"image/gif",
		"mage/gif",
		"image/jpeg",
		STREAM,
		FONT,
	}
}

func hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

func saveFile(typ, url string, data []byte) {

	dirname := fmt.Sprintf(`%s/%s`, CACHEDIR, typ)

	if _, err := os.Stat(dirname); os.IsNotExist(err) {
		_ = os.MkdirAll(dirname, os.ModePerm)
	}

	filename := fmt.Sprintf("%s/%s", dirname, hash(url))

	if err := ioutil.WriteFile(filename, data, 0666); err != nil {
		fmt.Println("缓存失败", err)
	}
}

func readFile(typ, url string) []byte {

	filename := fmt.Sprintf(`%s/%s/%s`, CACHEDIR, typ, hash(url))

	if _, err := os.Stat(filename); err == nil {
		if data, err := ioutil.ReadFile(filename); err == nil {
			return data
		}
	}
	return nil
}

func cacheResponse(req *http.Request, resp HttpResponse) {

	if resp.Body != nil {

		var typ string

		if len(resp.Headers["content-type"]) > 0 {
			typ = resp.Headers["content-type"][0]
		} else if len(resp.Headers["Content-Type"]) > 0 {
			typ = resp.Headers["Content-Type"][0]
		}

		if len(typ) > 0 {

			url := req.URL.String()

			if typ == STREAM || typ == FONT {
				if strings.HasSuffix(url, ".woff") {
					saveFile(typ, url, resp.Body)
				}
			} else {
				for _, key := range caches {

					if key == typ {
						saveFile(typ, url, resp.Body)
					}
				}
			}
		}
	}
}

func loadCache(req *http.Request) *http.Response {

	var data []byte
	var typ string

	url := req.URL.String()

	for _, typ = range caches {
		data = readFile(typ, url)
		if data != nil {
			break
		}
	}

	if data != nil {
		resp := HttpResponse{
			Headers: map[string][]string{
				"Content-Type": {typ},
			},
			Code: 200,
			Body: data,
		}
		return makeResponse(req, resp)
	}
	return nil
}
