package dropbox

//
//func (s *Storage) CreateLink(provider, code, token, providerID string) (*AccountLink, error) {
//	a := &AccountLink{
//		ID:          uuid.NewRandom().String(),
//		Provider:    provider,
//		AuthCode:    code,
//		AccessToken: token,
//		ProviderID:  providerID,
//	}
//	if err := s.getCollection(accountLinkCollection).Insert(a); err != nil {
//		return a, errors.Wrap(err, 0)
//	}
//	return a, nil
//}
//
//func (s *Storage) FindLinkByCode(provider, code string) (*AccountLink, error) {
//	var e AccountLink
//	c := s.getCollection(accountLinkCollection)
//	if err := c.Find(bson.M{"authcode": code, "provider": provider}).One(&e); err != nil {
//		return &e, errors.Wrap(err, 0)
//	}
//	return &e, nil
//}
//
//func (s *Storage) FindLinkByToken(provider, token string) (*AccountLink, error) {
//	e := new(AccountLink)
//	c := s.getCollection(accountLinkCollection)
//	if err := c.Find(bson.M{"accesstoken": token, "provider": provider}).One(e); err != nil {
//		return e, errors.Errorf("Could not find token %s in %s: %s", token, provider, err.Error())
//	}
//	return e, nil
//}
//
//func (s *Storage) CreateAccount(email, username string) (*Account, error) {
//	a := &Account{
//		ID:           uuid.NewRandom().String(),
//		Email:        email,
//		Username:     username,
//		RegisteredAt: time.Now(),
//	}
//	if err := s.getCollection(accountCollection).Insert(a); err != nil {
//		return a, errors.Wrap(err, 0)
//	}
//	return a, nil
//}
//
//func (s *Storage) LinkAccount(account *Account, link *AccountLink) error {
//	if err := s.getCollection(accountCollection).Update(bson.M{"id": account.ID}, bson.M{"$push": bson.M{"links": link.ID}}); err != nil {
//		return errors.Wrap(err, 0)
//	} else if err := s.getCollection(accountLinkCollection).Update(bson.M{"id": link.ID}, bson.M{"$set": bson.M{"accountid": account.ID}}); err != nil {
//		return errors.Wrap(err, 0)
//	}
//	return nil
//}
//
//func (s *Storage) FindAccountByProviderID(provider, id string) (*Account, error) {
//	l := new(AccountLink)
//	e := new(Account)
//	log.Printf("Trying to fetch %s %s", provider, id)
//	if err := s.getCollection(accountLinkCollection).Find(bson.M{"providerid": id, "provider": provider, "accountid": bson.M{"$ne": ""}}).One(l); err != nil {
//		log.Printf("No account at %s with provider id %s found", provider, id)
//		return e, errors.Wrap(err, 0)
//	} else if err := s.getCollection(accountCollection).Find(bson.M{"id": l.AccountID}).One(e); err != nil {
//		log.Printf("No account with id %s found. %s", l.AccountID, l.ID)
//		return e, errors.Wrap(err, 0)
//	} else {
//		log.Printf("Account for provider %s with id %s found", provider, id)
//		return e, nil
//	}
//}
//
//func (s *Storage) FindSharesByAccount(account *Account) (shares []*Share, err error) {
//	c := s.getCollection(shareCollection)
//	if err := c.Find(bson.M{"account": account.ID}).All(&shares); err != nil {
//		return shares, errors.Wrap(err, 0)
//	}
//	return shares, nil
//}
//
//func (s *Storage) GetAccount(id string) (*Account, error) {
//	var e Account
//	c := s.getCollection(accountCollection)
//	if err := c.Find(bson.M{"id": id}).One(&e); err != nil {
//		return &e, errors.Wrap(err, 0)
//	}
//	return &e, nil
//}
//
//func (s *Storage) GetShare(id string) (*Share, error) {
//	var e Share
//	c := s.getCollection(shareCollection)
//	if err := c.Find(bson.M{"id": id}).One(&e); err != nil {
//		return &e, errors.Wrap(err, 0)
//	}
//	return &e, nil
//}
//
//func (s *Storage) FindShareBySlug(slug string) (*Share, error) {
//	var e Share
//	c := s.getCollection(shareCollection)
//	if err := c.Find(bson.M{"slug": slug}).One(&e); err != nil {
//		return &e, errors.Wrap(err, 0)
//	}
//	return &e, nil
//}
//
//func (s *Storage) CreateShare(account, destination, title, name, email, parent, thumbnail string, expiresAt time.Time) (*Share, error) {
//	slug, err := sequence.RuneSequence(10, runes)
//	if err != nil {
//		return new(Share), errors.Wrap(err, 0)
//	}
//
//	var filename string
//	if d, err := url.Parse(destination); err != nil {
//		return new(Share), errors.Wrap(err, 0)
//	} else {
//		t := strings.Split(d.Path, "/")
//		filename = t[len(t)-1]
//	}
//
//	share := &Share{
//		ID:          uuid.NewRandom().String(),
//		Account:     account,
//		Destination: destination,
//		Name:        name,
//		Title:       title,
//		Email:       email,
//		Timestamp:   time.Now(),
//		ExpiresAt:   expiresAt,
//		Slug:        string(slug),
//		Parent:      parent,
//		Thumbnail:   thumbnail,
//		Filename:    filename,
//	}
//	if parent != "" {
//		if err := s.getCollection(shareCollection).Update(bson.M{"id": parent}, bson.M{"$push": bson.M{"children": share.ID}}); err != nil {
//			return nil, err
//		}
//	}
//	if err := s.getCollection(shareCollection).Insert(share); err != nil {
//		return nil, errors.Wrap(err, 0)
//	}
//
//	return share, nil
//}
//
//func (s *Storage) GetValidShare(slug string) (*Share, error) {
//	share, err := s.FindShareBySlug(slug)
//	if err != nil {
//		return share, errors.Wrap(err, 0)
//	}
//	if share.ExpiresAt.Before(time.Now()) {
//		return new(Share), ErrExpiredShare
//	}
//	return share, nil
//}
//
//func New(session *mgo.Session) *Storage {
//	s := &Storage{
//		session:  session,
//	}
//	s.ensureUnique(shareCollection, []string{"slug"})
//	s.ensureUnique(accountCollection, []string{"email"})
//	for _, k := range collections {
//		s.ensureUnique(k, []string{"id"})
//	}
//	return s
//}
//
//func (s *Storage) getCollection(name string) *mgo.Collection {
//	if s.sessionCopy != nil {
//		s.sessionCopy.Close()
//	}
//	s.sessionCopy = s.session.Copy()
//	return s.sessionCopy.DB("").C(name)
//}
//
//func (s *Storage) ensureUnique(collection string, keys []string) {
//	c := s.getCollection(collection)
//	if err := c.EnsureIndex(mgo.Index{Key: keys, Unique: true}); err != nil {
//		log.Fatalf("Could not ensure index: %s", err)
//	}
//}
