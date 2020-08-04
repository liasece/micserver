// Copyright (c) 2016 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package observer_test

import (
	"testing"

	"github.com/liasece/micserver/log"
	"github.com/liasece/micserver/log/core"
	. "github.com/liasece/micserver/log/test/observer"
	"github.com/stretchr/testify/assert"
)

func assertEmpty(t testing.TB, logs *ObservedLogs) {
	assert.Equal(t, 0, logs.Len(), "Expected empty ObservedLogs to have zero length.")
	assert.Equal(t, []LoggedEntry{}, logs.All(), "Unexpected LoggedEntries in empty ObservedLogs.")
}

func TestObserver(t *testing.T) {
	observer, logs := New(log.InfoLevel)
	assertEmpty(t, logs)

	assert.NoError(t, observer.Sync(), "Unexpected failure in no-op Sync")

	obs := log.NewLogger(observer).With(log.Int("i", 1))
	obs.Info("foo")
	obs.Debug("bar")
	// want := []LoggedEntry{{
	// 	Entry:   core.Entry{Level: log.InfoLevel, Message: "foo"},
	// 	Context: []core.Field{log.Int("i", 1)},
	// }}

	// assert.Equal(t, 1, logs.Len(), "Unexpected observed logs Len.")
	// assert.Equal(t, want, logs.AllUntimed(), "Unexpected contents from AllUntimed.")

	// all := logs.All()
	// require.Equal(t, 1, len(all), "Unexpected numbed of LoggedEntries returned from All.")
	// assert.NotEqual(t, time.Time{}, all[0].Time, "Expected non-zero time on LoggedEntry.")

	// copy & zero time for stable assertions
	// untimed := append([]LoggedEntry{}, all...)
	// untimed[0].Time = time.Time{}
	// assert.Equal(t, want, untimed, "Unexpected LoggedEntries from All.")

	// assert.Equal(t, all, logs.TakeAll(), "Expected All and TakeAll to return identical results.")
	assertEmpty(t, logs)
}

func TestObserverWith(t *testing.T) {
	sf1, logs := New(log.InfoLevel)

	// need to pad out enough initial fields so that the underlying slice cap()
	// gets ahead of its len() so that the sf3/4 With append's could choose
	// not to copy (if the implementation doesn't force them)
	sf1 = sf1.With([]core.Field{log.Int("a", 1), log.Int("b", 2)})

	sf2 := sf1.With([]core.Field{log.Int("c", 3)})
	sf3 := sf2.With([]core.Field{log.Int("d", 4)})
	sf4 := sf2.With([]core.Field{log.Int("e", 5)})
	ent := core.Entry{Level: log.InfoLevel, Message: "hello"}

	for i, core := range []core.Core{sf2, sf3, sf4} {
		if ce := core.Check(ent, nil); ce != nil {
			ce.Write(log.Int("i", i))
		}
	}

	assert.Equal(t, []LoggedEntry{
		{
			Entry: ent,
			Context: []core.Field{
				log.Int("a", 1),
				log.Int("b", 2),
				log.Int("c", 3),
				log.Int("i", 0),
			},
		},
		{
			Entry: ent,
			Context: []core.Field{
				log.Int("a", 1),
				log.Int("b", 2),
				log.Int("c", 3),
				log.Int("d", 4),
				log.Int("i", 1),
			},
		},
		{
			Entry: ent,
			Context: []core.Field{
				log.Int("a", 1),
				log.Int("b", 2),
				log.Int("c", 3),
				log.Int("e", 5),
				log.Int("i", 2),
			},
		},
	}, logs.All(), "expected no field sharing between With siblings")
}

func TestFilters(t *testing.T) {
	logs := []LoggedEntry{
		{
			Entry:   core.Entry{Level: log.InfoLevel, Message: "log a"},
			Context: []core.Field{log.String("fStr", "1"), log.Int("a", 1)},
		},
		{
			Entry:   core.Entry{Level: log.InfoLevel, Message: "log a"},
			Context: []core.Field{log.String("fStr", "2"), log.Int("b", 2)},
		},
		{
			Entry:   core.Entry{Level: log.InfoLevel, Message: "log b"},
			Context: []core.Field{log.Int("a", 1), log.Int("b", 2)},
		},
		{
			Entry:   core.Entry{Level: log.InfoLevel, Message: "log c"},
			Context: []core.Field{log.Int("a", 1), log.Namespace("ns"), log.Int("a", 2)},
		},
		{
			Entry:   core.Entry{Level: log.InfoLevel, Message: "msg 1"},
			Context: []core.Field{log.Int("a", 1), log.Namespace("ns")},
		},
		{
			Entry:   core.Entry{Level: log.InfoLevel, Message: "any map"},
			Context: []core.Field{log.Any("map", map[string]string{"a": "b"})},
		},
		{
			Entry:   core.Entry{Level: log.InfoLevel, Message: "any slice"},
			Context: []core.Field{log.Any("slice", []string{"a"})},
		},
	}

	logger, sink := New(log.InfoLevel)
	for _, log := range logs {
		logger.Write(log.Entry, log.Context)
	}

	tests := []struct {
		msg      string
		filtered *ObservedLogs
		want     []LoggedEntry
	}{
		{
			msg:      "filter by message",
			filtered: sink.FilterMessage("log a"),
			want:     logs[0:2],
		},
		{
			msg:      "filter by field",
			filtered: sink.FilterField(log.String("fStr", "1")),
			want:     logs[0:1],
		},
		{
			msg:      "filter by message and field",
			filtered: sink.FilterMessage("log a").FilterField(log.Int("b", 2)),
			want:     logs[1:2],
		},
		{
			msg:      "filter by field with duplicate fields",
			filtered: sink.FilterField(log.Int("a", 2)),
			want:     logs[3:4],
		},
		{
			msg:      "filter doesn't match any messages",
			filtered: sink.FilterMessage("no match"),
			want:     []LoggedEntry{},
		},
		{
			msg:      "filter by snippet",
			filtered: sink.FilterMessageSnippet("log"),
			want:     logs[0:4],
		},
		{
			msg:      "filter by snippet and field",
			filtered: sink.FilterMessageSnippet("a").FilterField(log.Int("b", 2)),
			want:     logs[1:2],
		},
		{
			msg:      "filter for map",
			filtered: sink.FilterField(log.Any("map", map[string]string{"a": "b"})),
			want:     logs[5:6],
		},
		{
			msg:      "filter for slice",
			filtered: sink.FilterField(log.Any("slice", []string{"a"})),
			want:     logs[6:7],
		},
	}

	for _, tt := range tests {
		got := tt.filtered.AllUntimed()
		assert.Equal(t, tt.want, got, tt.msg)
	}
}
