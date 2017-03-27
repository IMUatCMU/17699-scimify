package processor

import (
	"fmt"
	"github.com/go-scim/scimify/helper"
	"github.com/spf13/viper"
	"gopkg.in/mgo.v2"
	"testing"
)

func benchmarkDbDeleteProcessor(poolSize int, b *testing.B) {
	viper.Set("mongo.address", "localhost:32768")
	viper.Set("mongo.database", "benchmark")
	viper.Set("mongo.userCollectionName", "users")
	viper.Set(fmt.Sprintf("scim.threadPool.%s", DbUserDelete), poolSize)

	r, _, err := helper.LoadResource("../test_data/single_test_user_david.json")
	if err != nil {
		b.Fatal(err)
	}
	identity := r.GetId()

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
		for i := 0; i < b.N; i++ {
			c.Insert(r.Data())
		}
	}

	cleanUp()
	prepare()

	processor := GetWorkerBean(DbUserDelete)
	makeContext := func() *ProcessorContext {
		return &ProcessorContext{Identity: identity}
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

func BenchmarkDbDeleteProcessorWithPoolSize1(b *testing.B) {
	benchmarkDbDeleteProcessor(1, b)
}

func BenchmarkDbDeleteProcessorWithPoolSize2(b *testing.B) {
	benchmarkDbDeleteProcessor(2, b)
}

func BenchmarkDbDeleteProcessorWithPoolSize3(b *testing.B) {
	benchmarkDbDeleteProcessor(3, b)
}

func BenchmarkDbDeleteProcessorWithPoolSize4(b *testing.B) {
	benchmarkDbDeleteProcessor(4, b)
}

func BenchmarkDbDeleteProcessorWithPoolSize5(b *testing.B) {
	benchmarkDbDeleteProcessor(5, b)
}

func BenchmarkDbDeleteProcessorWithPoolSize6(b *testing.B) {
	benchmarkDbDeleteProcessor(6, b)
}

func BenchmarkDbDeleteProcessorWithPoolSize7(b *testing.B) {
	benchmarkDbDeleteProcessor(7, b)
}

func BenchmarkDbDeleteProcessorWithPoolSize8(b *testing.B) {
	benchmarkDbDeleteProcessor(8, b)
}

func BenchmarkDbDeleteProcessorWithPoolSize9(b *testing.B) {
	benchmarkDbDeleteProcessor(9, b)
}
