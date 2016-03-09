package ldap

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

var mock = ldapconf{
	id:  "123",
	host: "ldap://ldap.forumsys.com",
	port: "389",
	bind_dn: "cn=read-only-admin,dc=example,dc=com",
	base_dn: "dc=example,dc=com",
	uid: "uid",
	password: "password",
	filter: "(objectClass=organizationalPerson)",
}

func TestNew(t *testing.T) {
	m := New("321", "client", "secret", "/callback")
	assert.Equal(t, "321", m.id)

}

func TestGetID(t *testing.T) {
	assert.Equal(t, "123", mock.GetID())
}

func TestFetchSession(t *testing.T) {
	m := New("321", "client", "secret", "/callback")
	ses, err := m.FetchSession("riemann")

	assert.NoError(t, err)
	assert.Equal(t, "uid=riemann,dc=example,dc=com", ses.GetRemoteSubject())
	//fmt.Printf(ses.GetExtra())
}
