package xspider

import (
	"math/rand"
	"net/url"
	"time"
)

//ProxyPool ProxyPool
type ProxyPool struct {
	URLs        []*url.URL
	MaxDuration int64 //更新间隔
	Lasttime    int64 //最后更新时间
}

//Refresh Refresh
func (p *ProxyPool) Refresh() {
	type ProxHTTP struct {
		Error   string `json:"error"`
		Message []struct {
			ID         int    `json:"id"`
			IP         string `json:"ip"`
			Port       string `json:"port"`
			SchemeType int    `json:"scheme_type"`
			Content    string `json:"content"`
		} `json:"message"`
	}

	now := time.Now().Unix()
	if now-p.Lasttime > p.MaxDuration {
		weburl := "http://localhost:9999/sql?query=SELECT%20*%20FROM%20proxy%20order%20by%20score%20%20desc%20limit%20100"
		// weburl := "http://127.0.0.1:5010/get_all/"
		p.URLs = make([]*url.URL, 0)

		spider := Spider{URL: weburl}
		var proxhttp ProxHTTP
		// var proxhttp []string
		_ = spider.JSON(&proxhttp)
		// for _, v := range proxhttp {
		// 	p.URLs = append(p.URLs, &url.URL{Host: v})
		// }
		for _, v := range proxhttp.Message {
			scheme := ""
			if v.SchemeType == 0 {
				scheme = "http"
			} else {
				scheme = "https"
			}
			p.URLs = append(p.URLs, &url.URL{Scheme: scheme, Host: v.Content})
		}
		p.Lasttime = now
	}
}

//GetOne  刷新并随机获取一个
func (p ProxyPool) GetOne() *url.URL {
	p.Refresh()
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	return p.URLs[random.Intn(len(p.URLs))]
}
