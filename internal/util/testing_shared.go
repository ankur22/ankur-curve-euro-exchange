package util

import (
	"bytes"
	"math"
	"net/http"
	"testing"
	"time"
)

const margin float32 = 0.0001

type ClientMock struct {
	Timeout  time.Duration
	Response http.Response
}

func (c *ClientMock) Do(req *http.Request) (*http.Response, error) {
	return &c.Response, nil
}

type ClosingBuffer struct {
	*bytes.Buffer
}

func (cb ClosingBuffer) Close() error {
	return nil
}

func GetValidResponseForExchangeRate() http.Response {
	return http.Response{
		StatusCode: 200,
		Body:       ClosingBuffer{bytes.NewBufferString("{\"rates\":{\"USD\":1.1043,\"GBP\":0.87518},\"base\":\"EUR\",\"date\":\"2019-10-11\"}")},
	}
}

func GetInalid200ResponseForExchangeRate() http.Response {
	return http.Response{
		StatusCode: 200,
		Body:       ClosingBuffer{bytes.NewBufferString("<>")},
	}
}

func GetInvalidResponseForExchangeRate() http.Response {
	return http.Response{
		StatusCode: 404,
	}
}

// AssertEquals - assert expected is equal to actual
func AssertEquals(t *testing.T, expected, actual float32) {
	t.Helper()
	if float32(math.Abs(float64(expected-actual))) > margin {
		t.Fatalf("expected '%f' to equal to actual '%f'", expected, actual)
	}
}

// AssertNotEquals - assert expected is not equal to actual
func AssertNotEquals(t *testing.T, expected, actual float32) {
	t.Helper()
	if float32(math.Abs(float64(expected-actual))) < margin {
		t.Fatalf("expected '%f' not to equal to actual '%f'", expected, actual)
	}
}

// AssertTrue - assert actual is true
func AssertTrue(t *testing.T, actual bool) {
	t.Helper()
	if !actual {
		t.Fatal("expected 'true' but actual is 'false'")
	}
}

// AssertFalse - assert actual is false
func AssertFalse(t *testing.T, actual bool) {
	t.Helper()
	if actual {
		t.Fatal("expected 'false' but actual is 'true'")
	}
}

// AssertNotNil - asserts actual is not nil
func AssertNotNil(t *testing.T, actual time.Time) {
	t.Helper()
	if actual.IsZero() {
		t.Fatal("expected non zero value for time")
	}
}

// AssertErrorNil - asserts actual is nil
func AssertErrorNil(t *testing.T, actual error) {
	t.Helper()
	if actual != nil {
		t.Fatal("expected nil value for error")
	}
}

// AssertErrorNotNil - asserts actual is not nil
func AssertErrorNotNil(t *testing.T, actual error) {
	t.Helper()
	if actual == nil {
		t.Fatal("expected not nil value for error")
	}
}
