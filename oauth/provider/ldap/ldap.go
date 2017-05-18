package ldap

import (
	"fmt"
	. "github.com/ory-am/hydra/oauth/provider"
	"gopkg.in/ldap.v2"
	"crypto/tls"
)

type ldapconf struct {
	id       string
	host     string
	port     string
	tls      bool
	base_dn  string
	bind_dn  string
	uid      string
	password string
	filter   string
}

func New(id, client, secret, redirectURL string) *ldapconf {
	return &ldapconf{
		id:  id,
		host: "ldap.forumsys.com",
		tls: false,
		port: "389",
		bind_dn: "cn=read-only-admin,dc=example,dc=com",
		base_dn: "dc=example,dc=com",
		uid: "uid",
		password: "password",
		filter: "(objectClass=organizationalPerson)",
	}
}

func (l *ldapconf) GetAuthenticationURL(state string) string {
	return ""
}

func (l *ldapconf) FetchSession(code string) (Session, error) {
	var connection *ldap.Conn
	var err error

	if !l.tls {
		connection, err = ldap.Dial("tcp", fmt.Sprintf("%s:%s", l.host, l.port))
	} else {
		connection, err = ldap.DialTLS("tcp", fmt.Sprintf("%s:%s", l.host, l.port), &tls.Config{InsecureSkipVerify: true})
	}

	if err != nil {
		return nil, err
	}

	defer connection.Close()

	if err = connection.Bind(l.bind_dn, l.password); err != nil {
		return nil, err
	}

	searchRequest := ldap.NewSearchRequest(
		l.base_dn,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(&%s(%s=%s))", l.filter, l.uid, code),
		[]string{l.uid},
		nil,
	)

	sr, err := connection.Search(searchRequest)
	if err != nil {
		return nil, err
	}

	if len(sr.Entries) != 1 {
		return nil, err
	}

	acc := make(map[string]interface{}, 10)
	for _, attr := range sr.Entries[0].Attributes {
		acc[attr.Name] = attr.Values
	}

	return &DefaultSession{
		RemoteSubject: fmt.Sprintf("%s", sr.Entries[0].DN),
		Extra:         acc,
	}, nil
}

func (l *ldapconf) GetID() string {
	return l.id
}