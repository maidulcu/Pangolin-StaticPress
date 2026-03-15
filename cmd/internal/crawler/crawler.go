package crawler

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"

	"github.com/pangolin-cms/staticpress/cmd/internal/config"

	"github.com/PuerkitoBio/goquery"
)

const (
	UserAgent = "Pangolin/1.0 (Static Site Generator)"
)

type Page struct {
	URL     string
	HTML    string
	Headers http.Header
}

type Asset struct {
	URL  string
	Path string
	Type string
}

var (
	assetsDownloaded = make(map[string]bool)
	assetsLock       sync.Mutex
)

func FetchPage(url string, distDir string) (*Page, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", UserAgent)

	cfg, _ := config.LoadConfig()
	if cfg != nil && cfg.APIKey != "" {
		req.Header.Set("Authorization", "Bearer "+cfg.APIKey)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("failed to fetch %s: status %d", url, resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	cfg, _ = config.LoadConfig()
	baseURL := cfg.SiteURL

	downloadAssets(doc, baseURL, distDir)
	rewriteLinks(doc, baseURL)

	html, err := doc.Html()
	if err != nil {
		return nil, err
	}

	return &Page{
		URL:     url,
		HTML:    html,
		Headers: resp.Header,
	}, nil
}

func downloadAssets(doc *goquery.Document, baseURL, distDir string) {
	var wg sync.WaitGroup
	assetChan := make(chan Asset, 100)

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for asset := range assetChan {
				downloadAsset(asset.URL, distDir+asset.Path, asset.Type)
			}
		}()
	}

	sel := doc.Find("img[src]")
	sel.Each(func(i int, s *goquery.Selection) {
		src, exists := s.Attr("src")
		if !exists {
			return
		}

		if strings.HasPrefix(src, "http") {
			assetPath := "/images/" + filepath.Base(src)
			assetChan <- Asset{URL: src, Path: assetPath, Type: "image"}
			s.SetAttr("src", assetPath)
		} else if strings.HasPrefix(src, baseURL) {
			assetPath := "/images/" + filepath.Base(strings.TrimPrefix(src, baseURL))
			assetChan <- Asset{URL: baseURL + src, Path: assetPath, Type: "image"}
			s.SetAttr("src", assetPath)
		}
	})

	sel = doc.Find("link[href]")
	sel.Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if !exists {
			return
		}

		if !strings.HasSuffix(href, ".css") {
			return
		}

		if strings.HasPrefix(href, "http") {
			assetPath := "/assets/" + filepath.Base(href)
			assetChan <- Asset{URL: href, Path: assetPath, Type: "css"}
			s.SetAttr("href", assetPath)
		} else if strings.HasPrefix(href, baseURL) {
			cleanPath := strings.TrimPrefix(href, baseURL)
			assetPath := "/assets/" + filepath.Base(cleanPath)
			assetChan <- Asset{URL: baseURL + cleanPath, Path: assetPath, Type: "css"}
			s.SetAttr("href", assetPath)
		}
	})

	sel = doc.Find("script[src]")
	sel.Each(func(i int, s *goquery.Selection) {
		src, exists := s.Attr("src")
		if !exists {
			return
		}

		if strings.HasPrefix(src, "http") {
			assetPath := "/assets/" + filepath.Base(src)
			assetChan <- Asset{URL: src, Path: assetPath, Type: "js"}
			s.SetAttr("src", assetPath)
		} else if strings.HasPrefix(src, baseURL) {
			cleanPath := strings.TrimPrefix(src, baseURL)
			assetPath := "/assets/" + filepath.Base(cleanPath)
			assetChan <- Asset{URL: baseURL + cleanPath, Path: assetPath, Type: "js"}
			s.SetAttr("src", assetPath)
		}
	})

	close(assetChan)
	wg.Wait()
}

func downloadAsset(url, destPath, assetType string) {
	assetsLock.Lock()
	if assetsDownloaded[url] {
		assetsLock.Unlock()
		return
	}
	assetsDownloaded[url] = true
	assetsLock.Unlock()

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}
	req.Header.Set("User-Agent", UserAgent)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return
	}

	dir := path.Dir(destPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return
	}

	file, err := os.Create(destPath)
	if err != nil {
		return
	}
	defer file.Close()

	io.Copy(file, resp.Body)
}

func rewriteLinks(doc *goquery.Document, baseURL string) {
	sel := doc.Find("a[href]")
	sel.Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if !exists {
			return
		}

		if strings.HasPrefix(href, baseURL) {
			newHref := strings.TrimPrefix(href, baseURL)
			if newHref == "" || newHref == "/" {
				newHref = "/"
			}
			s.SetAttr("href", newHref)
		}
	})
}
