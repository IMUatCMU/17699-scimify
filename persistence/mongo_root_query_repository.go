package persistence

import (
	"github.com/go-scim/scimify/resource"
	"gopkg.in/mgo.v2"
)

type MongoRootQueryRepository struct {
	address         string
	databaseName    string
	collectionNames []string
}

func (repo *MongoRootQueryRepository) Create(resource.ScimObject) error {
	panic("not supported")
}

func (repo *MongoRootQueryRepository) GetAll() ([]resource.ScimObject, error) {
	panic("not supported")
}

func (repo *MongoRootQueryRepository) Get(string) (resource.ScimObject, error) {
	panic("not supported")
}

func (repo *MongoRootQueryRepository) Replace(string, resource.ScimObject) error {
	panic("not supported")
}

func (repo *MongoRootQueryRepository) Delete(string) error {
	panic("not supported")
}

func (repo *MongoRootQueryRepository) Query(filter interface{}, sortBy string, ascending bool, pageStart int, pageSize int) ([]resource.ScimObject, error) {
	// get session
	session := repo.getSession()
	defer session.Close()

	// prepare plans
	skipQuota, limitQuota := 0, 0

	if pageStart < 1 {
		skipQuota = 0
	} else {
		skipQuota = pageStart - 1
	}

	if pageSize > 0 {
		limitQuota = pageSize
	} else {
		limitQuota = 0
	}

	plans := repo.generateExecutionPlans()
	for _, plan := range plans {
		count, err := repo.getCollection(session, plan.collectionName).Find(filter).Count()
		if err != nil {
			plan.skipAll = true
		} else {
			plan.resultCount = count
		}
	}

	// devise plans
	for _, plan := range plans {
		if plan.skipAll {
			continue
		}

		if limitQuota <= 0 {
			plan.skipAll = true
			continue
		}

		if plan.resultCount <= 0 {
			plan.skipAll = true
			continue
		}

		if skipQuota <= 0 {
			// take a look at limit
			plan.skipAll = false
			plan.skip = 0
			if limitQuota < plan.resultCount {
				plan.limit = limitQuota
			} else {
				plan.limit = plan.resultCount
			}
			limitQuota -= plan.limit
		} else {
			// continue to skip
			if plan.resultCount > skipQuota {
				// partial results needed
				plan.skipAll = false
				plan.skip = skipQuota
				if limitQuota < plan.resultCount-plan.skip {
					plan.limit = limitQuota
				} else {
					plan.limit = plan.resultCount - plan.skip
				}
				limitQuota -= plan.limit
			} else {
				plan.skipAll = true
			}
			skipQuota -= plan.resultCount
		}
	}

	// execution plan
	allData := make([]map[string]interface{}, 0)
	for _, plan := range plans {
		if plan.skipAll || plan.limit == 0 {
			continue
		}

		query := repo.getCollection(session, plan.collectionName).Find(filter)
		if len(sortBy) > 0 {
			if ascending {
				query = query.Sort(sortBy)
			} else {
				query = query.Sort("-" + sortBy)
			}
		}
		query = query.Skip(plan.skip)
		query = query.Limit(plan.limit)

		data := make([]map[string]interface{}, 0)
		query.Iter().All(&data)
		allData = append(allData, data...)
	}

	resources := make([]resource.ScimObject, 0, len(allData))
	for _, data := range allData {
		resources = append(resources, resource.NewResourceFromMap(data))
	}

	return resources, nil
}

func (repo *MongoRootQueryRepository) getSession() *mgo.Session {
	session, err := mgo.Dial(repo.address)
	if err != nil {
		panic(err)
	}
	return session
}

func (repo *MongoRootQueryRepository) getCollection(session *mgo.Session, name string) *mgo.Collection {
	return session.DB(repo.databaseName).C(name)
}

func (repo *MongoRootQueryRepository) generateExecutionPlans() []*queryExecutionPlan {
	plans := make([]*queryExecutionPlan, 0, len(repo.collectionNames))
	for _, n := range repo.collectionNames {
		plans = append(plans, &queryExecutionPlan{
			collectionName: n,
		})
	}
	return plans
}

type queryExecutionPlan struct {
	collectionName string
	resultCount    int
	skipAll        bool
	skip           int
	limit          int
}
