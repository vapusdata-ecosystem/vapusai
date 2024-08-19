package main

import (
	"log"

	vapuspublish "github.com/vapusdata-ecosystem/vapusai-studio/scripts/goscripts/publish"
)

func main() {
	chartVersion := vapuspublish.PushOciImages()
	log.Println(chartVersion)
}
