package main

import (
	"fmt"
	"github.com/streamlist/streamlist/internal/youtube"
	"os"
)

func main() {
	youtube.SetDebug()

	videos, err := youtube.Search(os.Args[1])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, v := range videos {
		fmt.Printf("%s\n", v.ID)
		fmt.Printf("%s\n", v.Title)
		fmt.Printf("%s\n", v.Thumbnail)
		fmt.Printf("%d seconds\n", v.Length)
	}
}
