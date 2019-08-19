package xspider

import (
	"fmt"
	"path/filepath"

	"github.com/PuerkitoBio/goquery"
)

//Video Video
type Video struct {
	Title    string  `json:"title"`
	CoverURL string  `json:"coverurl"`
	PlayURL  string  `json:"playurl"`
	Height   int     `json:"height"`
	Width    int     `json:"width"`
	Fps      float64 `json:"fps"`
	Duration float64 `json:"duration"`
	Bitrate  float64 `json:"bitrate"`
	Size     int     `json:"size"`
	Format   float64 `json:"format"`
}

type Zhihu struct{}

func (n Zhihu) Video(weburl string) (video []Video) {
	PlayURL := "https://lens.zhihu.com/api/v4/videos/%s"

	type _zhihuVideo struct {
		Playlist struct {
			LD struct {
				Format   string  `json:"format"`
				PlayURL  string  `json:"play_url"`
				Height   int     `json:"height"`
				Width    int     `json:"width"`
				Fps      float64 `json:"fps"`
				Duration float64 `json:"duration"`
				Bitrate  float64 `json:"bitrate"`
				Size     int     `json:"size"`
			} `json:"LD"`
		} `json:"playlist"`
		CoverURL string `json:"cover_url"`
		Title    string `json:"title"`
	}

	spider := Spider{Type: "GET", Referer: "http://www.zhihu.com/", URL: weburl}
	res, err := spider.Fetch()
	if err != nil {
		fmt.Println(err)
	}
	defer res.Close()
	doc, err := goquery.NewDocumentFromResponse(res.RawResponse)

	doc.Find("a[class='video-box'][href]").Each(func(i int, sel *goquery.Selection) {
		// href="https://link.zhihu.com/?target=https%3A//www.zhihu.com/video/1118600026741563392"
		vid, _ := sel.Attr("href")
		spider.URL = fmt.Sprintf(PlayURL, filepath.Base(vid))

		var v _zhihuVideo
		err = spider.JSON(&v)
		if err != nil {
			fmt.Println(err)
		}

		video = append(video, Video{
			Title:    v.Title,
			CoverURL: v.CoverURL,
			PlayURL:  v.Playlist.LD.PlayURL,
			Height:   v.Playlist.LD.Height,
			Width:    v.Playlist.LD.Width,
			Fps:      v.Playlist.LD.Fps,
			Duration: v.Playlist.LD.Duration,
			Bitrate:  v.Playlist.LD.Bitrate,
			Size:     v.Playlist.LD.Size})
	})
	return video
}

type WeChat struct{}

func (w WeChat) Video(weburl string) (video []Video) {

	// https://mp.weixin.qq.com/s/9zMYS4YXdse3nCqYC_63lQ
	// http://v.qq.com/boke/page/d/0/v/y08118njwd0.html
	// https://v.qq.com/iframe/preview.html?width=500&amp;height=375&amp;auto=0&amp;vid=y08118njwd0

	PlayURL := "https://mp.weixin.qq.com/mp/videoplayer?action=get_mp_video_play_url&preview=0&__biz=&mid=&idx=&vid=%s&uin=&key=&pass_ticket=&wxtoken=&appmsg_token=&x5=0&f=json"

	type _wxVideo struct {
		Title   string `json:"title"`
		URLInfo []struct {
			URL        string `json:"url"`
			FormatID   int    `json:"format_id"`
			DurationMs int    `json:"duration_ms"`
			Filesize   int    `json:"filesize"`
			Width      int    `json:"width"`
			Height     int    `json:"height"`
		} `json:"url_info"`
	}

	spider := Spider{Type: "GET", Referer: "https://mp.weixin.qq.com", URL: weburl}

	res, err := spider.Fetch()
	if err != nil {
		fmt.Println(err)
	}
	defer res.Close()

	// raw := res.Bytes()
	// ioutil.WriteFile("demo.html", raw, 0666)

	// br := bytes.NewReader(raw)
	doc, err := goquery.NewDocumentFromResponse(res.RawResponse)

	// doc, err := goquery.NewDocumentFromResponse(res.RawResponse)
	if err != nil {
		fmt.Println(err)
	}

	doc.Find("iframe[class='video_iframe rich_pages'][data-src]").Each(func(i int, sel *goquery.Selection) {
		src, _ := sel.Attr("data-src")
		// https://v.qq.com/iframe/preview.html?width=500&height=375&auto=0&vid=i089578m3l0
		fmt.Println(src)

	})

	doc.Find("iframe[class='video_iframe rich_pages'][data-mpvid]").Each(func(i int, sel *goquery.Selection) {
		vid, _ := sel.Attr("data-mpvid")

		spider.URL = fmt.Sprintf(PlayURL, vid)
		var v _wxVideo
		err = spider.JSON(&v)
		if err != nil {
			fmt.Println(err)
		}
		for _, vv := range v.URLInfo {
			video = append(video, Video{
				Title:    v.Title,
				PlayURL:  vv.URL,
				Height:   vv.Height,
				Width:    vv.Width,
				Duration: float64(vv.DurationMs),
				Size:     vv.Filesize})
		}

	})
	doc.Find("div#video_container txpdiv[class=txp_video_container] video[src]").Each(func(i int, sel *goquery.Selection) {
		src, _ := sel.Attr("src")
		video = append(video, Video{
			Title:   "",
			PlayURL: src})
	})

	return video
}
