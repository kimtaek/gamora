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
		println(GetLanguage(c))
	})
	req, _ := http.NewRequest(http.MethodGet, "/", nil)

	// IE: ko, FF: ko-kr,ko;q=0.8,en-us;q=0.5,en;q=0.3
	// IE: ja-JP, FF: ja,en-US;q=0.7,en;q=0.3
	w := httptest.NewRecorder()

	req.Header.Add("Accept-Language", "ja-JP")
	r.ServeHTTP(w, req)

	req, _ = http.NewRequest(http.MethodGet, "/", nil)
	req.Header.Add("Accept-Language", "ja,en-US;q=0.7,en;q=0.3")
	r.ServeHTTP(w, req)

	req, _ = http.NewRequest(http.MethodGet, "/", nil)
	req.Header.Add("Accept-Language", "ko")
	r.ServeHTTP(w, req)

	req, _ = http.NewRequest(http.MethodGet, "/", nil)
	req.Header.Add("Accept-Language", "ko-kr,ko;q=0.8,en-us;q=0.5,en;q=0.3")
	r.ServeHTTP(w, req)

	req, _ = http.NewRequest(http.MethodGet, "/", nil)
	r.ServeHTTP(w, req)

	req, _ = http.NewRequest(http.MethodGet, "/", nil)
	req.Header.Add("Accept-Language", "en")
	r.ServeHTTP(w, req)
}

func TestUserAgentTags(t *testing.T) {
	Setup()
	gin.SetMode(gin.TestMode)

	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		println(GetLanguage(c))
	})
	req, _ := http.NewRequest(http.MethodGet, "/", nil)

	// IE: ko, FF: ko-kr,ko;q=0.8,en-us;q=0.5,en;q=0.3
	// IE: ja-JP, FF: ja,en-US;q=0.7,en;q=0.3
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
	req.Header.Add("User-Agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 14_0_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148/lang-jp")
	r.ServeHTTP(w, req)

	req, _ = http.NewRequest(http.MethodGet, "/", nil)
	req.Header.Add("User-Agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 14_0_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148/lang-unknown")
	r.ServeHTTP(w, req)
}
