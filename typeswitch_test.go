// Description:
// Author: liming.one@bytedance.com
package gruptor

import "testing"

// Golang type assertion / type switch will cost more time
type eventInter interface {
}

type event struct {
	data int64
}

func BenchmarkNoTypeSwitch(b *testing.B) {
	e := &event{}
	var holder []*event
	holder = append(holder, e)
	times := int64(b.N)
	b.ResetTimer()
	b.ReportAllocs()
	for i := int64(0); i < times; i++ {
		holder[0].data = i
	}
}

func BenchmarkUsingTypeSwitch(b *testing.B) {
	e := &event{}
	var holder []eventInter
	holder = append(holder, e)
	times := int64(b.N)
	b.ResetTimer()
	b.ReportAllocs()
	for i := int64(0); i < times; i++ {
		holder[0].(*event).data = i
	}
}

func BenchmarkUsingTypeAssertion(b *testing.B) {
	e := &event{}
	var holder []eventInter
	holder = append(holder, e)
	times := int64(b.N)
	b.ResetTimer()
	b.ReportAllocs()
	for i := int64(0); i < times; i++ {
		if ev, ok := holder[0].(*event); ok {
			ev.data = i
		}
	}
}
