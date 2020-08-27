package helper

import (
	"crypto/rand"
	"errors"
	"github.com/fatih/color"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"io"
	mrand "math/rand"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// GeneratePassword generate password with bcrypt
func GeneratePassword(password string) string {
	encryptPassword, err := bcrypt.GenerateFromPassword([]byte(password), 0)
	if err != nil {
		return err.Error()
	}
	return string(encryptPassword)
}

// ComparePassword verify password with bcrypt
func ComparePassword(old string, new string) bool {
	return bcrypt.CompareHashAndPassword([]byte(old), []byte(new)) == nil
}

// String2Int64 helpful string convert to int64
func String2Int64(s string) uint64 {
	i, _ := strconv.ParseUint(s, 10, 64)
	return i
}

// Interface2Int64 ...
func Interface2Int64(s interface{}) uint64 {
	return uint64(s.(float64))
}

// IntContains ...
func IntContains(s []int, x int) bool {
	for _, n := range s {
		if x == n {
			return true
		}
	}
	return false
}

// StringContains ...
func StringContains(s []string, x string) bool {
	for _, n := range s {
		if x == n {
			return true
		}
	}
	return false
}

// IsSupportHTTPMethod ...
func IsSupportHTTPMethod(x string) bool {
	switch strings.ToUpper(x) {
	case "GET", "PUT", "POST", "DELETE":
		return true
	}
	return false
}

// IsNilInterface ...
func IsNilInterface(i interface{}) bool {
	return reflect.ValueOf(i).IsNil() == true
}

// BooleanHelper ...
func BooleanHelper(v bool) *bool {
	return &v
}

// TimeHelper ...
func TimeHelper(v time.Time) *time.Time {
	return &v
}

// ReplaceDateFormat ...
func ReplaceDateFormat(s string, r string) string {
	return strings.Replace(s, r, "", 1)
}

// ConvertUint ...
func ConvertUint(i string) (uint, error) {
	ok, err := strconv.ParseUint(i, 10, 64)
	if err != nil {
		return 0, errors.New("validate.param.parse.error")
	}
	return uint(ok), err
}

// GetSegment ...
func GetSegment(idx int, c *gin.Context) string {
	return strings.Split(c.Request.URL.Path, `/`)[idx+1]
}

// GetLastSegment ...
func GetLastSegment(c *gin.Context) string {
	a := strings.Split(c.Request.URL.Path, `/`)
	return a[len(a)-1]
}

// Info ...
func Info(v ...interface{}) {
	c := color.New(color.FgWhite)
	_, _ = c.Println(time.Now().Format(time.RFC3339), "[info]", v)
}

// Success ...
func Success(v ...interface{}) {
	c := color.New(color.FgHiGreen)
	_, _ = c.Println(time.Now().Format(time.RFC3339), "[success]", v)
}

// Error ...
func Error(v ...interface{}) {
	c := color.New(color.FgHiRed)
	_, _ = c.Println(time.Now().Format(time.RFC3339), "[error]", v)
}

// Panic ...
func Panic(v ...interface{}) {
	c := color.New(color.FgHiRed)
	_, _ = c.Println(time.Now().Format(time.RFC3339), "[error]", v)
	os.Exit(1)
}

var table = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9'}

// GenerateRandomDigitNumber ...
func GenerateRandomDigitNumber(max int) string {
	b := make([]byte, max)
	n, err := io.ReadAtLeast(rand.Reader, b, max)
	if n != max {
		panic(err)
	}
	for i := 0; i < len(b); i++ {
		b[i] = table[int(b[i])%len(table)]
	}
	return string(b)
}

//const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
var seededRand *mrand.Rand = mrand.New(mrand.NewSource(time.Now().UnixNano()))

// GenerateRandomStringWithCharset ...
func GenerateRandomStringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}
