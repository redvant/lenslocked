package main

import (
	"fmt"

	"github.com/redvant/lenslocked/models"
)

func main() {
	gs := models.GalleryService{}
	fmt.Println(gs.Images(1))
}
