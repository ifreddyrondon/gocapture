package geocoding_test

import (
	"testing"

	"github.com/ifreddyrondon/gocapture/geocoding"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewPointSuccess(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name     string
		lat, lng float64
	}{
		{"valid with lng upper limit", 75, 180},
		{"valid with lat upper limit", 90, -147.45},
		{"valid with decimals", 77.11112223331, 149.99999999},
		{"valid both upper limits", 90, 180},
		{"valid both lower limits", -90.00000, -180.0000},
		{"valid with just point decimal", 90., 180.},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			result, err := geocoding.New(tc.lat, tc.lng)
			require.Nil(t, err)
			require.NotNil(t, result)
			assert.Equal(t, tc.lng, *result.LNG)
			assert.Equal(t, tc.lat, *result.LAT)
		})
	}
}

func TestNewPointFailure(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name     string
		lat, lng float64
		expected error
	}{
		{"invalid lat > 90", 95, 280, geocoding.ErrorLATRange},
		{"invalid lat < -95", -95, 280, geocoding.ErrorLATRange},
		{"invalid lng > 180", 75, 280, geocoding.ErrorLONRange},
		{"invalid lng with decimals", 77.11112223331, 249.99999999, geocoding.ErrorLONRange},
		{"invalid lng for 2 decimals points", 90, 180.2, geocoding.ErrorLONRange},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			point, err := geocoding.New(tc.lat, tc.lng)
			require.Nil(t, point)
			assert.Equal(t, tc.expected, err)
		})
	}
}

func TestUnmarshalJSONSuccess(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name     string
		payload  []byte
		expected *geocoding.Point
	}{
		{
			"lat and lng",
			[]byte(`{"lat":40.7486, "lng":-73.9864}`),
			getPoint(40.7486, -73.9864),
		},
		{
			"latitude and longitude",
			[]byte(`{"latitude":1, "longitude":1}`),
			getPoint(1, 1),
		},
		{
			"mixed latitude and lng",
			[]byte(`{"latitude":1, "lng":1}`),
			getPoint(1, 1),
		},
		{
			"mixed lat and longitude",
			[]byte(`{"lat":1, "longitude":1}`),
			getPoint(1, 1),
		},
		{
			"with elevation",
			[]byte(`{"lat":1, "longitude":1, "elevation": 1}`),
			func() *geocoding.Point {
				p := getPoint(1, 1)
				elevation := 1.0
				p.Elevation = &elevation
				return p
			}(),
		},
		{
			"with altitude",
			[]byte(`{"lat":1, "longitude":1, "altitude": 1}`),
			func() *geocoding.Point {
				p := getPoint(1, 1)
				elevation := 1.0
				p.Elevation = &elevation
				return p
			}(),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			result := geocoding.Point{}
			err := result.UnmarshalJSON(tc.payload)
			require.Nil(t, err)
			require.NotNil(t, result)
			assert.Equal(t, tc.expected.LNG, result.LNG)
			assert.Equal(t, tc.expected.LAT, result.LAT)
			assert.Equal(t, tc.expected.Elevation, result.Elevation)
		})
	}
}

func TestUnmarshalJSONFailure(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name     string
		payload  []byte
		expected error
	}{
		{"invalid lat", []byte(`{"lat":100, "lng": 1}`), geocoding.ErrorLATRange},
		{"invalid lng", []byte(`{"lat":1, "lng": 190}`), geocoding.ErrorLONRange},
		{"invalid json", []byte("`"), geocoding.ErrorUnmarshalPoint},
		{"missing lat", []byte(`{"lng":1}`), geocoding.ErrorLATMissing},
		{"missing lng", []byte(`{"lat":1}`), geocoding.ErrorLNGMissing},
		{"missing latitude", []byte(`{"longitude":1}`), geocoding.ErrorLATMissing},
		{"missing longitude", []byte(`{"latitude":1}`), geocoding.ErrorLNGMissing},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			p := geocoding.Point{}
			err := p.UnmarshalJSON(tc.payload)
			assert.Equal(t, tc.expected, err)
		})
	}
}

func TestUnmarshalJSONMissingBody(t *testing.T) {
	t.Parallel()

	result := geocoding.Point{}
	err := result.UnmarshalJSON([]byte("{}"))
	require.Nil(t, err)
	require.Nil(t, result.LAT)
	require.Nil(t, result.LNG)
}

func getPoint(lat, lng float64) *geocoding.Point {
	p, _ := geocoding.New(lat, lng)
	return p
}
