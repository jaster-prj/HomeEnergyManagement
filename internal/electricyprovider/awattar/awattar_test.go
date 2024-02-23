package awattar

import (
	"reflect"
	"testing"
	"time"
)

func TestAwattarElectricyProvider_GetMarketData(t *testing.T) {
	start := time.Now().Add(-24 * time.Hour)
	end := time.Now().Add(24 * time.Hour)
	type fields struct {
		location string
	}
	type args struct {
		start *time.Time
		end   *time.Time
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []AwattarData
		wantErr bool
	}{
		{"Awattar Request", fields{location: "de"}, args{&start, &end}, []AwattarData{}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &AwattarElectricyProvider{
				location: tt.fields.location,
			}
			got, err := a.GetMarketData(tt.args.start, tt.args.end)
			if (err != nil) != tt.wantErr {
				t.Errorf("AwattarElectricyProvider.GetMarketData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AwattarElectricyProvider.GetMarketData() = %v, want %v", got, tt.want)
			}
		})
	}
}
