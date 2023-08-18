package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/lennart-k/mumble-status/pkg/status"
)

func main() {
	hostPtr := flag.String("h", "", "hostname")
	portPtr := flag.Uint("p", 64738, "port")
	flag.Parse()
	if *hostPtr == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	status, err := status.GetServerStatus(*hostPtr, *portPtr)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", status)
}
