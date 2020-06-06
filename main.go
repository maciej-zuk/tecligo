// Terraria bot
package main

import (
	"io"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/maciej-zuk/tecligo/temap"
	"github.com/maciej-zuk/tecligo/tenet"
)

func botRun() {
}

func main() {
	file, err := os.OpenFile("info.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	multi := io.MultiWriter(file, os.Stdout)
	log.SetOutput(multi)

	defer func() {
		if r := recover(); r != nil {
			log.Println("Panic")
			log.Println(r)
		} else {
			log.Println("Exit (main)")
		}
	}()

	temap.InitTextures()
	defer temap.DestroyTextures()

	rand.Seed(time.Now().UTC().UnixNano())
	c := tenet.Connection{}
	c.Connect()
}
