package main

import (
	"flag"
	"fmt"
	"github.com/feimingxliu/quicksearch/pkg/util/bcrypt"
	"log"
)

var plaintext *string

func init() {
	plaintext = flag.String("p", "admin", "the plaintext password")
}

func main() {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*plaintext))
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("%s\n", hashedPassword)
}
