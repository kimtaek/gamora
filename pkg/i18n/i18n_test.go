package i18n

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAcceptLanguages(t *testing.T) {
	Setup()
	Config.UserAgentTags = nil
	gin.SetMode(gin.TestMode)

	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		println(c.Request.Header.Get("Accept-Language"), " ==> ", GetLanguage(c))
	})
	req, _ := http.NewRequest(http.MethodGet, "/", nil)

	// IE: ko, FF: ko-kr,ko;q=0.8,en-us;q=0.5,en;q=0.3
	// IE: ja-JP, FF: ja,en-US;q=0.7,en;q=0.3
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	req.Header.Add("Accept-Language", "ja-JP,ja;q=0.8,en-us;q=0.5,en;q=0.3")
	r.ServeHTTP(w, req)

	req, _ = http.NewRequest(http.MethodGet, "/", nil)
	req.Header.Add("Accept-Language", "en-US,en;q=0.5")
	r.ServeHTTP(w, req)

	req, _ = http.NewRequest(http.MethodGet, "/", nil)
	req.Header.Add("Accept-Language", "en,zh;q=0.9,zh-TW;q=0.8,zh-HK;q=0.7,zh-CN;q=0.6,ko;q=0.5")
	r.ServeHTTP(w, req)

	req, _ = http.NewRequest(http.MethodGet, "/", nil)
	req.Header.Add("Accept-Language", "ko-KR,ko;q=0.9,en-US;q=0.8,en;q=0.7")
	r.ServeHTTP(w, req)

	req, _ = http.NewRequest(http.MethodGet, "/", nil)
	req.Header.Add("Accept-Language", "ko,en;q=0.9,zh;q=0.8,zh-TW;q=0.7,zh-HK;q=0.6,zh-CN;q=0.5")
	r.ServeHTTP(w, req)

	req, _ = http.NewRequest(http.MethodGet, "/", nil)
	r.ServeHTTP(w, req)

	req, _ = http.NewRequest(http.MethodGet, "/", nil)
	req.Header.Add("Accept-Language", "zh,zh-TW;q=0.9,zh-HK;q=0.8,zh-CN;q=0.7,ko;q=0.6,en;q=0.5")
	r.ServeHTTP(w, req)

	req, _ = http.NewRequest(http.MethodGet, "/", nil)
	req.Header.Add("Accept-Language", "zh-CN,zh;q=0.9,ko;q=0.8,en;q=0.7")
	r.ServeHTTP(w, req)

	req, _ = http.NewRequest(http.MethodGet, "/", nil)
	req.Header.Add("Accept-Language", "zh-HK,zh-CN;q=0.9,zh;q=0.8,ko;q=0.7,en;q=0.6,zh-TW;q=0.5")
	r.ServeHTTP(w, req)

	req, _ = http.NewRequest(http.MethodGet, "/", nil)
	req.Header.Add("Accept-Language", "zh-TW,zh-HK;q=0.9,zh-CN;q=0.8,zh;q=0.7,ko;q=0.6,en;q=0.5")
	r.ServeHTTP(w, req)
}

func TestUserAgentTags(t *testing.T) {
	Setup()
	gin.SetMode(gin.TestMode)

	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		println(c.Request.Header.Get("User-Agent"), " ==> ", GetLanguage(c))
	})
	req, _ := http.NewRequest(http.MethodGet, "/", nil)

	w := httptest.NewRecorder()
	req.Header.Add("User-Agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 14_0_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148/lang-en")
	r.ServeHTTP(w, req)

	req, _ = http.NewRequest(http.MethodGet, "/", nil)
	req.Header.Add("User-Agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 14_0_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148/lang-ko")
	r.ServeHTTP(w, req)

	req, _ = http.NewRequest(http.MethodGet, "/", nil)
	req.Header.Add("User-Agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 14_0_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148/lang-zh")
	r.ServeHTTP(w, req)

	req, _ = http.NewRequest(http.MethodGet, "/", nil)
	req.Header.Add("User-Agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 14_0_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148/lang-ja")
	r.ServeHTTP(w, req)

	req, _ = http.NewRequest(http.MethodGet, "/", nil)
	req.Header.Add("User-Agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 14_0_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148/lang-unknown")
	r.ServeHTTP(w, req)
}
