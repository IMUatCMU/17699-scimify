package processor

import (
	"fmt"
	"github.com/go-scim/scimify/helper"
	"github.com/spf13/viper"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"testing"
)

func benchmarkDbQueryProcessor(poolSize int, b *testing.B) {
	viper.Set("mongo.address", "localhost:32768")
	viper.Set("mongo.database", "benchmark")
	viper.Set("mongo.userCollectionName", "users")
	viper.Set(fmt.Sprintf("scim.threadPool.%s", DbUserQuery), poolSize)

	r, _, err := helper.LoadResource("../test_data/single_test_user_david.json")
	if err != nil {
		b.Fatal(err)
	}

	getCollection := func() *mgo.Collection {
		session, err := mgo.Dial(viper.GetString("mongo.address"))
		if err != nil {
			b.Fatal(err)
		}

		return session.DB(viper.GetString("mongo.database")).C(viper.GetString("mongo.userCollectionName"))
	}

	cleanUp := func() {
		c := getCollection()
		c.RemoveAll(nil)
	}

	prepare := func() {
		c := getCollection()
		c.Insert(r.Data())
	}

	cleanUp()
	prepare()

	processor := GetWorkerBean(DbUserQuery)
	makeContext := func() *ProcessorContext {
		return &ProcessorContext{
			ParsedFilter: bson.M{
				"userName": bson.M{
					"$eq": "david",
				},
			},
			QueryPageStart: 0,
			QueryPageSize:  10,
		}
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			context := makeContext()
			err := processor.Process(context)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkDbQueryProcessorWithPoolSize1(b *testing.B) {
	benchmarkDbQueryProcessor(1, b)
}

func BenchmarkDbQueryProcessorWithPoolSize2(b *testing.B) {
	benchmarkDbQueryProcessor(2, b)
}

func BenchmarkDbQueryProcessorWithPoolSize3(b *testing.B) {
	benchmarkDbQueryProcessor(3, b)
}

func BenchmarkDbQueryProcessorWithPoolSize4(b *testing.B) {
	benchmarkDbQueryProcessor(4, b)
}

func BenchmarkDbQueryProcessorWithPoolSize5(b *testing.B) {
	benchmarkDbQueryProcessor(5, b)
}

func BenchmarkDbQueryProcessorWithPoolSize6(b *testing.B) {
	benchmarkDbQueryProcessor(6, b)
}

func BenchmarkDbQueryProcessorWithPoolSize7(b *testing.B) {
	benchmarkDbQueryProcessor(7, b)
}

func BenchmarkDbQueryProcessorWithPoolSize8(b *testing.B) {
	benchmarkDbQueryProcessor(8, b)
}

func BenchmarkDbQueryProcessorWithPoolSize9(b *testing.B) {
	benchmarkDbQueryProcessor(9, b)
}
