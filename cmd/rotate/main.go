package main

import (
	"fmt"

	_ "github.com/moov-io/identity" // need to import the embedded files

	"github.com/moov-io/identity/pkg/webkeys"
)

/*
This is used to rotate the jwks by our systems. It will generate a RSA256 public and private key for used.
*/
func main() {
	js, err := webkeys.NewGenerateJwksService()
	if err != nil {
		fmt.Printf("Unable to generate JWKS service - %s\n", err.Error())
		return
	}

	s, _ := js.(*webkeys.GenerateJwksService)

	fmt.Printf("Writing out generated JWKS files into ./\n")
	s.Save("./configs/")
}
