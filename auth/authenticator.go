package auth

import (
	"errors"
	"sanguosha.com/sgs_herox/proto/gameconf"
	gamedef "sanguosha.com/sgs_herox/proto/def"
)

var ErrTicketInvalid = errors.New("ticket not valid")
var ErrLoginTypeNotSupport = errors.New("loginType not support")

type Authenticator interface {
	Auth() (*gamedef.AuthInfo, string, error)
}

func NewAuthenticator(typ gameconf.AccountLoginTyp, token string) (Authenticator, error) {
	switch typ {
	case gameconf.AccountLoginTyp_ALTTest:
		return &InsideTestAuthenticator{
			token: token,
		}, nil
	case gameconf.AccountLoginTyp_ALTTablePark:
		return &TableParkAuthenticator{
			token: token,
		}, nil
	default:
		return nil, ErrLoginTypeNotSupport
	}
}
