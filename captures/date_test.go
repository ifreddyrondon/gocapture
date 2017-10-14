package captures_test

import (
	"os"
	"testing"
	"time"

	"github.com/ifreddyrondon/gocapture/captures"
)

func TestUnmarshalJSON(t *testing.T) {
	defer os.Setenv("TZ", os.Getenv("TZ"))
	os.Setenv("TZ", "UTC")

	tt := []struct {
		name    string
		payload []byte
		result  string
	}{
		{
			"valid RFC3339 date with date key",
			[]byte(`{"date": "1989-12-26T06:01:00.00Z"}`),
			"1989-12-26 06:01:00 +0000 UTC",
		},
		{
			"valid RFC3339 date with timestamp key",
			[]byte(`{"timestamp": "1989-12-26T06:01:00.00Z"}`),
			"1989-12-26 06:01:00 +0000 UTC",
		},
		{
			"valid RFC1123 date with date key",
			[]byte(`{"date": "Tue, 26 Dec 1989 06:01:00 UTC"}`),
			"1989-12-26 06:01:00 +0000 UTC",
		},
		{
			"valid RFC1123 date with timestamp key",
			[]byte(`{"timestamp": "Tue, 26 Dec 1989 06:01:00 UTC"}`),
			"1989-12-26 06:01:00 +0000 UTC",
		},
		{
			"valid Unix timestamp as string date with date key",
			[]byte(`{"date": "630655260"}`),
			"1989-12-26 06:01:00 +0000 UTC",
		},
		{
			"valid Unix timestamp as string date with timestamp key",
			[]byte(`{"timestamp": "630655260"}`),
			"1989-12-26 06:01:00 +0000 UTC",
		},
		{
			"valid Unix timestamp as integer date with date key",
			[]byte(`{"date": 630655260}`),
			"1989-12-26 06:01:00 +0000 UTC",
		},
		{
			"valid Unix timestamp as integer date with timestamp key",
			[]byte(`{"timestamp": 630655260}`),
			"1989-12-26 06:01:00 +0000 UTC",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			result := captures.NewDate()
			result.UnmarshalJSON(tc.payload)

			if result.Timestamp.String() != tc.result {
				t.Errorf("Expected CaptureDate to be '%v'. Got '%v'", tc.result, result.Timestamp.String())
			}
		})
	}
}

func TestUnmarshalJSONWhenFails(t *testing.T) {
	timeMock := time.Date(1989, time.Month(12), 26, 6, 1, 0, 0, time.UTC)
	captures.SetClockInstance(captures.NewMockClock(timeMock))

	defer os.Setenv("TZ", os.Getenv("TZ"))
	os.Setenv("TZ", "UTC")

	tt := []struct {
		name    string
		payload []byte
		result  string
	}{
		{"empty json", []byte("{}"), "1989-12-26 06:01:00 +0000 UTC"},
		{"invalid json", []byte("`"), "1989-12-26 06:01:00 +0000 UTC"},
		{
			"missing date or timestamp",
			[]byte(`{"foo": "630655260"}`),
			"1989-12-26 06:01:00 +0000 UTC",
		},
		{
			"invalid date",
			[]byte(`{"date": "asd"}`),
			"1989-12-26 06:01:00 +0000 UTC",
		},
		{
			"invalid timestamp",
			[]byte(`{"timestamp": "asd"}`),
			"1989-12-26 06:01:00 +0000 UTC",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			result := captures.NewDate()
			result.UnmarshalJSON(tc.payload)

			if result.Timestamp.String() != tc.result {
				t.Errorf("Expected CaptureDate to be '%v'. Got '%v'", tc.result, result.Timestamp.String())
			}
		})
	}
}
