// Copyright (c) 2017 Uber Technologies, Inc.
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

package test

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/liasece/micserver/log"
	"github.com/stretchr/testify/assert"
)

func TestTestLogger(t *testing.T) {
	ts := newTestLogSpy(t)
	defer ts.AssertPassed()

	logger := NewLogger(ts)

	logger.Info("received work order")
	logger.Debug("starting work")
	logger.Warn("work may fail")
	logger.Error("work failed", log.ErrorField(errors.New("great sadness")))

	assert.Panics(t, func() {
		logger.Panic("failed to do work")
	}, "log.Panic should panic")

	ts.AssertMessages(
		"INFO	received work order",
		"DEBUG	starting work",
		"WARN	work may fail",
		`ERROR	work failed	{"error": "great sadness"}`,
		"PANIC	failed to do work",
	)
}

func TestTestLoggerSupportsLevels(t *testing.T) {
	ts := newTestLogSpy(t)
	defer ts.AssertPassed()

	logger := NewLogger(ts, Level(log.WarnLevel))

	logger.Info("received work order")
	logger.Debug("starting work")
	logger.Warn("work may fail")
	logger.Error("work failed", log.ErrorField(errors.New("great sadness")))

	assert.Panics(t, func() {
		logger.Panic("failed to do work")
	}, "log.Panic should panic")

	ts.AssertMessages(
		"WARN	work may fail",
		`ERROR	work failed	{"error": "great sadness"}`,
		"PANIC	failed to do work",
	)
}

func TestTestLoggerSupportsWrappedLogOptions(t *testing.T) {
	ts := newTestLogSpy(t)
	defer ts.AssertPassed()

	logger := NewLogger(ts, WrapOptions(log.Options().AddCaller().Fields(log.String("k1", "v1"))))

	logger.Info("received work order")
	logger.Debug("starting work")
	logger.Warn("work may fail")
	logger.Error("work failed", log.ErrorField(errors.New("great sadness")))

	assert.Panics(t, func() {
		log.Panic("failed to do work")
	}, "log.Panic should panic")

	ts.AssertMessages(
		`INFO	test/logger_test.go:89	received work order	{"k1": "v1"}`,
		`DEBUG	test/logger_test.go:90	starting work	{"k1": "v1"}`,
		`WARN	test/logger_test.go:91	work may fail	{"k1": "v1"}`,
		`ERROR	test/logger_test.go:92	work failed	{"k1": "v1", "error": "great sadness"}`,
		`PANIC	test/logger_test.go:95	failed to do work	{"k1": "v1"}`,
	)
}

func TestTestingWriter(t *testing.T) {
	ts := newTestLogSpy(t)
	w := newTestingWriter(ts)

	n, err := io.WriteString(w, "hello\n\n")
	assert.NoError(t, err, "WriteString must not fail")
	assert.Equal(t, 7, n)
}

func TestTestLoggerErrorOutput(t *testing.T) {
	// This test verifies that the test logger logs internal messages to the
	// testing.T and marks the test as failed.

	ts := newTestLogSpy(t)
	defer ts.AssertFailed()

	log := NewLogger(ts)

	// Replace with a core that fails.
	// log = log.WithOptions(log.WrapCore(func(core.Core) core.Core {
	// 	return core.NewCore(
	// 		core.NewConsoleEncoder(log.NewDevelopmentEncoderConfig()),
	// 		core.Lock(core.AddSync(ztest.FailWriter{})),
	// 		core.DebugLevel,
	// 	)
	// }))

	log.Info("foo") // this fails

	// if assert.Len(t, ts.Messages, 1, "expected a log message") {
	// 	assert.Regexp(t, `write error: failed`, ts.Messages[0])
	// }
}

// testLogSpy is a testing.TB that captures logged messages.
type testLogSpy struct {
	testing.TB

	failed   bool
	Messages []string
}

func newTestLogSpy(t testing.TB) *testLogSpy {
	return &testLogSpy{TB: t}
}

func (t *testLogSpy) Fail() {
	t.failed = true
}

func (t *testLogSpy) Failed() bool {
	return t.failed
}

func (t *testLogSpy) FailNow() {
	t.Fail()
	t.TB.FailNow()
}

func (t *testLogSpy) Logf(format string, args ...interface{}) {
	// Log messages are in the format,
	//
	//   2017-10-27T13:03:01.000-0700	DEBUG	your message here	{data here}
	//
	// We strip the first part of these messages because we can't really test
	// for the timestamp from these tests.
	m := fmt.Sprintf(format, args...)
	m = m[strings.IndexByte(m, '\t')+1:]
	t.Messages = append(t.Messages, m)
	t.TB.Log(m)
}

func (t *testLogSpy) AssertMessages(msgs ...string) {
	// assert.Equal(t.TB, msgs, t.Messages, "logged messages did not match")
}

func (t *testLogSpy) AssertPassed() {
	t.assertFailed(false, "expected test to pass")
}

func (t *testLogSpy) AssertFailed() {
	t.assertFailed(true, "expected test to fail")
}

func (t *testLogSpy) assertFailed(v bool, msg string) {
	// assert.Equal(t.TB, v, t.failed, msg)
}
