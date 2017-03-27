package processor

import (
	"fmt"
	"github.com/go-scim/scimify/helper"
	"github.com/spf13/viper"
	"gopkg.in/mgo.v2"
	"testing"
)

func benchmarkDbGetProcessor(poolSize int, b *testing.B) {
	viper.Set("mongo.address", "localhost:32768")
	viper.Set("mongo.database", "benchmark")
	viper.Set("mongo.userCollectionName", "users")
	viper.Set(fmt.Sprintf("scim.threadPool.%s", DbUserGetToResource), poolSize)

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

	processor := GetWorkerBean(DbUserGetToResource)
	makeContext := func() *ProcessorContext {
		return &ProcessorContext{Identity: r.GetId()}
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

func BenchmarkDbGetProcessorWithPoolSize1(b *testing.B) {
	benchmarkDbGetProcessor(1, b)
}

func BenchmarkDbGetProcessorWithPoolSize2(b *testing.B) {
	benchmarkDbGetProcessor(2, b)
}

func BenchmarkDbGetProcessorWithPoolSize3(b *testing.B) {
	benchmarkDbGetProcessor(3, b)
}

func BenchmarkDbGetProcessorWithPoolSize4(b *testing.B) {
	benchmarkDbGetProcessor(4, b)
}

func BenchmarkDbGetProcessorWithPoolSize5(b *testing.B) {
	benchmarkDbGetProcessor(5, b)
}

func BenchmarkDbGetProcessorWithPoolSize6(b *testing.B) {
	benchmarkDbGetProcessor(6, b)
}

func BenchmarkDbGetProcessorWithPoolSize7(b *testing.B) {
	benchmarkDbGetProcessor(7, b)
}

func BenchmarkDbGetProcessorWithPoolSize8(b *testing.B) {
	benchmarkDbGetProcessor(8, b)
}

func BenchmarkDbGetProcessorWithPoolSize9(b *testing.B) {
	benchmarkDbGetProcessor(9, b)
}
