package internal

import (
	"net/http"
	"net/url"
	"strings"
	"sync"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

type Crawler struct {
	Mu          sync.Mutex
	ValidUrls   map[string]map[string]struct{}
	InvalidUrls map[string]struct{}
	Visited     map[string]struct{}
}

const webScraperUrl string = "scrape-me.dreamsofcode.io"

func (c *Crawler) BuildUrl(path string) url.URL {
	return url.URL{Scheme: "https", Host: webScraperUrl, Path: path}
}

func (c *Crawler) TraverseLinks(path string, wg *sync.WaitGroup, sem chan struct{}) {
	defer wg.Done()

	sem <- struct{}{}
	defer func() { <-sem }()

	c.Mu.Lock()
	if _, ok := c.Visited[path]; ok {
		c.Mu.Unlock()
		return
	}
	c.Visited[path] = struct{}{}
	c.Mu.Unlock()

	baseUrl := c.BuildUrl(path)
	resp, err := http.Get(baseUrl.String())

	if err != nil || resp.StatusCode != http.StatusOK {
		c.Mu.Lock()
		c.InvalidUrls[path] = struct{}{}
		c.Mu.Unlock()
		if resp != nil {
			resp.Body.Close()
		}
		return
	}
	defer resp.Body.Close()

	root, err := html.Parse(resp.Body)
	if err != nil {
		return
	}

	body := c.TraverseBody(root)
	if body == nil {
		return
	}

	nodes := []html.Node{}
	nodes = c.CollectNodes(body, atom.A, nodes)
	urls := c.FilterLinksByHostname(nodes)

	c.Mu.Lock()
	c.ValidUrls[path] = urls
	c.Mu.Unlock()

	for url := range urls {
		wg.Add(1)
		go c.TraverseLinks(url, wg, sem)
	}
}

func (cr *Crawler) CollectNodes(n *html.Node, a atom.Atom, nodes []html.Node) []html.Node {
	if n.Type == html.ElementNode && n.DataAtom == a {
		nodes = append(nodes, *n)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		nodes = cr.CollectNodes(c, a, nodes)
	}
	return nodes
}

func (cr *Crawler) TraverseBody(n *html.Node) *html.Node {
	if n.Type == html.ElementNode && n.DataAtom == atom.Body {
		return n
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if body := cr.TraverseBody(c); body != nil {
			return body
		}
	}
	return nil
}

func (cr *Crawler) FilterLinksByHostname(nodes []html.Node) map[string]struct{} {
	result := map[string]struct{}{}

	for _, n := range nodes {
		for _, a := range n.Attr {
			key, val := a.Key, a.Val
			if key == "href" {
				u, _ := url.Parse(val)
				_, ok := result[u.Path]
				if strings.TrimSpace(u.Hostname()) == "" && !ok {
					result[u.Path] = struct{}{}
				}
			}
		}
	}
	return result
}

func (c *Crawler) RenderNode(n *html.Node) string {
	var buf strings.Builder
	err := html.Render(&buf, n)
	if err != nil {
		return ""
	}
	return buf.String()
}
