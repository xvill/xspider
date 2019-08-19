package xspider

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"strings"

	jsoniter "github.com/json-iterator/go"
	"github.com/json-iterator/go/extra"
)

// baseurl   = "http://search.kuwo.cn/r.s?all=%s&ft=music&itemset=web_2018&client=kt&pn=%s&rn=10&rformat=json&encoding=utf8"
// searchurl = "http://m.kuwo.cn/newh5/singles/songinfoandlrc?musicId=%s"
// palyurl   = "http://antiserver.kuwo.cn/anti.s?type=convert_url&rid=%s&format=aac|mp3&response=url"
// // commenturl  = "http://comment.kuwo.cn/com.s?type=get_comment&uid=0&prod=newWeb&digest=15&sid=%s&page=1&rows=10&f=web"
// // songlisturl = "http://yinyue.kuwo.cn/yy/cinfo_%s.htm"

//Kuwo Kuwo
type Kuwo struct{}

//Search Search
func (k Kuwo) Search(keyword string) []Song {
	type Abslist struct {
		Songs []struct {
			Album    string `json:"ALBUM"`
			Albumid  string `json:"ALBUMID"`
			Artist   string `json:"ARTIST"`
			Artistid string `json:"ARTISTID"`
			Musicrid string `json:"MUSICRID"`
			Songname string `json:"SONGNAME"`
		} `json:"abslist"`
	}

	var abslist Abslist
	page_ := 0
	baseurl := "http://search.kuwo.cn/r.s?all=%s&ft=music&itemset=web_2018&client=kt&pn=%d&rn=10&rformat=json&encoding=utf8"

	spider := Spider{Type: "GET", Referer: "http://player.kuwo.cn/webmusic/play",
		URL: fmt.Sprintf(baseurl, keyword, page_)}

	raw, err := spider.FetchBytes()
	if err != nil {
		fmt.Println(err)
	}

	extra.RegisterFuzzyDecoders()
	raw = bytes.ReplaceAll(raw, []byte("'"), []byte("\""))
	err = jsoniter.Unmarshal(raw, &abslist)
	if err != nil {
		fmt.Println(err)
	}
	var songs []Song
	for _, v := range abslist.Songs {
		songs = append(songs, k.Song(v.Musicrid))
	}
	return songs
}

// Song Song
func (k Kuwo) Song(rid string) Song {
	type SongXML struct {
		MusicID   string `xml:"music_id"`
		Name      string `xml:"name"`
		Artist    string `xml:"artist"`
		ArtistPic string `xml:"artist_pic"`
		Mp3dl     string `xml:"mp3dl"`
		Aacdl     string `xml:"aacdl"`
		Aacpath   string `xml:"aacpath"`
		Mp3path   string `xml:"mp3path"`
	}
	palyurl := "http://player.kuwo.cn/webmusic/st/getNewMuiseByRid?rid=%s"
	spider := Spider{Type: "GET", Referer: "http://player.kuwo.cn/webmusic/play", URL: fmt.Sprintf(palyurl, rid)}
	raw, err := spider.FetchBytes()
	if err != nil {
		fmt.Println(err)
	}
	raw = bytes.ReplaceAll(raw, []byte("&"), []byte("#"))
	var song SongXML
	err = xml.Unmarshal(raw, &song)
	if err != nil {
		fmt.Println(err)
	}

	music := Song{WebSite: "kuwo", ID: rid, Title: song.Name, Artist: song.Artist, Cover: song.ArtistPic}

	if song.Mp3dl != "" {
		music.URL = fmt.Sprintf("http://%s/resource/%s", song.Mp3dl, song.Mp3path)
	} else {
		music.URL = fmt.Sprintf("http://%s/resource/%s", song.Aacdl, song.Aacpath)
	}
	music.URL = strings.ReplaceAll(music.URL, "#", "&")

	return music
}
