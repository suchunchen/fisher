package main

import (
	// "fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"time"
)

func search(keyword string) {
	// https://s.taobao.com/search?q=%E5%A5%B3%E8%A3%85&imgfile=&js=1&stats_click=search_radio_all%3A1&initiative_id=staobaoz_20190509&ie=utf8
	u := url.Values{}
	u.Set("q", keyword)
	u.Set("initiative_id", "staobaoz_"+time.Now().Format("20060102"))
	u.Set("imgfile", "")
	u.Set("js", "1")
	u.Set("stats_click", "search_radio_all:1")
	u.Set("ie", "utf8")

	requestUrl := "https://s.taobao.com/search?" + u.Encode()

	spew.WDump(requestUrl)

	resp, err := http.Get(requestUrl)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// 返回的请求体 resp.Header 中以下2个字段有何作用，待研究
	// Eagleeye-Traceid: 0b838cef15573910940247287e4d06
	// Set-Cookie: JSESSIONID=7CC8CEAFCE14ACD185488BCF4A5F650B; Path=/; HttpOnly
	if body, err := ioutil.ReadAll(resp.Body); err != nil {
		panic(err)
	} else {
		// spew.WDump(string(body))
		re := regexp.MustCompile("g_page_config = ({.*});")

		matches := re.FindAllStringSubmatch(string(body), -1)

		if len(matches) != 1 || len(matches[0]) != 2 {
			spew.WDump("没找到")
			return
		}

		itemlist := gjson.Get(matches[0][1], "mods.itemlist.data.auctions")

		if itemlist.String() == "" {
			spew.WDump("没有商品信息")
			return
		}

		var detailUrlList []string

		// 获取每个搜索页的detailUrlList
		itemlist.ForEach(func(key, item gjson.Result) bool {
			detailUrlList = append(detailUrlList, item.Get("detail_url").String())
			return true
		})

		spew.WDump(detailUrlList)
	}

}

func main() {
	search("女装")
}
