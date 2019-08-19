package xspider

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/json-iterator/go"
)

type neteaseSearch struct {
	Result struct {
		Songs []struct {
			ID          int      `json:"id"`
			Name        string   `json:"name"`
			Track       int      `json:"no"`
			PublishTime int64    `json:"publishTime"`
			Alias       []string `json:"alias"`
			Ar          []struct {
				ID   int    `json:"id"`
				Name string `json:"name"`
			} `json:"ar"`
			Al struct {
				ID     int    `json:"id"`
				Name   string `json:"name"`
				PicURL string `json:"picUrl"`
			} `json:"al"`
		}
	}
}

type neteaseSong struct {
	Songs []struct {
		ID      int      `json:"id"`
		Name    string   `json:"name"`
		Track   int      `json:"no"`
		Alias   []string `json:"alias"`
		Artists []struct {
			ID     int    `json:"id"`
			Name   string `json:"name"`
			PicURL string `json:"picUrl"`
		} `json:"artists"`
		Album struct {
			Name        string `json:"name"`
			ID          int    `json:"id"`
			Type        string `json:"type"`
			PicURL      string `json:"picUrl"`
			PublishTime int64  `json:"publishTime"`
		} `json:"album"`
	} `json:"songs"`
}

type neteaseSongURL struct {
	Data []struct {
		ID         int      `json:"id"`
		URL        string   `json:"url"`
		Alias      []string `json:"alias"`
		Br         int      `json:"br"`
		Size       int      `json:"size"`
		EncodeType string   `json:"encodeType"`
	} `json:"data"`
}

type Netease struct{}

//Search Search
func (n Netease) Search(keyword string) []Song {
	spider := Spider{
		Type:    "GET",
		URL:     "http://music.163.com/api/cloudsearch/pc",
		Referer: "http://music.163.com/",
		Params: map[string]string{
			"s":      keyword,
			"type":   "1",
			"offset": "0",
			"limit":  "10"}}
	//-----------------------------------
	var songs neteaseSearch
	err := spider.JSON(&songs)
	if err != nil {
		log.Println(err)
	}
	//-----------------------------------
	ret := make([]Song, 0)
	for _, v := range songs.Result.Songs {
		artists := []string{}
		artistsID := []string{}
		for _, a := range v.Ar {
			artists = append(artists, a.Name)
			artistsID = append(artistsID, strconv.Itoa(a.ID))
		}
		s := Song{
			WebSite:     "netease",
			ID:          strconv.Itoa(v.ID),
			Title:       v.Name,
			Track:       v.Track,
			Alias:       strings.Join(v.Alias, "/"),
			Artist:      strings.Join(artists, "/"),
			ArtistID:    strings.Join(artistsID, "/"),
			AlbumID:     v.Al.ID,
			Album:       v.Al.Name,
			Year:        time.Unix(v.PublishTime/1000, 0).Format("2006"),
			Publishtime: time.Unix(v.PublishTime/1000, 0).Format("2006-01-02"),
			Cover:       v.Al.PicURL}

		lyric := n.LRC([]int{v.ID})
		if v, exists := lyric[v.ID]; exists {
			s.Lyric = v
		}
		songurl := n.SongURL([]int{v.ID})
		if v, exists := songurl[v.ID]; exists {
			s.URL = v.URL
			s.Size = v.Size
			s.Encodetype = v.Encodetype
			s.Basicrate = v.Basicrate
		}
		ret = append(ret, s)
	}
	return ret
}

//Playlist Playlist
func (n Netease) Playlist(playlistid int) []Song {
	type PlaylistDATA struct {
		Result struct {
			Tracks []struct {
				ID      int      `json:"id"`
				Name    string   `json:"name"`
				Alias   []string `json:"alias"`
				Artists []struct {
					Name string `json:"name"`
					ID   int    `json:"id"`
				} `json:"artists"`
				Album struct {
					Name        string `json:"name"`
					ID          int    `json:"id"`
					PicURL      string `json:"picUrl"`
					PublishTime int64  `json:"publishTime"`
				} `json:"album"`
			} `json:"tracks"`
		} `json:"result"`
	}
	spider := Spider{
		Type:    "GET",
		URL:     fmt.Sprintf("http://music.163.com/api/playlist/detail?id=%d&ids=[%d]", playlistid, playlistid),
		Referer: "http://music.163.com/"}

	var songs PlaylistDATA
	err := spider.JSON(&songs)
	if err != nil {
		log.Println(err)
	}
	//-----------------------------------
	ret := make([]Song, 0)
	for _, v := range songs.Result.Tracks {
		artists := []string{}
		artistsID := []string{}
		for _, a := range v.Artists {
			artists = append(artists, a.Name)
			artistsID = append(artistsID, strconv.Itoa(a.ID))
		}
		s := Song{
			WebSite:     "netease",
			ID:          strconv.Itoa(v.ID),
			Title:       v.Name,
			Alias:       strings.Join(v.Alias, "/"),
			Artist:      strings.Join(artists, "/"),
			ArtistID:    strings.Join(artistsID, "/"),
			AlbumID:     v.Album.ID,
			Album:       v.Album.Name,
			Year:        time.Unix(v.Album.PublishTime/1000, 0).Format("2006"),
			Publishtime: time.Unix(v.Album.PublishTime/1000, 0).Format("2006-01-02"),
			Cover:       v.Album.PicURL}

		lyric := n.LRC([]int{v.ID})
		if v, exists := lyric[v.ID]; exists {
			s.Lyric = v
		}
		songurl := n.SongURL([]int{v.ID})
		if v, exists := songurl[v.ID]; exists {
			s.URL = v.URL
			s.Size = v.Size
			s.Encodetype = v.Encodetype
			s.Basicrate = v.Basicrate
		}
		ret = append(ret, s)
	}

	return ret
}

//Artist Artist
func (n Netease) Artist(artistid int) []Song {
	type ArtistDATA struct {
		HotSongs []struct {
			ID      int      `json:"id,omitempty"`
			Name    string   `json:"name,omitempty"`
			Alias   []string `json:"alias"`
			Artists []struct {
				Name string `json:"name"`
				ID   int    `json:"id"`
			} `json:"artists,omitempty"`
			Album struct {
				Status      int    `json:"status"`
				PublishTime int64  `json:"publishTime"`
				PicURL      string `json:"picUrl"`
				Name        string `json:"name"`
				ID          int    `json:"id"`
			} `json:"album,omitempty"`
		} `json:"hotSongs"`
	}

	spider := Spider{Type: "GET",
		URL:     fmt.Sprintf("http://music.163.com/api/artist/%d", artistid),
		Referer: "http://music.163.com/"}

	var songs ArtistDATA
	err := spider.JSON(&songs)
	if err != nil {
		log.Println(err)
	}
	//-----------------------------------
	ret := make([]Song, 0)
	for _, v := range songs.HotSongs {
		if v.Album.Status < 0 {
			continue
		}
		artists := []string{}
		artistsID := []string{}
		for _, a := range v.Artists {
			artists = append(artists, a.Name)
			artistsID = append(artistsID, strconv.Itoa(a.ID))
		}
		s := Song{
			WebSite:     "netease",
			ID:          strconv.Itoa(v.ID),
			Title:       v.Name,
			Alias:       strings.Join(v.Alias, "/"),
			Artist:      strings.Join(artists, "/"),
			ArtistID:    strings.Join(artistsID, "/"),
			AlbumID:     v.Album.ID,
			Album:       v.Album.Name,
			Year:        time.Unix(v.Album.PublishTime/1000, 0).Format("2006"),
			Publishtime: time.Unix(v.Album.PublishTime/1000, 0).Format("2006-01-02"),
			Cover:       v.Album.PicURL}

		lyric := n.LRC([]int{v.ID})
		if v, exists := lyric[v.ID]; exists {
			s.Lyric = v
		}
		songurl := n.SongURL([]int{v.ID})
		if v, exists := songurl[v.ID]; exists {
			s.URL = v.URL
			s.Size = v.Size
			s.Encodetype = v.Encodetype
			s.Basicrate = v.Basicrate
		}
		ret = append(ret, s)
	}

	return ret
}

//Album 获取Album
func (n Netease) Album(album int) []Song {
	type AlbumDATA struct {
		Album struct {
			Songs []struct {
				ID      int      `json:"id"`
				Name    string   `json:"name"`
				Alias   []string `json:"alias"`
				Artists []struct {
					Name string `json:"name"`
					ID   int    `json:"id"`
				} `json:"artists"`
			} `json:"songs"`
			PublishTime int64  `json:"publishTime"`
			PicURL      string `json:"picUrl"`
			Name        string `json:"name"`
			ID          int    `json:"id"`
		} `json:"album"`
	}

	spider := Spider{Type: "GET",
		URL:     fmt.Sprintf("http://music.163.com/api/album/%d", album),
		Referer: "http://music.163.com/"}

	var songs AlbumDATA
	err := spider.JSON(&songs)
	if err != nil {
		log.Println(err)
	}
	//-----------------------------------
	ret := make([]Song, 0)

	for _, v := range songs.Album.Songs {
		artists := []string{}
		artistsID := []string{}
		for _, a := range v.Artists {
			artists = append(artists, a.Name)
			artistsID = append(artistsID, strconv.Itoa(a.ID))
		}
		s := Song{
			WebSite:     "netease",
			ID:          strconv.Itoa(v.ID),
			Title:       v.Name,
			Alias:       strings.Join(v.Alias, "/"),
			Artist:      strings.Join(artists, "/"),
			ArtistID:    strings.Join(artistsID, "/"),
			AlbumID:     songs.Album.ID,
			Album:       songs.Album.Name,
			Year:        time.Unix(songs.Album.PublishTime/1000, 0).Format("2006"),
			Publishtime: time.Unix(songs.Album.PublishTime/1000, 0).Format("2006-01-02"),
			Cover:       songs.Album.PicURL}

		lyric := n.LRC([]int{v.ID})
		if v, exists := lyric[v.ID]; exists {
			s.Lyric = v
		}
		songurl := n.SongURL([]int{v.ID})
		if v, exists := songurl[v.ID]; exists {
			s.URL = v.URL
			s.Size = v.Size
			s.Encodetype = v.Encodetype
			s.Basicrate = v.Basicrate
		}
		ret = append(ret, s)
	}
	return ret
}

//Song 获取Song
func (n Netease) Song(songids []int) map[int]Song {
	songidsarr := []string{}
	ret := make(map[int]Song, 0)
	for _, v := range songids {
		songidsarr = append(songidsarr, strconv.Itoa(v))
	}
	//-----------------------------------

	spider := Spider{Type: "GET",
		URL:     fmt.Sprintf("http://music.163.com/api/song/detail?ids=[%s]", strings.Join(songidsarr, ",")),
		Referer: "http://music.163.com/"}

	var songs neteaseSong
	err := spider.JSON(&songs)
	if err != nil {
		log.Println(err)
	}

	//-----------------------------------
	for _, v := range songs.Songs {
		artists := []string{}
		artistsID := []string{}
		for _, a := range v.Artists {
			artists = append(artists, a.Name)
			artistsID = append(artistsID, strconv.Itoa(a.ID))
		}
		ret[v.ID] = Song{
			WebSite:     "netease",
			ID:          strconv.Itoa(v.ID),
			Title:       v.Name,
			Track:       v.Track,
			Alias:       strings.Join(v.Alias, "/"),
			Artist:      strings.Join(artists, "/"),
			ArtistID:    strings.Join(artistsID, "/"),
			AlbumID:     v.Album.ID,
			Album:       v.Album.Name,
			Year:        time.Unix(v.Album.PublishTime/1000, 0).Format("2006"),
			Publishtime: time.Unix(v.Album.PublishTime/1000, 0).Format("2006-01-02"),
			Cover:       v.Album.PicURL}
	}

	for k, v := range n.SongURL(songids) {
		s := ret[k]
		s.URL = v.URL
		s.Size = v.Size
		s.Encodetype = v.Encodetype
		s.Basicrate = v.Basicrate
		ret[k] = s
	}

	for k, v := range n.LRC(songids) {
		s := ret[k]
		s.Lyric = v
		ret[k] = s
	}

	return ret
}

//SongURL 获取真实播放地址
func (n Netease) SongURL(songids []int) map[int]Song {
	songidsarr := []string{}
	for _, v := range songids {
		songidsarr = append(songidsarr, strconv.Itoa(v))
	}
	//-----------------------------------
	spider := Spider{Type: "POST",
		URL:     "http://music.163.com/api/song/enhance/player/url",
		Referer: "http://music.163.com/",
		Params: map[string]string{
			"br":  "32000",
			"ids": "[" + strings.Join(songidsarr, ",") + "]"}}

	var songs neteaseSongURL
	err := spider.JSON(&songs)
	if err != nil {
		fmt.Println(err)
	}
	//-----------------------------------
	ret := make(map[int]Song, 0)
	for _, v := range songs.Data {
		ret[v.ID] = Song{
			ID:         strconv.Itoa(v.ID),
			WebSite:    "netease",
			Alias:      strings.Join(v.Alias, "/"),
			URL:        v.URL,
			Size:       v.Size,
			Encodetype: v.EncodeType,
			Basicrate:  v.Br}
	}
	return ret
}

//LRC 获取歌词
func (n Netease) LRC(songids []int) map[int]string {
	URI := "http://music.163.com/api/song/lyric?id=%d&lv=1"
	spider := Spider{Type: "GET", Referer: "http://music.163.com/"}

	lyrics := make(map[int]string, 0)
	for _, songid := range songids {
		spider.URL = fmt.Sprintf(URI, songid)
		raw, err := spider.FetchBytes()
		if err != nil {
			fmt.Println(err)
		}
		lyric := jsoniter.Get(raw, "lrc", "lyric").ToString()
		lyrics[songid] = lyric
	}
	return lyrics
}

/***
// TODO
[]会员歌曲去除

***/
// http://music.163.com/api/playlist/detail?id=2414026629&ids=[2414026629]
// http://music.163.com/api/song/detail?ids=[20789751]
// http://music.163.com/api/song/detail?id=32628933&ids=[32628933]
// http://music.163.com/api/song/detail?id=569213220&ids=[569213220,32628933]
// http://music.163.com/api/song/lyric?id=32628933&lv=1
