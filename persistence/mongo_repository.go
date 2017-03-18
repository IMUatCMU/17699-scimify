package persistence

import (
	"github.com/go-scim/scimify/resource"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func NewMongoRepository(mongoAddress, database, collection string) *MongoRepository {
	return &MongoRepository{
		address:        mongoAddress,
		databaseName:   database,
		collectionName: collection,
	}
}

type MongoRepository struct {
	address        string
	databaseName   string
	collectionName string
}

func (m *MongoRepository) Create(resource resource.ScimObject) error {
	session := m.getSession()
	defer session.Close()

	return m.getCollection(session).Insert(resource.Data())
}

func (m *MongoRepository) GetAll() ([]resource.ScimObject, error) {
	return nil, resource.CreateError(resource.NotImplemented, "get all is not implemented for monogo repository.")
}

func (m *MongoRepository) Get(id string) (resource.ScimObject, error) {
	session := m.getSession()
	defer session.Close()

	query := m.getCollection(session).Find(bson.M{
		"id": id,
	})

	data := make(map[string]interface{}, 0)
	err := query.One(&data)
	if err != nil {
		return nil, err
	}

	return resource.NewResourceFromMap(data), nil
}

func (m *MongoRepository) Replace(id string, resource resource.ScimObject) error {
	session := m.getSession()
	defer session.Close()

	criteria := bson.M{"id": resource.GetId()}
	update := bson.M{"$set": resource.Data()}

	return m.getCollection(session).Update(criteria, update)
}

func (m *MongoRepository) Delete(id string) error {
	session := m.getSession()
	defer session.Close()

	return m.getCollection(session).Remove(bson.M{
		"id": id,
	})
}

// Query mongoDB for entries.
// - filter: mongo styled filters in mgo.M
// - sortBy: empty or a valid resource full path
// - ascending: sort order, ignored when sortBy is empty
// - pageStart: skip how many entries, if less than 0, will be defaulted to 0
// - pageSize: collect how many entries, if less than 0, will be ignored
func (m *MongoRepository) Query(filter interface{}, sortBy string, ascending bool, pageStart int, pageSize int) ([]resource.ScimObject, error) {
	// get session
	session := m.getSession()
	defer session.Close()

	// prepare query
	query := m.getCollection(session).Find(filter)

	// sort order
	if len(sortBy) > 0 {
		if ascending {
			query = query.Sort(sortBy)
		} else {
			query = query.Sort("-" + sortBy)
		}
	}

	// page start
	if pageStart < 1 {
		query = query.Skip(0)
	} else {
		query = query.Skip(pageStart - 1)
	}

	// page size
	if pageSize > 0 {
		query = query.Limit(pageSize)
	}

	// execute query
	rawData := make([]map[string]interface{}, 0)
	query.Iter().All(&rawData)

	// parse data
	resources := make([]resource.ScimObject, 0, len(rawData))
	for _, data := range rawData {
		//resources = append(resources, parseResource(data))
		resources = append(resources, resource.NewResourceFromMap(data))
	}

	return resources, nil
}

func (m *MongoRepository) getSession() *mgo.Session {
	session, err := mgo.Dial(m.address)
	if err != nil {
		panic(err)
	}
	return session
}

func (m *MongoRepository) getCollection(session *mgo.Session) *mgo.Collection {
	return session.DB(m.databaseName).C(m.collectionName)
}
