package main

import (
	"fisher/lib"
	"github.com/davecgh/go-spew/spew"
	"io"
	"os"
)

const savePath = "/opt/download/"

func dealWithHtml(keyword, url, index, content string) {
	var path = savePath + keyword + "/" + index

	if _, err := os.Stat(savePath + keyword); os.IsNotExist(err) {
		// 必须分成两步：先创建文件夹、再修改权限
		os.Mkdir(savePath+keyword, 0777)
		// os.Chmod(savePath+keyword, 0777)
	}

	//创建文件
	f, _ := os.Create(path + ".html")
	//写入文件(字符串)
	if _, err := io.WriteString(f, content); err != nil {
		lib.Log("保存关键词:", keyword, "第", index, "个html出错，错误信息为:", err)
	}

	if videoUrl := lib.GetVideoUrl(content); videoUrl != "" {
		if err := lib.DownloadFile(path+".mp4", videoUrl); err != nil {
			lib.Log("下载url:", videoUrl, "出错，错误信息为:", err)
		}
	}
}

func start(keyword string) {
	// 抓取第一页到第三页，总共 44*3条记录
	// for i := 1; i < 4; i++ {
	for i := 1; i < 2; i++ {
		lib.SearchTB(keyword, i, func(url, index, html string) { dealWithHtml(keyword, url, index, html) })
	}
}

func main() {
	start("猫粮")
}

func test() {
	// 格式1
	detailContent, _ := lib.HttpGet("https://detail.tmall.com/item.htm?id=534714077614&ali_trackid=2:mm_26632614_0_0:1556445756_275_1897617264&spm=a21bo.7925826.192013.1.62504c0d2Hbyvd")
	// 格式2
	// detailContent, _ := lib.HttpGet("https://item.taobao.com/item.htm?spm=a217f.8051907.312171.37.7c4c3308ygq3wn&id=576633151569&qq-pf-to=pcqq.c2c")
	spew.WDump(lib.GetVideoUrl(detailContent))
}
