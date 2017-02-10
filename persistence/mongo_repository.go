package persistence

import (
	fltr "github.com/go-scim/scimify/filter"
	"github.com/go-scim/scimify/resource"
	"gopkg.in/mgo.v2"
)

func NewMongoRepository(mongoAddress string) *MongoRepository {
	return &MongoRepository{
		address: mongoAddress,
	}
}

type MongoRepository struct {
	address string
}

func (m *MongoRepository) Create(resource *resource.Resource, context resource.Context) error {
	session := m.getSession()
	defer session.Close()

	return nil
}

func (m *MongoRepository) Get(id string, context resource.Context) (*resource.Resource, error) {
	session := m.getSession()
	defer session.Close()

	return nil, nil
}

func (m *MongoRepository) Replace(id string, resource *resource.Resource, context resource.Context) error {
	session := m.getSession()
	defer session.Close()

	return nil
}

func (m *MongoRepository) Delete(id string, context resource.Context) error {
	session := m.getSession()
	defer session.Close()

	return nil
}

func (m *MongoRepository) Query(filter string, sortBy string, ascending bool, pageStart int32, pageSize int32, context resource.Context) ([]*resource.Resource, error) {
	session := m.getSession()
	defer session.Close()

	tokens, err := fltr.Tokenize(filter)
	if err != nil {
		return err
	}

	root, err := fltr.Parse(tokens)
	if err != nil {
		return err
	}

	if root { // remove this line
		// TODO transform root to mongo query
	}

	return nil, nil
}

func (m *MongoRepository) getSession() *mgo.Session {
	session, err := mgo.Dial(m.address)
	if err != nil {
		panic(err)
	}
	return session
}
