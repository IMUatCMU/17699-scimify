package persistence

import (
	"context"
	"github.com/go-scim/scimify/resource"
	"gopkg.in/mgo.v2"
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

func (m *MongoRepository) Create(resource *resource.Resource, context context.Context) error {
	session := m.getSession()
	defer session.Close()

	return m.getCollection(session).Insert(resource.ToMap())
}

func (m *MongoRepository) Get(id string, context context.Context) (*resource.Resource, error) {
	session := m.getSession()
	defer session.Close()

	return nil, nil
}

func (m *MongoRepository) Replace(id string, resource *resource.Resource, context context.Context) error {
	session := m.getSession()
	defer session.Close()

	return nil
}

func (m *MongoRepository) Delete(id string, context context.Context) error {
	session := m.getSession()
	defer session.Close()

	return nil
}

// Query mongoDB for entries.
// - filter: mongo styled filters in mgo.M
// - sortBy: empty or a valid resource full path
// - ascending: sort order, ignored when sortBy is empty
// - pageStart: skip how many entries, if less than 0, will be defaulted to 0
// - pageSize: collect how many entries, if less than 0, will be ignored
// - context: auxiliary information for the query
func (m *MongoRepository) Query(filter interface{}, sortBy string, ascending bool, pageStart int, pageSize int, context context.Context) ([]*resource.Resource, error) {
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
	if pageStart < 0 {
		query = query.Skip(0)
	} else {
		query = query.Skip(pageStart)
	}

	// page size
	if pageSize > 0 {
		query = query.Limit(pageSize)
	}

	// execute query
	rawData := make([]map[string]interface{}, 0)
	query.Iter().All(&rawData)

	// parse data
	resources := make([]*resource.Resource, 0, len(rawData))
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

//func parseResource(data map[string]interface{}) *resource.Resource {
//	resource := resource.NewResource()
//
//	if schemas, ok := data["schemas"].([]string); ok {
//		resource.Schemas = schemas
//		delete(data, "schemas")
//	}
//
//	if id, ok := data["id"].(string); ok {
//		resource.Id = id
//		delete(data, "id")
//	}
//
//	if externalId, ok := data["externalId"].(string); ok {
//		resource.ExternalId = externalId
//		delete(data, "externalId")
//	}
//
//	if meta, ok := data["meta"].(map[string]interface{}); ok {
//		if metaResourceType, ok := meta["resourceType"].(string); ok {
//			resource.Meta.ResourceType = metaResourceType
//		}
//		if metaCreated, ok := meta["created"].(string); ok {
//			resource.Meta.Created = metaCreated
//		}
//		if metaLastModified, ok := meta["lastModified"].(string); ok {
//			resource.Meta.LastModified = metaLastModified
//		}
//		if metaLocation, ok := meta["location"].(string); ok {
//			resource.Meta.Location = metaLocation
//		}
//		if metaVersion, ok := meta["version"].(string); ok {
//			resource.Meta.Version = metaVersion
//		}
//		delete(data, "meta")
//	}
//
//	resource.Attributes = data
//	return resource
//}
