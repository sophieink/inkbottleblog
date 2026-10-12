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

	// 排序==[date:: YYYY-MM-DD] 最新post在前
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].Date > posts[j].Date //降序，最新前置
	})

	// 2. 两组html模板：layout+index / layout+post
	// 组 1：layout + index
  indexTpl, err := htmlTpl.ParseFiles(
	  filepath.Join(conf.TemplateDir, "layout.html"),
	  filepath.Join(conf.TemplateDir, "index.html"),   // ← 必须一起
  )

  // 组 2：layout + post
  postTpl, err := htmlTpl.ParseFiles(
	  filepath.Join(conf.TemplateDir, "layout.html"),
	  filepath.Join(conf.TemplateDir, "post.html"),    // ← 必须一起
  )

	// 3. 渲染每篇（用 postTpl）
	for _, p := range posts {
		data := struct {
			Title string
			Post  *Post
		}{Title: p.Title, Post: p}
		if err := render(postTpl, "layout", p.Slug+".html", data); err != nil {
			return err
		}
	}

	// 4. 渲染首页（排除 about）
	var listed []*Post
	for _, p := range posts {
		if p.Slug == "about" {
			continue
		}
		listed = append(listed, p)
	}
	data := struct {
		Title string
		Posts []*Post
	}{Title: "Blog", Posts: listed}
	if err := render(indexTpl, "layout", "index.html", data); err != nil {
		return err
	}

	// 渲染feed.xml ( use text/template )
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
		SiteURL: "https://sophieink.github.io/inkbottle-blog/index.html",
		Posts: listed, // 复用index
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
