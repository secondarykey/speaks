package logic

import (
	"fmt"

	ldap "gopkg.in/ldap.v3"
)

type LDAP struct {
	Protocol     string
	Host         string
	Port         int
	BaseDN       string
	BindUser     string
	BindPassword string
}

type LDAPResult struct {
	LoginName string
	Name      string
}

func NewLDAP(host string, bDN string, u string, p string) *LDAP {
	return &LDAP{
		Protocol:     "tcp",
		Host:         host,
		Port:         389,
		BaseDN:       bDN,
		BindUser:     u,
		BindPassword: p,
	}
}

//ActiveDirectoryを想定
func (l *LDAP) Login(u string, p string) (*LDAPResult, error) {

	s, err := ldap.Dial(l.Protocol, fmt.Sprintf("%s:%d", l.Host, l.Port))
	if err != nil {
		return nil, fmt.Errorf("LDAP Dial Error[%v]", err)
	}
	defer s.Close()

	err = s.Bind(l.BindUser, l.BindPassword)
	if err != nil {
		return nil, fmt.Errorf("Bind Error[%v]", err.Error())
	}

	account := fmt.Sprintf("(sAMAccountName=%s)", u)
	rtnType := []string{"sAMAccountName", "displayName"}

	req := ldap.NewSearchRequest(l.BaseDN, ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases, 0, 0, false,
		account, rtnType, nil)

	sr, err := s.Search(req)
	if err != nil {
		return nil, fmt.Errorf("Search Error[%v]", err)
	}

	if len(sr.Entries) == 0 {
		return nil, fmt.Errorf("Not Found:%s", u)
	} else if len(sr.Entries) > 1 {
		return nil, fmt.Errorf("many found... user: %s", u)
	}

	entity := sr.Entries[0]
	err = s.Bind(entity.DN, p)
	if err != nil {
		return nil, fmt.Errorf("Authorization Error[%s]", u)
	}

	result := LDAPResult{
		LoginName: entity.GetAttributeValue("sAMAccountName"),
		Name:      entity.GetAttributeValue("displayName"),
	}

	return &result, nil
}
