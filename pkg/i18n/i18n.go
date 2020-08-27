package i18n

import (
	"github.com/caarlos0/env"
	"github.com/gin-gonic/gin"
	"github.com/go-ini/ini"
	"github.com/kimtaek/gamora/constants/i18n"
	"golang.org/x/text/language"
)

type config struct {
	SupportLanguages []language.Tag
}

var c config

func Setup() {
	_ = env.Parse(&c)

	c.SupportLanguages = []language.Tag{
		language.English,            // en
		language.Korean,             // ko
		language.MustParse("ko-Kr"), //ko-Kr
		language.SimplifiedChinese,  // zh-Hans
	}
}

func GetLanguage(ctx *gin.Context) string {
	accept := ctx.GetHeader("Accept-Language")
	if accept == "" {
		return "ko"
	}

	matcher := language.NewMatcher(c.SupportLanguages)
	t, i := language.MatchStrings(matcher, accept)

	switch i {
	case 0:
		return "en"
	case 1, 2:
		return "ko"
	}

	return t.String()
}

func GetI18nMessage(code string, lang string) string {
	l := getLanguageFile(lang)
	if m := l.Section(ini.DefaultSection).Key(code).String(); m != "" {
		return m
	}
	return l.Section(ini.DefaultSection).Key(i18n.LangCodeNotFoundMessage).String()
}

func getLanguageFile(lang string) *ini.File {
	if ok, err := ini.Load("i18n/" + lang + ".ini"); err == nil {
		return ok
	}
	if ok, err := ini.Load("i18n/en.ini"); err == nil {
		return ok
	}
	return nil
}
