package xspider

import (
	"bytes"
	"fmt"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

//DoubanMovie DoubanMovie
type DoubanMovie struct {
	ID            string    `json:"id"`
	Title         string    `json:"title"`
	Cover         string    `json:"cover"`
	Director      string    `json:"director"`
	Author        string    `json:"author"`
	Actor         string    `json:"actor"`
	Genre         string    `json:"genre"`
	Area          string    `json:"area"`
	Language      string    `json:"language"`
	DatePublished string    `json:"datepublished"`
	Duration      string    `json:"duration"`
	Alias         string    `json:"alias"`
	IMDb          string    `json:"imdb"`
	Rate          float64   `json:"rate"`
	RatingCount   int       `json:"ratingcount"`
	RatingWeight  []float64 `json:"ratingweight"`
	Tag           []string  `json:"tag"`
	Recommend     []string  `json:"recommend"`
}

//Douban Douban
type Douban struct{}

//MovieID MovieID
func (n Douban) MovieID(raw []byte) DoubanMovie {
	r := bytes.NewReader(raw)
	doc, _ := goquery.NewDocumentFromReader(r)
	var douban DoubanMovie
	pageid := doc.Find("#mainpic a.nbgnbg").AttrOr("href", "")
	movieid := regexp.MustCompile(`https://movie.douban.com/subject/(.*)/.*`).FindSubmatch([]byte(pageid))
	if len(movieid) > 0 {
		douban.ID = string(movieid[1])
	}

	// douban.ID = filepath.Base(res.Request.URL.String())
	douban.Title, _ = doc.Find("#content h1 span").Eq(0).Html()
	douban.Cover = doc.Find("#mainpic a img[src]").AttrOr("src", "")
	doc.Find("div.tags-body a[href]").Each(func(i int, sel *goquery.Selection) {
		douban.Tag = append(douban.Tag, filepath.Base(sel.AttrOr("href", "")))
	})
	doc.Find("div.recommendations-bd dd a[href]").Each(func(i int, sel *goquery.Selection) {
		douban.Recommend = append(douban.Recommend, filepath.Base(
			strings.ReplaceAll(sel.AttrOr("href", ""), "/?from=subject-page", "")))
	})
	rating := doc.Find("#interest_sectl div[class='rating_wrap clearbox'] div[class='rating_self clearfix'] strong[class='ll rating_num']").Text()
	douban.Rate, _ = strconv.ParseFloat(rating, 64)

	ratingCount := doc.Find("#interest_sectl div[class='rating_wrap clearbox'] div[class='rating_self clearfix'] div.rating_right div.rating_sum a span").Text()
	douban.RatingCount, _ = strconv.Atoi(ratingCount)

	doc.Find("#interest_sectl div.rating_wrap.clearbox div.ratings-on-weight div.item span.rating_per").Each(
		func(i int, sel *goquery.Selection) {
			weight := strings.ReplaceAll(sel.Text(), "%", "")
			weightfloat, _ := strconv.ParseFloat(weight, 64)
			douban.RatingWeight = append(douban.RatingWeight, weightfloat)
		})

	info := []byte(doc.Find("#info").Text())
	info1 := regexp.MustCompile(`导演: (.*)`).FindSubmatch(info)
	info2 := regexp.MustCompile(`编剧: (.*)`).FindSubmatch(info)
	info3 := regexp.MustCompile(`主演: (.*)`).FindSubmatch(info)
	info4 := regexp.MustCompile(`类型: (.*)`).FindSubmatch(info)
	info5 := regexp.MustCompile(`制片国家/地区: (.*)`).FindSubmatch(info)
	info6 := regexp.MustCompile(`语言: (.*)`).FindSubmatch(info)
	info7 := regexp.MustCompile(`上映日期: (.*)`).FindSubmatch(info)
	info8 := regexp.MustCompile(`片长: (.*)`).FindSubmatch(info)
	info9 := regexp.MustCompile(`又名: (.*)`).FindSubmatch(info)
	info10 := regexp.MustCompile(`IMDb链接: (.*)`).FindSubmatch(info)

	if len(info1) > 1 {
		douban.Director = string(info1[1])
	}
	if len(info2) > 1 {
		douban.Author = string(info2[1])
	}
	if len(info3) > 1 {
		douban.Actor = string(info3[1])
	}
	if len(info4) > 1 {
		douban.Genre = string(info4[1])
	}
	if len(info5) > 1 {
		douban.Area = string(info5[1])
	}
	if len(info6) > 1 {
		douban.Language = string(info6[1])
	}

	if len(info7) > 1 {
		douban.DatePublished = string(info7[1])
	}

	if len(info8) > 1 {
		douban.Duration = string(info8[1])
	}

	if len(info9) > 1 {
		douban.Alias = string(info9[1])
	}
	if len(info10) > 1 {
		douban.IMDb = string(info10[1])
	}
	return douban
}

//GetMovie250 解析URL获取Movie
func _Movie250() {
	type Movie struct {
		URL   string //URL
		Title string
		Cover string
		Info  string
		Rate  string
		Quote string
	}
	spider := Spider{URL: "https://movie.douban.com/top250?start=125&filter="}
	raw, _ := spider.FetchBytes()
	r := bytes.NewReader(raw)
	doc, _ := goquery.NewDocumentFromReader(r)

	movies := make([]Movie, 0)
	doc.Find(`div.article>ol.grid_view>li>div.item`).Each(func(i int, s *goquery.Selection) {
		r := Movie{}
		r.Title = s.Find(`div.info div.hd a span`).Eq(0).Text()
		r.Cover, _ = s.Find(`div[class="pic"] a img`).Eq(0).Attr("src")
		r.Info = s.Find(`div[class="info"] div[class="bd"] p[1]`).Eq(0).Text()
		r.Rate = s.Find(`div[class="info"] div[class="bd"] div[class="star"] span[class="rating_num"]`).Eq(0).Text()
		r.Quote = s.Find(`div[class="info"] div[class="bd"] p[class="quote"] span[class="inq"]')`).Eq(0).Text()
		movies = append(movies, r)
	})
	fmt.Println(movies)
}
