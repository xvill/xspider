package xspider

import (
	"bytes"
	"crypto/md5"
	"errors"
	"fmt"
	"go-common/app/service/ops/log-agent/pkg/bufio"
	"math/rand"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/grafov/m3u8"
	"github.com/levigross/grequests"
	"github.com/xvill/xutil"
)

// https://github.com/storyicon/golang-proxy
// go get -u github.com/gocolly/colly/...

type Spider struct {
	Type    string
	URL     string
	UA      string
	Referer string
	Params  map[string]string
	Proxy   *url.URL
	ro      grequests.RequestOptions
}

func (s *Spider) init() {
	if s.Type == "" {
		s.Type = "GET"
	}
	if s.UA == "" {
		s.UA = genRandUA()
	}
	if s.Proxy != nil {
		s.ro.Proxies = map[string]*url.URL{s.Proxy.Scheme: s.Proxy}
	}
	s.ro = grequests.RequestOptions{
		Params: s.Params,
		Headers: map[string]string{
			"Referer":    s.Referer,
			"User-Agent": s.UA}}
}

//Fetch Fetch
func (s *Spider) Fetch() (*grequests.Response, error) {
	s.init()
	switch strings.ToUpper(s.Type) {
	case "POST":
		return grequests.Post(s.URL, &s.ro)
	case "GET":
		return grequests.Get(s.URL, &s.ro)
	}
	return nil, errors.New("Spider.Type not support")
}

//FetchBytes FetchBytes
func (s *Spider) FetchBytes() ([]byte, error) {
	res, err := s.Fetch()
	if err != nil {
		return nil, err
	}
	defer res.Close()
	return res.Bytes(), nil
}

//JSON JSON
func (s *Spider) JSON(userStruct interface{}) error {
	res, err := s.Fetch()
	if err != nil {
		return err
	}
	defer res.Close()
	return res.JSON(userStruct)
}

//Download Download
func (s *Spider) Download(savepath string) error {
	res, err := s.Fetch()
	if err != nil {
		return err
	}

	err = res.DownloadToFile(savepath)
	if err != nil {
		return err
	}
	return nil
}

//=======================================================================================================

// DownloadSong DownloadSong
func (s *Spider) DownloadSong(song Song, savepath string) error {
	if song.URL == "" {
		return errors.New("empty url")
	}

	s.URL = song.URL
	err := s.Download(savepath)
	if err != nil {
		return err
	}

	if song.Cover != "" && len(song.CoverImage) == 0 {
		s.URL = song.Cover
		res, _ := s.Fetch()
		song.CoverImage = res.Bytes()
	}
	return song.UpdateFile(savepath)
}

//=======================================================================================================

//HLSM3u8 HLSM3u8
type HLSM3u8 struct {
	URL       string
	SubURL    []string
	KeyMethod string
	KeyURI    string
	Segment   []string
	Raw       []byte
}

//FetchHLSM3u8 FetchHLSM3u8
func (s *Spider) FetchHLSM3u8(inurl string) (h HLSM3u8, err error) {
	if inurl == "" {
		return h, errors.New("empty url")
	}
	s.URL = inurl
	raw, err := s.FetchBytes()
	if err != nil {
		return h, err
	}
	h.URL = inurl
	h.Raw = raw

	playlist, listType, err := m3u8.DecodeFrom(bytes.NewReader(raw), true)
	if err != nil {
		return h, err
	}

	switch listType {
	case m3u8.MASTER:
		masterpl := playlist.(*m3u8.MasterPlaylist)
		for _, variant := range masterpl.Variants {
			if variant != nil {
				h.SubURL = append(h.SubURL, variant.URI)
			}
		}
		return
	case m3u8.MEDIA:
		mediapl := playlist.(*m3u8.MediaPlaylist)
		h.KeyMethod = mediapl.Key.Method
		h.KeyURI = mediapl.Key.URI
		h.Segment = make([]string, 0)
		for _, segment := range mediapl.Segments {
			if segment != nil {
				h.Segment = append(h.Segment, segment.URI)
			}
		}
	default:
		return h, errors.New("Not a valid playlist")
	}

	HostArr := strings.Split(h.URL, "/")
	if !strings.HasPrefix(h.KeyURI, "http://") && !strings.HasPrefix(h.KeyURI, "http://") {
		HostArr[len(HostArr)-1] = h.KeyURI
		h.KeyURI = strings.Join(HostArr, "/")
	}

	for i, v := range h.Segment {
		if strings.HasPrefix(v, "http://") || strings.HasPrefix(v, "http://") {
			continue
		}
		HostArr[len(HostArr)-1] = v
		h.Segment[i] = strings.Join(HostArr, "/")
	}
	return h, nil
}

// DownloadHLSM3u8 DownloadHLSM3u8
func (s *Spider) DownloadHLSM3u8(h HLSM3u8, savepath string) error {
	s.URL = h.KeyURI
	rawkey, err := s.FetchBytes()
	if err != nil {
		return err
	}
	c := xutil.NewCrypto(rawkey)

	mergefile := savepath
	isExist, isDir, err := xutil.IsFileExist(savepath)
	if !isExist {
		return errors.New(savepath + " not found")
	}
	if !isDir {
		mergefile = fmt.Sprintf("%s/%x.mp4", filepath.Dir(savepath), md5.Sum([]byte(h.URL)))
	}

	mergef, err := os.OpenFile(mergefile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0660)
	if err != nil {
		return err
	}
	defer mergef.Close()
	bw := bufio.NewWriter(mergef)

	for _, v := range h.Segment {
		s.URL = v
		subraw, err := s.FetchBytes()
		if err != nil {
			fmt.Println("Download", v, err.Error())
			continue
		}
		wraw := make([]byte, 0)
		if h.KeyURI != "" {
			dst, err := c.Decrypt(subraw)
			if err != nil {
				fmt.Println("Decrypt", v, err.Error())
				continue
			}
			wraw = dst
		} else {
			wraw = subraw
		}
		_, err = bw.Write(wraw)
		if err != nil {
			fmt.Println("WriteFile", v, err.Error())
			continue
		}
	}
	bw.Flush()
	return nil
}

//DownloadHLSM3u8FromURL DownloadHLSM3u8FromURL
func (s *Spider) DownloadHLSM3u8FromURL(inurl, savepath string) error {
	hlsm3u8, err := s.FetchHLSM3u8(inurl)
	if err != nil {
		return err
	}
	return s.DownloadHLSM3u8(hlsm3u8, savepath)
}

//=======================================================================================================

func genRandUA() string {
	var ffVersions = []float32{58.0, 57.0, 56.0, 52.0, 48.0, 40.0, 35.0}
	var chromeVersions = []string{"65.0.3325.146", "64.0.3282.0", "41.0.2228.0", "40.0.2214.93", "37.0.2062.124"}
	var osStrings = []string{
		"Macintosh; Intel Mac OS X 10_10",
		"Windows NT 10.0",
		"Windows NT 5.1",
		"Windows NT 6.1; WOW64",
		"Windows NT 6.1; Win64; x64",
		"X11; Linux x86_64",
	}
	if rand.New(rand.NewSource(time.Now().UnixNano())).Intn(10) > 5 {
		version := ffVersions[rand.Intn(len(ffVersions))]
		os := osStrings[rand.Intn(len(osStrings))]
		return fmt.Sprintf("Mozilla/5.0 (%s; rv:%.1f) Gecko/20100101 Firefox/%.1f", os, version, version)
	}
	version := chromeVersions[rand.Intn(len(chromeVersions))]
	os := osStrings[rand.Intn(len(osStrings))]
	return fmt.Sprintf("Mozilla/5.0 (%s) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/%s Safari/537.36", os, version)
}
