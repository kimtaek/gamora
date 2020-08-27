package respond

import (
	"github.com/gin-gonic/gin"
	"github.com/kimtaek/gamora/pkg/i18n"
	"net/http"
)

// Source defined response source
type Source struct {
	Code        int         `json:"-"`
	Data        interface{} `json:"data"`
	Message     string      `json:"message"`
	MessageCode string      `json:"-"`
}

// Data response data
func Data(c *gin.Context, s *Source) {
	if s.Code == 0 {
		s.Code = http.StatusOK
	}

	if s.MessageCode != "" {
		s.Message = i18n.GetI18nMessage(s.MessageCode, i18n.GetLanguage(c))
	}

	if s.Message == "" {
		s.Message = http.StatusText(s.Code)
	}

	c.AbortWithStatusJSON(s.Code, gin.H{
		"data":    s.Data,
		"message": s.Message,
	})

	return
}

// MessageByCode response by code
func MessageByCode(c *gin.Context, code string) {
	c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
		"data":    nil,
		"message": i18n.GetI18nMessage(code, i18n.GetLanguage(c)),
	})

	return
}

// MessageWithStatusByCode response with special status code
func MessageWithStatusByCode(c *gin.Context, code string, status int) {
	if status == 0 {
		status = http.StatusBadRequest
	}

	c.AbortWithStatusJSON(status, gin.H{
		"data":    nil,
		"message": i18n.GetI18nMessage(code, i18n.GetLanguage(c)),
	})

	return
}

// Message response message with 400 code
func Message(c *gin.Context, message string) {
	c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
		"data":    nil,
		"message": message,
	})

	return
}

// MessageWithStatus response message with special status code
func MessageWithStatus(c *gin.Context, message string, status int) {
	if status == 0 {
		status = http.StatusBadRequest
	}

	c.AbortWithStatusJSON(status, gin.H{
		"data":    nil,
		"message": message,
	})

	return
}
