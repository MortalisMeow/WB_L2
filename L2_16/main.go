package main

import (
	"fmt"
	"golang.org/x/net/html"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

// Отслеживание посещенных URL
var visited = make(map[string]bool)

func main() {
	startURL := "https://example.com"
	maxDepth := 2

	// Создаем корневую папку
	rootDir := "mirror"
	os.MkdirAll(rootDir, 0755)

	// Запускаем рекурсивное скачивание
	downloadRecursive(startURL, rootDir, 0, maxDepth)

	fmt.Println("Done!")
}

func downloadRecursive(rawURL string, rootDir string, depth int, maxDepth int) {
	// Проверяем глубину
	if depth > maxDepth {
		return
	}

	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return
	}
	parsedURL.Fragment = ""
	cleanURL := parsedURL.String()

	// Проверяем, не посещали ли уже
	if visited[cleanURL] {
		return
	}
	visited[cleanURL] = true

	fmt.Printf("[Depth %d] Downloading: %s\n", depth, cleanURL)

	// Определяем локальный путь для сохранения
	localPath := urlToFilePath(parsedURL, rootDir)

	// Создаем папки для файла
	os.MkdirAll(filepath.Dir(localPath), 0755)

	// Скачиваем файл
	resp, err := http.Get(cleanURL)
	if err != nil {
		fmt.Printf("Error downloading %s: %v\n", cleanURL, err)
		return
	}
	defer resp.Body.Close()

	// Читаем содержимое
	body, _ := io.ReadAll(resp.Body)

	// Проверяем тип
	contentType := resp.Header.Get("Content-Type")

	if strings.Contains(contentType, "text/html") {
		links := extractLinks(string(body), parsedURL)

		modifiedHTML := rewriteLinks(string(body), parsedURL, rootDir)

		os.WriteFile(localPath, []byte(modifiedHTML), 0644)

		for _, link := range links {
			linkURL, err := url.Parse(link)
			if err != nil {
				continue
			}

			if isSameDomain(linkURL, parsedURL) {
				downloadRecursive(link, rootDir, depth+1, maxDepth)
			}
		}
	} else {
		os.WriteFile(localPath, body, 0644)
	}
}

// Преобразует URL в путь
func urlToFilePath(u *url.URL, rootDir string) string {
	path := u.Host + u.Path

	if strings.HasSuffix(path, "/") {
		path += "index.html"
	}

	if filepath.Ext(path) == "" {
		path += ".html"
	}

	return filepath.Join(rootDir, path)
}

// Берем все ссылки из html
func extractLinks(htmlContent string, baseURL *url.URL) []string {
	var links []string

	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		return links
	}

	var extractFromNode func(*html.Node)
	extractFromNode = func(n *html.Node) {
		if n.Type == html.ElementNode {
			var href string

			switch n.Data {
			case "a", "link":
				href = getAttr(n, "href")
			case "img", "script", "source":
				href = getAttr(n, "src")
			}

			if href != "" {
				absoluteURL := resolveURL(baseURL, href)
				if absoluteURL != "" {
					links = append(links, absoluteURL)
				}
			}
		}

		// Рекурсивно обходим дочерние узлы
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			extractFromNode(c)
		}
	}

	extractFromNode(doc)
	return links
}

// Заменяет ссылки в html на локальные пути
func rewriteLinks(htmlContent string, baseURL *url.URL, rootDir string) string {
	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		return htmlContent
	}

	var rewriteNode func(*html.Node)
	rewriteNode = func(n *html.Node) {
		if n.Type == html.ElementNode {
			var attrName string

			switch n.Data {
			case "a", "link":
				attrName = "href"
			case "img", "script", "source":
				attrName = "src"
			}

			if attrName != "" {
				for i, attr := range n.Attr {
					if attr.Key == attrName {
						absoluteURL := resolveURL(baseURL, attr.Val)
						if absoluteURL != "" {
							parsedURL, _ := url.Parse(absoluteURL)
							localPath := urlToFilePath(parsedURL, rootDir)
							n.Attr[i].Val = localPath
						}
					}
				}
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			rewriteNode(c)
		}
	}

	rewriteNode(doc)

	// Сериализуем обратно
	var buf strings.Builder
	html.Render(&buf, doc)
	return buf.String()
}

func getAttr(n *html.Node, attrName string) string {
	for _, attr := range n.Attr {
		if attr.Key == attrName {
			return attr.Val
		}
	}
	return ""
}

func resolveURL(base *url.URL, href string) string {
	if strings.HasPrefix(href, "javascript:") ||
		strings.HasPrefix(href, "mailto:") ||
		strings.HasPrefix(href, "#") {
		return ""
	}

	parsed, err := url.Parse(href)
	if err != nil {
		return ""
	}

	resolved := base.ResolveReference(parsed)
	return resolved.String()
}

func isSameDomain(u1 *url.URL, u2 *url.URL) bool {
	return u1.Host == u2.Host
}
