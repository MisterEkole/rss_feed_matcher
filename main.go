package main

import (
	"log"
	_ "M1/matcher"
	"os"
	"M1/search"
)

func init() {
	log.SetOutput(os.Stdout)
}

func main() {
	search.Run("Deputy ministers")
}
