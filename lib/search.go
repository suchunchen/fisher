package lib

import (
	"fmt"
	"github.com/json-iterator/go"
	"github.com/tidwall/gjson"
	"net/url"
	"regexp"
	"time"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func ToJsonString(v interface{}) string {
	b, err := json.Marshal(v)
	if err == nil {
		return string(b[:])
	}
	return ""
}

func SearchTB(keyword string, page int, dealWithHtml func(detailUrl, index, content string)) {
	// https://s.taobao.com/search?q=%E5%A5%B3%E8%A3%85&imgfile=&js=1&stats_click=search_radio_all%3A1&initiative_id=staobaoz_20190509&ie=utf8
	u := url.Values{}
	u.Set("q", keyword)
	u.Set("initiative_id", "staobaoz_"+time.Now().Format("20060102"))
	u.Set("imgfile", "")
	u.Set("js", "1")
	u.Set("stats_click", "search_radio_all:1")
	u.Set("ie", "utf8")

	if page > 1 {
		if page == 3 {
			// 别问我为什么，就是这样的
			u.Set("ntoffset", "6")
		} else {
			u.Set("ntoffset", fmt.Sprintf("%d", 6-3*(page-1)))
		}

		u.Set("bcoffset", fmt.Sprintf("%d", 6-3*(page-1)))
		u.Set("p4ppushleft", "1,48")
		u.Set("s", fmt.Sprintf("%d", (page-1)*44))
	}

	// 返回的请求体 resp.Header 中以下2个字段有何作用，待研究
	// Eagleeye-Traceid: 0b838cef15573910940247287e4d06
	// Set-Cookie: JSESSIONID=7CC8CEAFCE14ACD185488BCF4A5F650B; Path=/; HttpOnly
	response, _ := HttpGet("https://s.taobao.com/search?" + u.Encode())

	matches := regexp.MustCompile("g_page_config = ({.*});").FindAllStringSubmatch(response, -1)
	if len(matches) != 1 || len(matches[0]) != 2 {
		Log("没找到 g_page_config 相关配置")
		return
	}

	itemlist := gjson.Get(matches[0][1], "mods.itemlist.data.auctions")

	if itemlist.String() == "" {
		Log("没有商品信息")
		return
	}

	var (
		detailUrlList []string
		index         int = 1
	)

	// 获取每个搜索页的detailUrlList
	itemlist.ForEach(func(i, item gjson.Result) bool {
		detailUrl := item.Get("detail_url").String()

		detailUrlList = append(detailUrlList, detailUrl)
		detailContent, _ := HttpGet(detailUrl)

		dealWithHtml(detailUrl, fmt.Sprintf("%d", index), detailContent)

		index++

		if index > 5 {
			return false
		}

		return true
	})

	Log("关键词：", keyword, "第", page, "页的连接为: \n", ToJsonString(detailUrlList))
}

// 获取详情页的视频链接
func GetVideoUrl(content string) string {
	// 格式1：
	// cloud.video.taobao.com/play/u/3351172141/p/1/e/1/t/8/218965885328.swf
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
