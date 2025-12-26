package markdown

import (
	"bytes"

	"github.com/microcosm-cc/bluemonday"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer/html"
)

var md = goldmark.New(
	goldmark.WithExtensions(extension.GFM),
	goldmark.WithRendererOptions(
		html.WithUnsafe(), // 先允许，后面用 bluemonday 过滤
	),
)

var policy = bluemonday.UGCPolicy()

func RenderToSafeHTML(markdownText string) (string, error) {
	var buf bytes.Buffer
	if err := md.Convert([]byte(markdownText), &buf); err != nil {
		return "", err
	}
	// XSS 清洗
	return policy.Sanitize(buf.String()), nil
}
