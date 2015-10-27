package dropbox

//
//import (
//	"github.com/arekkas/flitt/provider"
//	"github.com/arekkas/flitt/storage"
//	"github.com/go-errors/errors"
//	"github.com/ory-am/common/env"
//	"github.com/stacktic/dropbox"
//	"golang.org/x/oauth2"
//	"log"
//	"strconv"
//)
//
//var (
//	config = &oauth2.Config{
//		ClientID:     env.Getenv("OAUTH_DROPBOX_CLIENT", ""),
//		ClientSecret: env.Getenv("OAUTH_DROPBOX_SECRET", ""),
//		RedirectURL:  env.Getenv("OAUTH_DROPBOX_CALLBACK", ""),
//		Endpoint: oauth2.Endpoint{
//			AuthURL:  "https://www.dropbox.com/1/oauth2/authorize",
//			TokenURL: "https://www.dropbox.com/1/oauth2/token",
//		},
//	}
//)
//
//type dropboxProvider struct {
//	db    *dropbox.Dropbox
//	store *storage.Storage
//}
//
//func New(store *storage.Storage) provider.Provider {
//	db := dropbox.NewDropbox()
//	db.SetAppInfo(client, secret)
//	db.SetRedirectURL(redirect)
//	return &dropboxProvider{db, store}
//}
//
//func (p *dropboxProvider) GetOAuthConfig() *oauth2.Config {
//	return config
//}
//
//func (p *dropboxProvider) Login(token string) (acc *storage.Account, err error) {
//	if a, err := p.GetAccountInfo(token); err != nil {
//		return acc, errors.Wrap(err, 0)
//	} else {
//		id := strconv.FormatInt(int64(a.UID), 10)
//		if acc, err := p.store.FindAccountByProviderID("dropbox", id); err != nil {
//			return acc, provider.ErrUnauthorized
//		} else {
//			return acc, nil
//		}
//	}
//}
//
//func (p *dropboxProvider) ExchangeCodeForToken(code string) (*storage.AccountLink, error) {
//	if t, err := p.store.FindLinkByCode("dropbox", code); err != nil {
//		log.Printf("Info: Got error %s in ExchangeCodeForToken.", err.Error())
//	} else {
//		return t, nil
//	}
//
//	if err := p.db.AuthCode(code); err != nil {
//		return nil, errors.Wrap(err, 0)
//	} else {
//		token := p.db.AccessToken()
//		if a, err := p.GetAccountInfo(token); err != nil {
//			return nil, errors.Wrap(err, 0)
//		} else if link, err := p.store.CreateLink("dropbox", code, token, strconv.FormatInt(int64(a.UID), 10)); err != nil {
//			return nil, errors.Wrap(err, 0)
//		} else {
//			return link, nil
//		}
//	}
//}
//
//func (p *dropboxProvider) LinkAccount(account *storage.Account, token string) error {
//	if link, err := p.store.FindLinkByToken("dropbox", token); err != nil {
//		return errors.Wrap(err, 0)
//	} else if err := p.store.LinkAccount(account, link); err != nil {
//		return errors.Wrap(err, 0)
//	}
//	return nil
//}
//
//func (p *dropboxProvider) GetAccountInfo(token string) (*dropbox.Account, error) {
//	p.db.SetAccessToken(token)
//	if a, err := p.db.GetAccountInfo(); err != nil {
//		return a, errors.Wrap(err, 1)
//	} else {
//		return a, nil
//	}
//}
