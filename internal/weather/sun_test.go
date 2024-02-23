package weather

import (
	"testing"
	"time"
)

var (
	timezone, _ = time.LoadLocation("Europe/Berlin")
)

func TestSun_updateCalc(t *testing.T) {
	type fields struct {
		location    location
		date        time.Time
		sunMovement *sunMovement
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
		want    *sunMovement
	}{
		{"Calculate Sunrise", fields{location: location{lat: 52.5, lon: 13.5, tz: timezone}, date: time.Date(2024, time.January, 1, 0, 0, 0, 0, time.FixedZone("CET", 1)), sunMovement: nil}, false, &sunMovement{sunrise: time.Date(2024, time.January, 1, 8, 16, 47, 0, timezone), sunset: time.Date(2024, time.January, 1, 16, 01, 21, 0, timezone)}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Sun{
				location:    tt.fields.location,
				date:        tt.fields.date,
				sunMovement: tt.fields.sunMovement,
			}
			if err := s.updateCalcOld(); (err != nil) != tt.wantErr {
				t.Errorf("Sun.updateCalc() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.want != nil {
				if s.sunMovement.sunrise != tt.want.sunrise {
					t.Errorf("Sunrise is: %v, want: %v", s.sunMovement.sunrise, tt.want.sunrise)
				}
				if s.sunMovement.sunset != tt.want.sunset {
					t.Errorf("Sunset is: %v, want: %v", s.sunMovement.sunset, tt.want.sunset)
				}
			}
		})
	}
}
