package lib

import (
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func HttpGet(requestUrl string) (response string, header map[string][]string) {
	if strings.HasPrefix(requestUrl, "//") {
		requestUrl = "https:" + requestUrl
	}
	Log("请求url:", requestUrl)
	resp, err := http.Get(requestUrl)

	if resp != nil {
		defer resp.Body.Close()

		if loc, ok := resp.Header["Location"]; ok {
			Log("请求url:", requestUrl, "\n返回302，跳转到:\n", loc[0])
			return HttpGet(loc[0])
		}
	}

	if err != nil {
		panic(err)
	}

	if body, err := ioutil.ReadAll(resp.Body); err != nil {
		panic(err)
	} else {
		return string(body), resp.Header
	}
}

func DownloadFile(filepath string, url string) (err error) {
	Log("下载文件url:", url)
	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}
