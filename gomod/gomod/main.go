package main

import (
	"flag"
	"log"

	"go.x2ox.com/sorbifolia/gomod"
)

var filename string

func init() {
	flag.StringVar(&filename, "filename", "data.json", "")
	flag.StringVar(&filename, "f", "data.json", "")
}

func main() {
	flag.Parse()

	pkg, err := gomod.Parse(filename)
	if err != nil {
		log.Panicln(err)
	}
	for _, v := range pkg {
		var ms []string
		if ms, err = v.FindModule(); err != nil {
			log.Println(err)
		}

		for _, m := range ms {
			if err = v.Output(m); err != nil {
				log.Println(m, err)
			}
		}

		if err = v.Clean(); err != nil {
			log.Println(err)
		}
	}
}
