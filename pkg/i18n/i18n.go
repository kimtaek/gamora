package i18n

import (
	"github.com/caarlos0/env"
	"github.com/gin-gonic/gin"
	"github.com/go-ini/ini"
	"github.com/kimtaek/gamora/constants/i18n"
	"golang.org/x/text/language"
	"strings"
)

// Configure config for i18n
type Configure struct {
	UserAgentTags    []string `env:"I18N_USERAGENT_TAGS" envDefault:"" envSeparator:":"` // lang-en:lang-ko:lang-ko:lang-zh:lang-zh:lang-zh:lang-jp:lang-jp
	SupportLanguages []language.Tag
}

// Config global defined i18n config
var Config Configure

// Setup init i18n config
func Setup() {
	_ = env.Parse(&Config)

	Config.SupportLanguages = []language.Tag{
		language.English,            // en
		language.Korean,             // ko
		language.MustParse("ko-kr"), // ko-kr
		language.Chinese,            // zh
		language.SimplifiedChinese,  // zh-Hans-CN
		language.TraditionalChinese, // zh-Hant-HK
		language.Japanese,           // ja
		language.MustParse("ja-JP"), // ja-JP
	}
}

// GetLanguage parse request language
func GetLanguage(ctx *gin.Context) string {
	lang := ctx.GetHeader("Accept-Language")
	if Config.UserAgentTags != nil {
		userAgent := ctx.GetHeader("User-Agent")
		for i, v := range Config.UserAgentTags {
			if strings.Contains(userAgent, v) {
				return Config.SupportLanguages[i].String()
			}
		}
	}

	if lang == "" {
		return "ko"
	}

	matcher := language.NewMatcher(Config.SupportLanguages)
	t, i := language.MatchStrings(matcher, lang)

	switch i {
	case 0:
		return "en"
	case 1, 2:
		return "ko"
	case 3, 4, 5:
		return "zh"
	case 6, 7:
		return "ja"
	}

	return t.String()
}

// GetI18nMessage return i18n message
func GetI18nMessage(code string, lang string) string {
	l := getLanguageFile(lang)
	if m := l.Section(ini.DefaultSection).Key(code).String(); m != "" {
		return m
	}
	return l.Section(ini.DefaultSection).Key(i18n.LangCodeNotFoundMessage).String()
}

// getLanguageFile load i18n language files
func getLanguageFile(lang string) *ini.File {
	if ok, err := ini.Load("i18n/" + lang + ".ini"); err == nil {
		return ok
	}
	if ok, err := ini.Load("i18n/en.ini"); err == nil {
		return ok
	}
	return nil
}
