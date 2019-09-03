package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	termbox "github.com/nsf/termbox-go"
	"github.com/xvill/xspider"
)

func init() {
	if err := termbox.Init(); err != nil {
		panic(err)
	}
	termbox.SetCursor(0, 0)
	termbox.HideCursor()
}

func main() {
	log.Print("task begin-->>>>>>>>>>>>")
	fname := filepath.Join(filepath.Dir(os.Args[0]), "downloadm3u8.txt")
	raw, err := ioutil.ReadFile(fname)
	if err != nil {
		log.Println(err)
		ioutil.WriteFile(fname, []byte{}, 0666)
		log.Println("create ", fname)
	} else {
		bs := bufio.NewScanner(bytes.NewReader(raw))
		var spider xspider.Spider
		inurl, savepath := "", filepath.Base(fname)
		fmap := make(map[string]string, 0)
		cnt := 0
		for bs.Scan() {
			cnt++
			inurl = bs.Text()
			if _, exists := fmap[inurl]; exists {
				continue
				log.Println(cnt, inurl, " download skip")
			}
			fmap[inurl] = ""
			log.Println(cnt, inurl, " download start")
			err := spider.DownloadHLSM3u8FromURL(inurl, savepath)
			if err != nil {
				log.Println(cnt, inurl, err)
			} else {
				log.Println(cnt, inurl, " download success")
			}
		}
	}
	log.Print("task end--<<<<<<<<<<<<<")
	pause()
}

func pause() {
	fmt.Println("------------- Press any key to continue -------------")
Loop:
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			break Loop
		}
	}
}
