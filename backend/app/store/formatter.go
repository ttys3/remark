package store

import (
	"bytes"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	chromahtml "github.com/alecthomas/chroma/formatters/html"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

// CommentFormatter implements all generic formatting ops on comment
type CommentFormatter struct {
	colorScheme string // chroma style name
	converters []CommentConverter
}

// CommentConverter defines interface to convert some parts of commentHTML
// Passed at creation time and does client-defined conversions, like image proxy link change
type CommentConverter interface {
	Convert(text string) string
}

// CommentConverterFunc functional struct implementing CommentConverter
type CommentConverterFunc func(text string) string

// Convert calls func for given text
func (f CommentConverterFunc) Convert(text string) string {
	return f(text)
}

// NewCommentFormatter makes CommentFormatter
func NewCommentFormatter(converters ...CommentConverter) *CommentFormatter {
	return &CommentFormatter{colorScheme: "monokailight", converters: converters}
}

// Set chroma style name
func (f *CommentFormatter) WithStyle(colorScheme string) *CommentFormatter {
	f.colorScheme = colorScheme
	return f
}

// Format comment fields
func (f *CommentFormatter) Format(c Comment) Comment {
	c.Text = f.FormatText(c.Text)
	return c
}

// FormatText converts text with markdown processor, applies external converters and shortens links
func (f *CommentFormatter) FormatText(txt string) (res string) {
	if f.colorScheme == "" {
		f.colorScheme = "monokailight"
	}

	markdown := goldmark.New(
		goldmark.WithExtensions(extension.GFM), // Linkify | Table | Strikethrough | TaskList
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
		),
		goldmark.WithExtensions(
			highlighting.NewHighlighting(
				highlighting.WithStyle(f.colorScheme),
				highlighting.WithFormatOptions(
					chromahtml.TabWidth(2),
					chromahtml.WithClasses(true),
				),
			),
		),
		goldmark.WithExtensions(extension.Typographer), // substitutes punctuations with typographic entities like smartypants
	)

	var output bytes.Buffer
	if err := markdown.Convert([]byte(txt), &output); err != nil {
		panic(err)
	}
	res = f.unEscape(output.String())

	for _, conv := range f.converters {
		res = conv.Convert(res)
	}
	res = f.shortenAutoLinks(res, shortURLLen)
	return res
}

// Shortens all the automatic links in HTML: auto link has equal "href" and "text" attributes.
func (f *CommentFormatter) shortenAutoLinks(commentHTML string, max int) (resHTML string) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(commentHTML))
	if err != nil {
		return commentHTML
	}
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		if href, ok := s.Attr("href"); ok {
			if href != s.Text() || len(href) < max+3 || max < 3 {
				return
			}
			commentURL, e := url.Parse(href)
			if e != nil {
				return
			}
			commentURL.Path, commentURL.RawQuery, commentURL.Fragment = "", "", ""
			host := commentURL.String()
			if host == "" {
				return
			}
			short := href[:max-3]
			if len(short) < len(host) {
				short = host
			}
			s.SetText(short + "...")
		}
	})
	resHTML, err = doc.Find("body").Html()
	if err != nil {
		return commentHTML
	}
	return resHTML
}

func (f *CommentFormatter) unEscape(txt string) (res string) {
	elems := []struct {
		from, to string
	}{
		{`&amp;mdash;`, "—"},
	}
	res = txt
	for _, e := range elems {
		res = strings.Replace(res, e.from, e.to, -1)
	}
	return res
}
