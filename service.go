package main

import (
	"fmt"
	"log"
	"strings"

	ldap "github.com/go-ldap/ldap/v3"
)

type Service struct {
	url       string
	cn        string
	dn        string
	pass      string
	dcsString string
	conn      *ldap.Conn
}

func (svc *Service) Bind() (err error) {
	svc.conn, err = ldap.DialURL(svc.url)
	if err != nil {
		return err
	}
	bindStr := fmt.Sprintf("cn=%s", svc.cn)

	// split dn components
	dcs := strings.Split(svc.dn, ".")
	for i := range dcs {
		svc.dcsString = fmt.Sprintf("%s dc=%s", svc.dcsString, dcs[i])
		if i < len(dcs)-1 {
			svc.dcsString = fmt.Sprintf("%s,", svc.dcsString)
		}
	}
	bindStr = fmt.Sprintf("%s, %s", bindStr, svc.dcsString)
	fmt.Printf("connecting using \"%s\"\n", bindStr)

	err = svc.conn.Bind(bindStr, svc.pass)
	if err != nil {
		return err
	}
	return nil
}

func (svc *Service) Query(user string) error {
	baseDN := svc.dcsString
	filter := fmt.Sprintf("(CN=%s)", ldap.EscapeFilter(user))

	// Filters must start and finish with ()!
	searchReq := ldap.NewSearchRequest(baseDN, ldap.ScopeWholeSubtree, 0, 0, 0, false, filter, []string{"sAMAccountName"}, []ldap.Control{})

	result, err := svc.conn.Search(searchReq)
	if err != nil {
		return err
	}
	log.Println("Got", len(result.Entries), "search results")
	return nil
}

func (svc *Service) Close() {
	svc.conn.Close()
}
