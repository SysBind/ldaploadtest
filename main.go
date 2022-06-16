package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	ldap "github.com/go-ldap/ldap/v3"
)

func main() {
	ldapURL := flag.String("u", "ldap://localhost:389", "LDAP Server URL")
	cn := flag.String("c", "admin", "Common Name (i.e: login username)")
	dn := flag.String("dn", "sysbind.test", "Distinguished Name (i.e: domain name)")
	pass := flag.String("p", "p4ssw0rd", "Password")
	flag.Parse()

	l, err := ldap.DialURL(*ldapURL)
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	bindStr := fmt.Sprintf("cn=%s", *cn)
	dcs := strings.Split(*dn, ".")
	for i := range dcs {
		bindStr = fmt.Sprintf("%s, dc=%s", bindStr, dcs[i])
	}
	fmt.Printf("attempting to connect using %s\n", bindStr)

	err = l.Bind(bindStr, *pass)
	if err != nil {
		log.Fatal(err)
	}
	user := "fooUser"
	baseDN := "DC=sysbind,DC=test"
	filter := fmt.Sprintf("(CN=%s)", ldap.EscapeFilter(user))

	// Filters must start and finish with ()!
	searchReq := ldap.NewSearchRequest(baseDN, ldap.ScopeWholeSubtree, 0, 0, 0, false, filter, []string{"sAMAccountName"}, []ldap.Control{})

	result, err := l.Search(searchReq)
	if err != nil {
		fmt.Errorf("failed to query LDAP: %w", err)
	}
	log.Println("Got", len(result.Entries), "search results")
}
