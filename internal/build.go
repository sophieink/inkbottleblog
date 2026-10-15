package internal

import (
	"bytes"
	htmlTpl "html/template"
	textTpl "text/template"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"inkbottle/conf"
)

// Build 渲染 posts/*.md → docs/*.html
func Build() error {
	// 1. 解析所有 md
	entries, err := os.ReadDir(conf.PostsDir)
	if err != nil {
		return err
	}

	var posts []*Post
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".md") {
			continue
		}
		slug := strings.TrimSuffix(e.Name(), ".md")
		raw, err := ReadFile(filepath.Join(conf.PostsDir, e.Name()))
		if err != nil {
			return err
		}
		p, err := ParseMarkdown(raw, slug)
		if err != nil {
			return err
		}
		posts = append(posts, p)
	}

	// 排序：最新在前
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].Date > posts[j].Date
	})

	// 排除 about 的列表（首页用）
	var listed []*Post
	for _, p := range posts {
		if p.Slug == "about" {
			continue
		}
		listed = append(listed, p)
	}

	// 2. 渲染首页（独立解析 index 集合）
	indexTpl, err := htmlTpl.ParseFiles(
		filepath.Join(conf.TemplateDir, "base.html"),
		filepath.Join(conf.TemplateDir, "partials", "header.html"),
		filepath.Join(conf.TemplateDir, "partials", "footer.html"),
		filepath.Join(conf.TemplateDir, "index.html"),
	)
	if err != nil {
		return err
	}
	data := struct {
		Title string
		Posts []*Post
	}{Title: "Blog", Posts: listed}
	if err := render(indexTpl, "index.html", "index.html", data); err != nil {
		return err
	}

	// 3. 渲染每篇文章（独立解析 post 集合）
	postTpl, err := htmlTpl.ParseFiles(
		filepath.Join(conf.TemplateDir, "base.html"),
		filepath.Join(conf.TemplateDir, "partials", "header.html"),
		filepath.Join(conf.TemplateDir, "partials", "footer.html"),
		filepath.Join(conf.TemplateDir, "post.html"),
	)
	if err != nil {
		return err
	}
	for _, p := range posts {
		pdata := struct {
			Title string
			Post  *Post
		}{Title: p.Title, Post: p}
		if err := render(postTpl, "post.html", p.Slug+".html", pdata); err != nil {
			return err
		}
	}

	// 4. 渲染 feed.xml（text/template，独立）
	feedTpl, err := textTpl.ParseFiles(filepath.Join(conf.TemplateDir, "feed.xml"))
	if err != nil {
		return err
	}
	feedData := struct {
		SiteTitle string
		SiteURL   string
		Posts     []*Post
	}{
		SiteTitle: "Blog",
		SiteURL:   "https://sophieink.github.io/inkbottle-blog/index.html",
		Posts:     listed,
	}
	var feedBuf bytes.Buffer
	if err := feedTpl.ExecuteTemplate(&feedBuf, "feed.xml", feedData); err != nil {
		return err
	}
	if err := WriteFile(filepath.Join(conf.DocsDir, "feed.xml"), feedBuf.Bytes()); err != nil {
		return err
	}

	// 5. 拷贝静态资源
	return CopyDir(conf.StaticDir, conf.DocsStaticDir)
}

// render 渲染模板到 docs/<out>
func render(tpl *htmlTpl.Template, name, out string, data any) error {
	var buf bytes.Buffer
	if err := tpl.ExecuteTemplate(&buf, name, data); err != nil {
		return err
	}
	return WriteFile(filepath.Join(conf.DocsDir, out), buf.Bytes())
}
