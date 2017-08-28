package main

import (
	//"os"
	"fmt"
	"youtube"

	log "github.com/Sirupsen/logrus"
)

func main() {
	//id := os.Args[1]
	//filename := id+".mp4"
	//video, err := youtube.GetVideo(id)
	//if err != nil {
	//	log.Fatal(err)
	//}

	videos, err := youtube.Search("lukas graham apologize", 5)
	if err != nil {
		log.Fatalf("search failed: %s", err)
	}

	for _, v := range videos {
		fmt.Printf("%s\n", v.ID)
		fmt.Printf("%s\n", v.Title)
		fmt.Printf("%s\n", v.Thumbnail)
		fmt.Printf("%d views\n", v.Views)
		fmt.Printf("%d seconds\n", v.Length)
		fmt.Printf("%.1f/5\n", v.Rating)
		fmt.Printf("\n")

		//s := v.Streams[0]
		//fmt.Printf("%s %s %s %s %d kbit/s %s\n", s.Extension, s.Resolution, s.VideoEncoding, s.AudioEncoding, s.AudioBitrate, s.URL)

	}
	// if err := v.Download(filename); err != nil { log.Fatal(err) }
}
