package g2s2

import (
	"github.com/golang/geo/s2"
	geojson "github.com/paulmach/go.geojson"
	"reflect"
	"testing"
)

func Test_getCoordinates(t *testing.T) {
	var polygon = [][][]float64{
		{
			[]float64{
				112.7930098772049,
				-7.374314770570638,
			},
			[]float64{
				112.79335856437682,
				-7.374868055202253,
			},
			[]float64{
				112.7943617105484,
				-7.375048936566447,
			},
			[]float64{
				112.79475331306458,
				-7.374878695284546,
			},
			[]float64{
				112.79483914375305,
				-7.374527572433864,
			},
			[]float64{
				112.79471039772034,
				-7.374298810426767,
			},
			[]float64{
				112.79407739639282,
				-7.373857246218312,
			},
			[]float64{
				112.79323518276215,
				-7.373915766801375,
			},
			[]float64{
				112.7930098772049,
				-7.374314770570638,
			},
		},
	}

	expected := []s2.Point{
		{struct{ X, Y, Z float64 }{X: -0.38419882729703975, Y: 0.9142851170924182, Z: -0.12835102557783717}},
		{struct{ X, Y, Z float64 }{X: -0.38420391120313574, Y: 0.9142816362497596, Z: -0.12836060233839366}},
		{struct{ X, Y, Z float64 }{X: -0.3842197615746778, Y: 0.9142745358015657, Z: -0.1283637331972294}},
		{struct{ X, Y, Z float64 }{X: -0.3842261581724513, Y: 0.9142722613442001, Z: -0.1283607865066313}},
		{struct{ X, Y, Z float64 }{X: -0.3842278325352469, Y: 0.914272410934289, Z: -0.12835470895011494}},
		{struct{ X, Y, Z float64 }{X: -0.38422597667116376, Y: 0.9142737467514089, Z: -0.1283507493247735}},
		{struct{ X, Y, Z float64 }{X: -0.3842162590023453, Y: 0.9142789034969275, Z: -0.1283431063154489}},
		{struct{ X, Y, Z float64 }{X: -0.38420276882467275, Y: 0.9142844302965083, Z: -0.12834411924527192}},
		{struct{ X, Y, Z float64 }{X: -0.38419882729703975, Y: 0.9142851170924182, Z: -0.12835102557783717}},
	}

	type args struct {
		fc *geojson.Feature
	}
	tests := []struct {
		name string
		args args
		want []s2.Point
	}{
		{
			name: "Given correct feature, should return coordinates",
			args: args{fc: &geojson.Feature{
				Type: "Feature",
				Geometry: &geojson.Geometry{
					Type:    "Polygon",
					Polygon: polygon,
				},
			}},
			want: expected,
		},
		{
			name: "Given an empty feature, should return no coordinates",
			args: args{fc: &geojson.Feature{
				Type: "Feature",
				Geometry: &geojson.Geometry{
					Type:    "Polygon",
					Polygon: [][][]float64{},
				},
			}},
			want: make([]s2.Point, 0),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getCoordinates(tt.args.fc); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getCoordinates() = %v, want %v", got, tt.want)
			}
		})
	}
}
