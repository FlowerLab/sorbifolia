package main

import (
	"os"
	"time"
)

type Meta struct {
	Time time.Time `json:"time"`
	Name string    `json:"name"`
	Tag  string    `json:"tag"`
}

var meta = &Meta{
	Time: time.Now(),
	Name: os.Getenv("_SB_FA_NAME"),
	Tag:  os.Getenv("_SB_FA_TAG"),
}
