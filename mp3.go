package xspider

import (
	"fmt"

	"github.com/bogem/id3v2"
)

//Song Song
type Song struct {
	WebSite     string `json:"website"`
	Title       string `json:"title"`
	ID          string `json:"id"`
	Alias       string `json:"alias"`
	Artist      string `json:"artist"`
	URL         string `json:"src"`
	Cover       string `json:"pic"`
	Lyric       string `json:"lrc"`
	Album       string `json:"album"`
	ArtistID    string `json:"artist_id,omitempty"`
	AlbumID     int    `json:"albumid,omitempty"`
	CoverImage  []byte `json:"-"`
	Year        string `json:"year,omitempty"`
	Publishtime string `json:"publishtime,omitempty"`
	Lyclng      string `json:"lyclng,omitempty"`
	Lycdesc     string `json:"lycdesc,omitempty"`
	Track       int    `json:"track,omitempty"`
	Size        int    `json:"size,omitempty"`
	Encodetype  string `json:"encodetype,omitempty"`
	Basicrate   int    `json:"basicrate,omitempty"`
}

// Desc Desc
func (s *Song) Desc() map[string]string {

	ret := map[string]string{
		"website": s.WebSite,
		"id":      s.ID,
		"title":   s.Title,
		"artist":  s.Artist,
		"src":     s.URL,
		"pic":     s.Cover,
		"lrc":     s.Lyric,
		"album":   s.Album}
	if s.Alias != "" {
		ret["title"] = fmt.Sprintf("%s-(%s)", ret["title"], s.Alias)
	}
	return ret
}

//FromFile FromFile
func (s *Song) FromFile(fname string) error {
	tag, err := id3v2.Open(fname, id3v2.Options{Parse: true})
	if err != nil {
		return err
	}
	s.Album = tag.Album()
	s.Artist = tag.Artist()
	s.Title = tag.Title()
	s.Year = tag.Year()

	pictures := tag.GetFrames(tag.CommonID("Attached picture"))
	if len(pictures) > 0 {
		pic, ok := pictures[0].(id3v2.PictureFrame)
		if ok {
			s.CoverImage = pic.Picture
		}
	}
	lyrics := tag.GetFrames(tag.CommonID("Unsynchronised lyrics/text transcription"))
	if len(lyrics) > 0 {
		lrc, ok := lyrics[0].(id3v2.UnsynchronisedLyricsFrame)
		if ok {
			s.Lyclng = lrc.Language
			s.Lycdesc = lrc.ContentDescriptor
			s.Lyric = lrc.Lyrics
		}
	}
	return nil
}

// UpdateFile 更新标签
func (s Song) UpdateFile(fname string) error {
	tag, err := id3v2.Open(fname, id3v2.Options{Parse: true})
	if err != nil {
		return err
	}
	defer tag.Close()

	tag.SetDefaultEncoding(id3v2.EncodingUTF8)
	tag.SetTitle(s.Title)
	tag.SetArtist(s.Artist)
	tag.SetAlbum(s.Album)
	tag.SetYear(s.Year)

	if len(s.CoverImage) > 0 {
		pf := id3v2.PictureFrame{
			Encoding:    id3v2.EncodingUTF8,
			MimeType:    "image/jpg",
			PictureType: id3v2.PTOther,
			Picture:     s.CoverImage}

		tag.AddAttachedPicture(pf)
	}
	if s.Lyric != "" {
		if s.Lyclng == "" {
			s.Lyclng = "zho"
		}
		uslf := id3v2.UnsynchronisedLyricsFrame{
			Encoding:          id3v2.EncodingUTF8,
			Language:          s.Lyclng, // zho 中/eng 英
			ContentDescriptor: s.Lycdesc,
			Lyrics:            s.Lyric}
		tag.AddUnsynchronisedLyricsFrame(uslf)
	}
	return tag.Save()
}

/***
//TODO
xspider


****/
