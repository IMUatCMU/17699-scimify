package processor

import (
	"fmt"
	"github.com/go-scim/scimify/helper"
	"github.com/spf13/viper"
	"gopkg.in/mgo.v2"
	"testing"
)

func benchmarkDbCreateProcessor(poolSize int, b *testing.B) {
	viper.Set("mongo.address", "localhost:32768")
	viper.Set("mongo.database", "benchmark")
	viper.Set("mongo.userCollectionName", "users")
	viper.Set(fmt.Sprintf("scim.threadPool.%s", DbUserCreate), poolSize)

	cleanUp := func() {
		session, err := mgo.Dial(viper.GetString("mongo.address"))
		if err != nil {
			b.Fatal(err)
		}

		c := session.DB(viper.GetString("mongo.database")).C(viper.GetString("mongo.userCollectionName"))
		c.RemoveAll(nil)
	}

	r, _, err := helper.LoadResource("../test_data/single_test_user_david.json")
	if err != nil {
		b.Fatal(err)
	}

	processor := GetWorkerBean(DbUserCreate)
	makeContext := func() *ProcessorContext {
		return &ProcessorContext{Resource: r}
	}

	cleanUp()
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

func BenchmarkDbCreateProcessorWithPoolSize1(b *testing.B) {
	benchmarkDbCreateProcessor(1, b)
}

func BenchmarkDbCreateProcessorWithPoolSize2(b *testing.B) {
	benchmarkDbCreateProcessor(2, b)
}

func BenchmarkDbCreateProcessorWithPoolSize3(b *testing.B) {
	benchmarkDbCreateProcessor(3, b)
}

func BenchmarkDbCreateProcessorWithPoolSize4(b *testing.B) {
	benchmarkDbCreateProcessor(4, b)
}

func BenchmarkDbCreateProcessorWithPoolSize5(b *testing.B) {
	benchmarkDbCreateProcessor(5, b)
}

func BenchmarkDbCreateProcessorWithPoolSize6(b *testing.B) {
	benchmarkDbCreateProcessor(6, b)
}

func BenchmarkDbCreateProcessorWithPoolSize7(b *testing.B) {
	benchmarkDbCreateProcessor(7, b)
}

func BenchmarkDbCreateProcessorWithPoolSize8(b *testing.B) {
	benchmarkDbCreateProcessor(8, b)
}

func BenchmarkDbCreateProcessorWithPoolSize9(b *testing.B) {
	benchmarkDbCreateProcessor(9, b)
}
