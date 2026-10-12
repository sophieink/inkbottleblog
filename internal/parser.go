package internal

import (
	"bytes"
	"fmt"
	"html/template"
	"regexp"
	"strings"
	"time"

	"github.com/yuin/goldmark/v2/extension"
	"github.com/yuin/goldmark/v2/parser"
	"github.com/yuin/goldmark/v2/renderer/html"
	highlighting "github.com/yuin/goldmark-highlighting/v3"
)

type KV struct {
	Key   string
	Value string
}

type Post struct {
	Slug    string
	Title   string
	Date    string
	KVs     []KV
	Content template.HTML
}

// parser：GFM 全家桶（表格、删除线、任务列表、链接自动识别）
var p = parser.New(
	parser.WithExtensions(
		extension.GFMParser,
		highlighting.NewParser(),
	),
)

// renderer：按需注册，和 parser 里用到的 GFM 功能一一对应
var r = html.New(
	html.WithExtensions(
		extension.NewTableHTMLRenderer(),          // 表格
		extension.NewStrikethroughHTMLRenderer(),  // 删除线
		extension.NewTaskListItemHTMLRenderer(),   // 任务列表
		highlighting.NewHTMLRenderer(
			highlighting.WithStyle("monokailight"),
		),
	),
	html.WithUnsafe(),
	html.WithXHTML(),
)

var kvRe = regexp.MustCompile(`[\[\(\{<]\s*([^:\[\]\(\)\{\}<>\s]+)\s*::\s*([^\[\]\(\)\{\}<>]*?)\s*[\]\)\}>]`)

func ParseMarkdown(raw []byte, slug string) (*Post, error) {
	text := string(raw)
	post := &Post{Slug: slug}

	// 提取KV + 替换正文
	text = kvRe.ReplaceAllStringFunc(text, func(m string) string {
		sub := kvRe.FindStringSubmatch(m)
		if len(sub) != 3 {
			return m
		}
		k := strings.TrimSpace(sub[1])
		v := strings.TrimSpace(sub[2])
		post.KVs = append(post.KVs, KV{Key: k, Value: v})

		//title 正文不显示,只用于页面标题
		if k == "title" {
			return ""
		}
		//(k:: v)--> v
		if m[0] == '(' {
			return fmt.Sprintf(
				`<span class="kv"><span class="v">%s</span></span>`,
				template.HTMLEscapeString(v),
			)
		}

		return fmt.Sprintf(
			`<span class="kv"><span class="k">%s</span><span class="sep">::</span><span class="v">%s</span></span>`,
			template.HTMLEscapeString(k),
			template.HTMLEscapeString(v),
		)
	})

	// 渲染md -> html
	source := []byte(text)
	doc := p.Parse(source)
	var buf bytes.Buffer
	if err := r.Render(&buf, source, doc); err != nil {
		return nil, fmt.Errorf("render %s: %w", slug, err)
	}

	post.Content = template.HTML(buf.String())
	post.Title = slug
	for _, kv := range post.KVs {
		if kv.Key == "title" {
			post.Title = kv.Value	
		}
		if kv.Key == "date" {
			post.Date = kv.Value
			break
		}
	}
	return post, nil
}

// 日期合RFC 1123，支持所有RSS阅读器
func (p *Post) PubDate() string {
	t, err := time.Parse("2006-01-02", p.Date)
	if err != nil {
		return time.Now().UTC().Format(time.RFC1123Z)
	}
	return t.UTC().Format(time.RFC1123Z)
}
