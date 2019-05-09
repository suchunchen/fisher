package main

import (
	// "fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

func httpGet(requestUrl string) (response string, header map[string][]string) {
	if strings.HasPrefix(requestUrl, "//") {
		requestUrl = "https:" + requestUrl
	}

	spew.WDump(requestUrl)
	resp, err := http.Get(requestUrl)

	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if loc, ok := resp.Header["Location"]; ok {
		return httpGet(loc[0])
	}

	if body, err := ioutil.ReadAll(resp.Body); err != nil {
		panic(err)
	} else {
		return string(body), resp.Header
	}
}

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

	// 返回的请求体 resp.Header 中以下2个字段有何作用，待研究
	// Eagleeye-Traceid: 0b838cef15573910940247287e4d06
	// Set-Cookie: JSESSIONID=7CC8CEAFCE14ACD185488BCF4A5F650B; Path=/; HttpOnly
	response, _ := httpGet(requestUrl)

	matches := regexp.MustCompile("g_page_config = ({.*});").FindAllStringSubmatch(response, -1)
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
	itemlist.ForEach(func(index, item gjson.Result) bool {
		detailUrlList = append(detailUrlList, item.Get("detail_url").String())

		detailContent, _ := httpGet(item.Get("detail_url").String())
		getDetailVideo(detailContent)

		return true
	})

	spew.WDump(detailUrlList)
}

// 获取详情页的视频链接
func getDetailVideo(content string) string {
	// 格式1：
	// //cloud.video.taobao.com/play/u/3351172141/p/1/e/1/t/8/218965885328.swf
	matches := regexp.MustCompile("//cloud.video.taobao.com/play/u/(\\d+)/p/1/e/1/t/8/(\\d+).swf").FindAllStringSubmatch(content, -1)

	var videoId, videoOwnerId string

	if len(matches) == 1 && len(matches[0]) == 3 {
		videoOwnerId = matches[0][1]
		videoId = matches[0][2]
	} else {
		// 格式2
		// Hub.config.set('video', {"videoDuaration":"30","picUrl":"//img.alicdn.com/imgextra/i2/334724657/O1CN011kGwwXENhmxmvD7_!!334724657.jpg_310x310.jpg","videoId":"50254186813","autoplay":"1","videoOwnerId":"334724657","videoStatus":"0"})
		videoOwnerId = getInfoFromJson("videoOwnerId", content)
		videoId = getInfoFromJson("videoId", content)

		if videoOwnerId == "" || videoId == "" {
			return ""
		}
	}

	return "https://cloud.video.taobao.com/play/u/" + videoOwnerId + "/p/1/e/6/t/1/" + videoId + ".mp4"
}

func getInfoFromJson(key, content string) string {
	matches := regexp.MustCompile("\""+key+"\":\"(\\d+)\"").FindAllStringSubmatch(content, -1)

	if len(matches) == 1 && len(matches[0]) == 2 {
		return matches[0][1]
	}
	return ""
}

func main() {
	// search("女装")

	// 格式1
	detailContent, _ := httpGet("https://detail.tmall.com/item.htm?id=534714077614&ali_trackid=2:mm_26632614_0_0:1556445756_275_1897617264&spm=a21bo.7925826.192013.1.62504c0d2Hbyvd")

	// 格式2
	// detailContent, _ := httpGet("https://item.taobao.com/item.htm?spm=a217f.8051907.312171.37.7c4c3308ygq3wn&id=576633151569&qq-pf-to=pcqq.c2c")
	spew.WDump(getDetailVideo(detailContent))
}
