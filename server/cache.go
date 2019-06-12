package server

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sub-app-server/config"
	"sync"
	"time"
)

type Cache struct {
	Url    string      `json:"url"`
	Hash   string      `json:"hash"`
	Header http.Header `json:"header"`
	Path   string      `json:"path"`
	Expire time.Time   `json:"expire"`

	body []byte `json:"-"`
}

type CacheStore struct {
	Data  map[string]*Cache `json:"data"`
	mutex *sync.Mutex
}

func (c *CacheStore) store(cache *Cache) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.Data[cache.Hash] = cache

	if data, err := json.Marshal(c); err == nil {

		dirname := filepath.Dir(cache.Path)

		if _, err := os.Stat(dirname); os.IsNotExist(err) {
			_ = os.MkdirAll(dirname, os.ModePerm)
		}
		//缓存文件
		_ = ioutil.WriteFile(cache.Path, cache.body, 0666)

		//缓存对象
		_ = ioutil.WriteFile(config.C.Cache.Dir+"/meta.json", data, 0666)
	}
}

func (c *CacheStore) load(key string) *Cache {

	c.mutex.Lock()
	defer func() { c.mutex.Unlock() }()

	if ca, ok := c.Data[key]; ok {
		return ca
	}
	return nil
}

var store = CacheStore{
	Data: map[string]*Cache{},
}

func init() {

	if data, err := ioutil.ReadFile(config.C.Cache.Dir + "/meta.json"); err == nil {
		_ = json.Unmarshal(data, &store)
	}

	store.mutex = &sync.Mutex{}
}

func hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

func contentType(header http.Header) string {
	var typ string

	if len(header["content-type"]) > 0 {
		typ = header["content-type"][0]
	} else if len(header["Content-Type"]) > 0 {
		typ = header["Content-Type"][0]
	}

	return typ
}

func cacheResponse(req *http.Request, resp HttpResponse) bool {

	if resp.Body != nil {

		//跳过缓存
		for _, url := range config.C.Cache.Ignore {
			if strings.HasPrefix(req.URL.String(), url) {
				return false
			}
		}

		uri := req.URL

		//url := uri.Scheme + "://" + uri.Host + uri.Path
		url := uri.String()

		urlHash := hash(url)

		duration, err := time.ParseDuration(config.C.Cache.Expire)

		if err != nil {
			log.Println("parse cache expire failure !", config.C.Cache.Expire)
			return false
		}

		cache := &Cache{
			Url:    req.URL.String(),
			Hash:   urlHash,
			Header: resp.Headers,
			Path:   fmt.Sprintf(`%s/%s/%s`, config.C.Cache.Dir, req.Host, urlHash),
			Expire: time.Now().Add(duration),
			body:   resp.Body,
		}

		//根据content type缓存
		typ := contentType(resp.Headers)

		for _, val := range config.C.Cache.Types {

			if strings.HasPrefix(typ, val) {
				store.store(cache)
				return true
			}
		}

		//根据url缓存
		for _, val := range config.C.Cache.Urls {
			if strings.HasPrefix(url, val) {
				store.store(cache)
				return true
			}
		}

	}

	return false
}

func loadCache(req *http.Request) *http.Response {

	key := hash(req.URL.String())

	cache := store.load(key)

	if cache != nil {

		if time.Now().Before(cache.Expire) {

			if data, err := ioutil.ReadFile(cache.Path); err == nil {
				resp := HttpResponse{
					Headers: cache.Header,
					Code:    200,
					Body:    data,
				}
				return makeResponse(req, resp)
			}
		}
	}
	return nil
}
