package expr_test

import (
	"fmt"
	"github.com/antonmedv/expr"
	"github.com/antonmedv/expr/conf"
	"github.com/antonmedv/expr/vm"
	"reflect"
	"testing"
)

/*
go test -bench=BenchmarkEvalFunc1 -benchmem -cpuprofile profile1.out
go test -bench=BenchmarkEvalFunc2 -benchmem -cpuprofile profile2.out
go tool pprof -http=:9091 profile1.out
go tool pprof -http=:9092 profile2.out
*/

func TestEval(t *testing.T) {
	p, _ := expr.Compile("CalcWeight(20, 2, 100, 4.0, FollowRateDistinct)", func(c *conf.Config) {
		if c.Types == nil {
			c.Types = make(map[string]conf.Tag)
		}
		c.Types["CalcWeight"] = conf.Tag{Type: reflect.TypeOf(func(...interface{}) interface{} { return nil })}
	})

	v, _ := vm.Run(p, map[string]interface{}{
		"CalcWeight": func(...interface{}) interface{} {
			return 1.1
		},
		"FollowRateDistinct": 2.0,
	})
	fmt.Println(v)
}

func BenchmarkEvalFunc1(b *testing.B) {
	p, _ := expr.Compile("CalcWeight(20, 2, 100, 4.0, FollowRateDistinct)")
	m := map[string]interface{}{
		"CalcWeight": func(v1, v2, v3, v4, v0 interface{}) float64 {
			return 1.1
		},
		"FollowRateDistinct": 2.0,
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		vm.Run(p, m)
	}
}

func BenchmarkEvalFunc2(b *testing.B) {
	m := map[string]interface{}{
		"CalcWeight": func(...interface{}) interface{} {
			return 1.1
		},
		"FollowRateDistinct": 2.0,
	}

	p, _ := expr.Compile("CalcWeight(20, 2, 100, 4.0, FollowRateDistinct)", expr.Env(m))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		vm.Run(p, m)
	}
}
