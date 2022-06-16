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

	// split dn components
	var dcsstr string
	dcs := strings.Split(*dn, ".")
	for i := range dcs {
		dcsstr = fmt.Sprintf("%s dc=%s", dcsstr, dcs[i])
		fmt.Printf("length=%d, i=%d", len(dcs), i)
		if i < len(dcs)-1 {
			dcsstr = fmt.Sprintf("%s,", dcsstr)
		}
	}
	bindStr = fmt.Sprintf("%s, %s", bindStr, dcsstr)
	fmt.Printf("connecting using \"%s\"\n", bindStr)

	err = l.Bind(bindStr, *pass)
	if err != nil {
		log.Fatal(err)
	}
	user := "demo1"
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
