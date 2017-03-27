package processor

import (
	"fmt"
	"github.com/go-scim/scimify/helper"
	"github.com/spf13/viper"
	"gopkg.in/mgo.v2"
	"testing"
)

func benchmarkDbReplaceProcessor(poolSize int, b *testing.B) {
	viper.Set("mongo.address", "localhost:32768")
	viper.Set("mongo.database", "benchmark")
	viper.Set("mongo.userCollectionName", "users")
	viper.Set(fmt.Sprintf("scim.threadPool.%s", DbUserReplace), poolSize)

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

	processor := GetWorkerBean(DbUserReplace)
	makeContext := func() *ProcessorContext {
		return &ProcessorContext{Identity: r.GetId(), Resource: r}
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

func BenchmarkDbReplaceProcessorWithPoolSize1(b *testing.B) {
	benchmarkDbReplaceProcessor(1, b)
}

func BenchmarkDbReplaceProcessorWithPoolSize2(b *testing.B) {
	benchmarkDbReplaceProcessor(2, b)
}

func BenchmarkDbReplaceProcessorWithPoolSize3(b *testing.B) {
	benchmarkDbReplaceProcessor(3, b)
}

func BenchmarkDbReplaceProcessorWithPoolSize4(b *testing.B) {
	benchmarkDbReplaceProcessor(4, b)
}

func BenchmarkDbReplaceProcessorWithPoolSize5(b *testing.B) {
	benchmarkDbReplaceProcessor(5, b)
}

func BenchmarkDbReplaceProcessorWithPoolSize6(b *testing.B) {
	benchmarkDbReplaceProcessor(6, b)
}

func BenchmarkDbReplaceProcessorWithPoolSize7(b *testing.B) {
	benchmarkDbReplaceProcessor(7, b)
}

func BenchmarkDbReplaceProcessorWithPoolSize8(b *testing.B) {
	benchmarkDbReplaceProcessor(8, b)
}

func BenchmarkDbReplaceProcessorWithPoolSize9(b *testing.B) {
	benchmarkDbReplaceProcessor(9, b)
}
