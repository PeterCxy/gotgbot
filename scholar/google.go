// Google search functions
package scholar

import (
	"fmt"
	qs "net/url"

	"github.com/ddliu/go-httpclient"
	"github.com/PuerkitoBio/goquery"

	"github.com/paulrosania/go-charset/charset"
	_ "github.com/paulrosania/go-charset/data"
)

type GoogleResult struct {
	title string
	url string
	summary string
}

const (
	linkSel = "h3.r a"
	descSel = "div.s"
	itemSel = "li.g"
	nextSel = "td.b a span"
)

func Google(query string, start int, num int, ipv6 bool) (ret []GoogleResult, hasNext bool) {
	url := "https://google.com/search"

	if ipv6 {
		url = "https://ipv6.google.com/search"
	}

	res, err := httpclient.
		WithHeader("Accept-Charset", "utf-8").
		Get(url, map[string]string {
			"hl": "en",
			"q": query,
			"start": fmt.Sprintf("%d", start),
			"num": fmt.Sprintf("%d", num),
	})

	r, _ := charset.NewReader("iso-8859-1", res.Body)

	//defer r.Close()
	defer res.Body.Close()

	if (err != nil) || (res.StatusCode != 200) {
		return
	}

	doc, errDoc := goquery.NewDocumentFromReader(r)

	if (errDoc != nil) {
		return
	}

	doc.Find(itemSel).Each(func(i int, sel *goquery.Selection) {
		link := sel.Find(linkSel)
		desc := sel.Find(descSel)
		href, ok := link.Attr("href")

		if ok {
			u, err := qs.Parse(href)

			if err != nil {
				return
			}

			desc.Find("div").Remove()

			ret = append(ret, GoogleResult {
				title: link.First().Text(),
				url: u.Query().Get("q"),
				summary: desc.Text(),
			})
		}
	})

	if doc.Find(nextSel).Last().Text() == "Next" {
		hasNext = true
	}

	return
}
