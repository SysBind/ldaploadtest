package main

import (
	"flag"
	"fmt"
)

func main() {
	ldapURL := flag.String("u", "ldap://localhost:389", "LDAP Server URL")
	cn := flag.String("c", "admin", "Common Name (i.e: login username)")
	dn := flag.String("dn", "sysbind.test", "Distinguished Name (i.e: domain name)")
	pass := flag.String("p", "p4ssw0rd", "Password")
	flag.Parse()

	svc := Service{url: *ldapURL, cn: *cn, dn: *dn, pass: *pass}
	err := svc.Bind()
	if err != nil {
		fmt.Errorf("LDAP bind failed: %w", err)
	}
	defer svc.Close()

	loader := Loader{svc: &svc}

	err = loader.Run()
	if err != nil {
		fmt.Errorf("failed to run: %w", err)
	}
}
