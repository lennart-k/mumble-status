package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/lennart-k/mumble-status/pkg/status"
)

func main() {
	addressPtr := flag.String("address", "", "address")
	flag.Parse()
	if *addressPtr == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	status, err := status.GetServerStatus(*addressPtr)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", status)
}
