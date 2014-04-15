package main

import (
	"flag"
	"fmt"
	"github.com/davidnarayan/go-flake"
	"log"
)

var max = flag.Int("max", 1, "number of IDs to create")
var hex = flag.Bool("hex", false, "Show hex representation")
var integer = flag.Bool("integer", false, "Show integer representation")

func main() {
	flag.Parse()
	f, err := flake.New()
	if err != nil {
		log.Fatal(err)
	}

	if !*hex && !*integer {
		*hex = true
	}

	for i := 0; i < *max; i++ {
		id := f.NextId()

		if *integer {
			fmt.Println(id)
		}

		if *hex {
			fmt.Println(id.String())
		}
	}
}
