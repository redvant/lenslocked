package main

import (
	"fmt"

	"github.com/redvant/lenslocked/models"
)

func main() {
	is := models.ImageService{}
	fmt.Println(is.Images(1))
}
