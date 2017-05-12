package urlshortner

import (
	"errors"

	couchdb "github.com/gospackler/caddyshack-couchdb"
	"github.com/gospackler/caddyshack/resource"
)

type Factory struct {
	tinyUrlStore *couchdb.CouchStore
}

func NewFactory(host string, port int, db string, desDoc string) *Factory {
	res := &resource.Definition{
		Host:   host,
		Port:   port,
		Name:   db,
		DesDoc: desDoc,
	}

	couchStore := couchdb.NewCouchStore(res, &TinyURL{})
	return &Factory{
		tinyUrlStore: couchStore,
	}

}

func (f *Factory) GetURL(tinyUrl string) (string, error) {
	objs, err := f.tinyUrlStore.ReadByKey(tinyUrl)
	if err != nil {
		return "", err
	}

	if len(objs) != 1 {
		return "", errors.New("More than one tinyUrl for key")
	}

	t := objs[0].(*TinyURL)
	t.Accessed = t.Accessed + 1
	err = f.tinyUrlStore.UpdateOne(t)
	if err != nil {
		return "", err
	}
	return t.LongURL, nil
}

func (f *Factory) AddURL(tinyURL, longURL string) error {
	t := new(TinyURL)
	t.TinyURL = tinyURL
	t.LongURL = longURL
	return f.tinyUrlStore.Create(t)
}

type TinyURL struct {
	TinyURL  string `json:"tinyurl" by:"tinyurl"`
	LongURL  string `json:"longurl"`
	Accessed int    `json:"accessed"`
	Id       string `json:"id"`
}

func (t *TinyURL) GetKey() string {
	return t.Id
}

func (t *TinyURL) SetKey(id string) {
	t.Id = id
}
