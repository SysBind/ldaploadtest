package main

import (
	"flag"
	"fmt"
	"log"

	ldap "github.com/go-ldap/ldap/v3"
)

func main() {
	ldapURL := flag.String("u", "ldap://localhost:389", "LDAP Server URL")
	l, err := ldap.DialURL(*ldapURL)
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()
	fmt.Println("LDAP Load Testing Tool")
}
