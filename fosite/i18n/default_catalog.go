// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package i18n

import (
	"net/http"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

// DefaultMessage is a single message in the locale bundle
// identified by 'ID'.
type DefaultMessage struct {
	ID               string `json:"id"`
	FormattedMessage string `json:"msg"`
}

// DefaultLocaleBundle is a bundle of messages for the specified
// locale. The language tag can be arbitrary to allow for
// unsupported/unknown languages used by custom clients.
type DefaultLocaleBundle struct {
	LangTag  string            `json:"lang"`
	Messages []*DefaultMessage `json:"messages"`
}

// defaultMessageCatalog is a catalog of all locale bundles.
type defaultMessageCatalog struct {
	Bundles []*DefaultLocaleBundle

	matcher language.Matcher
}

func NewDefaultMessageCatalog(bundles []*DefaultLocaleBundle) MessageCatalog {
	c := &defaultMessageCatalog{
		Bundles: bundles,
	}

	for _, v := range c.Bundles {
		if err := v.Init(); err != nil {
			continue
		}
	}

	c.makeMatcher()
	return c
}

// Init initializes the default catalog with the
// list of messages. The lang tag must parse, otherwise this
// func will panic.
func (l *DefaultLocaleBundle) Init() error {
	tag := language.MustParse(l.LangTag)
	for _, m := range l.Messages {
		if err := message.SetString(tag, m.ID, m.FormattedMessage); err != nil {
			return err
		}
	}

	return nil
}

func (c *defaultMessageCatalog) GetMessage(ID string, tag language.Tag, v ...interface{}) string {
	matchedTag, _, _ := c.matcher.Match(tag)
	p := message.NewPrinter(matchedTag)

	result := p.Sprintf(ID, v...)
	if result == ID && tag != language.English {
		return c.GetMessage(ID, language.English, v...)
	}

	return result
}

func (c *defaultMessageCatalog) GetLangFromRequest(r *http.Request) language.Tag {
	lang, _ := r.Cookie("lang")
	accept := r.Header.Get("Accept-Language")
	tag, _ := language.MatchStrings(c.matcher, lang.String(), accept)

	return tag
}

func (c *defaultMessageCatalog) makeMatcher() {
	result := []language.Tag{language.English}
	defLangs := message.DefaultCatalog.Languages()
	// remove "en" if was already in the list of languages
	for i, t := range defLangs {
		if t == language.English {
			result = append(result, defLangs[:i]...)
			result = append(result, defLangs[i+1:]...)
		}
	}

	c.matcher = language.NewMatcher(defLangs)
}
