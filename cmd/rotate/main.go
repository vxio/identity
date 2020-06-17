package main

import (
	"fmt"

	"github.com/moov-io/tumbler/pkg/webkeys"
)

/*
This is used to rotate the jwks by our systems. It will generate a RSA256 public and private key for used.
*/
func main() {
	s, err := webkeys.NewGenerateJwksService()
	if err != nil {
		fmt.Printf("Unable to generate JWKS service - %s\n", err.Error())
		return
	}

	fmt.Printf("Writing out generated JWKS files into ./\n")
	s.Save("./configs/")
}
