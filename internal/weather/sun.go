package weather

import (
	"errors"
	"math"
	"time"
)

type location struct {
	lat float64
	lon float64
	tz  *time.Location
}

type sunMovement struct {
	sunrise time.Time
	sunset  time.Time
}

type Sun struct {
	location    location
	date        time.Time
	sunMovement *sunMovement
}

func (s *Sun) Date(date time.Time) {
	s.date = time.Date(date.Year(), date.Month(), date.Day(), 12, 0, 0, 0, s.location.tz)
	s.sunMovement = nil
}

func (s *Sun) Sunrise() time.Time {
	if s.sunMovement == nil {
		s.updateCalcOld()
	}
	return s.sunMovement.sunrise
}

func (s *Sun) Sunset() time.Time {
	if s.sunMovement == nil {
		s.updateCalcOld()
	}
	return s.sunMovement.sunset
}

func (s *Sun) updateCalcOld() error {
	julian := gregorianToJulian(s.date)
	julianCentury := (julian - 2451545.0) / 36525.0

	geomMeanLongSun := math.Mod(280.46646+julianCentury*(36000.76983+julianCentury*0.0003032), 360.0)
	geomMeanAnomSun := 357.52911 + julianCentury*(35999.05029-0.0001537*julianCentury)
	eccentEarthOrbit := 0.016708634 - julianCentury*(0.000042037+0.0000001267*julianCentury)
	meanObliqCorr := 23.0 + (26.0+(21.448-julianCentury*(46.815+julianCentury*(0.00059-julianCentury*0.001813)))/60.0)/60.0
	obliqCorr := meanObliqCorr + 0.00256*math.Cos(degToRad(125.04-1934.136*julianCentury))
	varY := math.Tan(degToRad(obliqCorr/2.0)) * math.Tan(degToRad(obliqCorr/2.0))
	eqOfTime := 4.0 * radToDeg(varY*math.Sin(2.0*degToRad(geomMeanLongSun))-2.0*eccentEarthOrbit*math.Sin(degToRad(geomMeanAnomSun))+4.0*eccentEarthOrbit*varY*math.Sin(degToRad(geomMeanAnomSun))*math.Cos(2.0*degToRad(geomMeanLongSun))-0.5*varY*varY*math.Sin(4.0*degToRad(geomMeanLongSun))-1.25*eccentEarthOrbit*eccentEarthOrbit*math.Sin(2.0*degToRad(geomMeanAnomSun)))

	_, offset := s.date.Zone()
	solarNoonLST := ((720.0 - 4.0*s.location.lon - eqOfTime + float64(offset)*60.0) / 1440.0)

	sunEqOfCtr := math.Sin(degToRad(geomMeanAnomSun))*(1.914602-julianCentury*(0.004817+0.000014*julianCentury)) + math.Sin(degToRad(2.0*geomMeanAnomSun))*(0.019993-0.000101*julianCentury) + math.Sin(degToRad(3.0*geomMeanAnomSun))*0.000289
	sunTrueLon := geomMeanLongSun + sunEqOfCtr
	sunAppLon := sunTrueLon - 0.00569 - 0.00478*math.Sin(degToRad(125.04-1934.136*julianCentury))
	sunDecl := radToDeg(math.Asin(math.Sin(degToRad(obliqCorr)) * math.Sin(degToRad(sunAppLon))))
	haSunrise := radToDeg(math.Acos(math.Cos(degToRad(90.833))/(math.Cos(degToRad(s.location.lat))*math.Cos(degToRad(sunDecl))) - math.Tan(degToRad(s.location.lat))*math.Tan(degToRad(sunDecl))))

	sunriseLST := (solarNoonLST*1440.0 - haSunrise*4.0) / 1440.0
	sunsetLST := (solarNoonLST*1440.0 + haSunrise*4.0) / 1440.0

	s.sunMovement = &sunMovement{
		sunrise: time.Date(s.date.Year(), s.date.Month(), s.date.Day(), int(math.Floor(sunriseLST*24.0)), int(math.Floor(math.Mod(sunriseLST*1440.0, 60))), int(math.Floor(math.Mod((sunriseLST*86400), 60))), 0, s.location.tz),
		sunset:  time.Date(s.date.Year(), s.date.Month(), s.date.Day(), int(math.Floor(sunsetLST*24.0)), int(math.Floor(math.Mod(sunsetLST*1440.0, 60))), int(math.Floor(math.Mod((sunsetLST*86400), 60))), 0, s.location.tz),
	}
	return nil
}

func (s *Sun) updateCalc() error {
	julian := gregorianToJulian(s.date)
	n := julianDay(julian)
	meanSolarTime := meanSolarTime(n, s.location.lon)
	mDegrees := math.Mod(357.5291+0.98560028*meanSolarTime, 360.0)
	mRadians := degToRad(mDegrees)
	cDegrees := 1.9148*math.Sin(mRadians) + 0.02*math.Sin(2.0*mRadians) + 0.0003*math.Sin(3.0*mRadians)
	lDegrees := math.Mod(mDegrees+cDegrees+180.0+102.9372, 360.0)
	lRadians := degToRad(lDegrees)
	julianTransit := 2451545.0 + meanSolarTime + 0.0053*math.Sin(mRadians) - 0.0069*math.Sin(2.0*lRadians)

	dSin := math.Sin(lRadians) * math.Sin(degToRad(23.4397))
	dCos := math.Cos(math.Asin(dSin))

	someCos := (math.Sin(degToRad(-0.833-2.076*math.Sqrt(0.0)/60.0)) - math.Sin(degToRad(s.location.lat))*dSin/(math.Cos(degToRad(s.location.lat))*dCos))
	w0Radians := math.Acos(someCos)
	if math.IsNaN(w0Radians) {
		return errors.New("w0Radians is not a number")
	}
	w0Degrees := radToDeg(w0Radians)
	s.sunMovement = &sunMovement{
		sunrise: julianToGregorian(julianTransit-w0Degrees/360.0, s.date.Location()),
		sunset:  julianToGregorian(julianTransit+w0Degrees/360.0, s.date.Location()),
	}
	return nil
}

func gregorianToJulian(greg time.Time) float64 {
	m := int(greg.Month())
	y := greg.Year()
	if int(greg.Month()) <= 2 {
		y = y - 1
		m = 12 + m
	}
	b := float64(2) - math.Floor(float64(y)/float64(100)) + math.Floor(float64(y)/float64(400))
	d := float64(greg.Day()) + float64(greg.Hour())/24.0 + float64(greg.Minute())/1440.0 + float64(greg.Second())/86400.0
	jd := math.Floor(365.25*(float64(y)+float64(4716))) + math.Floor(30.6001*(float64(m+1))) + d + b - 1524.5
	_, offset := greg.Zone()
	return (jd - float64(offset)/24.0)
}

func julianToGregorian(jd float64, tz *time.Location) time.Time {
	z := math.Floor(jd + 0.5)
	f := jd + 0.5 - z

	alpha := math.Floor((z - 1867216.25) / 36524.25)
	a := z + 1.0 + alpha - math.Floor(alpha/4.0)
	b := a + 1524.0
	c := math.Floor((b - 122.1) / 365.25)
	d := math.Floor(365.25 * c)
	e := math.Floor((b - d) / 30.6001)

	rawDay := b - d - math.Floor(30.6001*e) + f
	day := math.Floor(rawDay)
	hours := math.Floor((rawDay - day) * 24.0)
	minutes := math.Floor(math.Mod((rawDay-day)*1440.0, 60))
	seconds := math.Floor(math.Mod((rawDay-day)*86400.0, 60))
	nanos := (math.Mod((rawDay-day)*86400.0, 60) - seconds) * 1000000000
	var month int
	var year int
	if e <= 13 {
		month = int(e - 1.0)
		year = int(c - 4716)
	} else {
		month = int(e - 13.0)
		year = int(c - 4715)
	}
	return time.Date(year, time.Month(month), int(day), int(hours), int(minutes), int(seconds), int(nanos), tz)
}

func julianDay(jDate float64) int {
	return int(math.Ceil(jDate - (2451545.0 + 0.0009) + 69.184/86400.0))
}

func meanSolarTime(jDay int, lon float64) float64 {
	return float64(jDay) + 0.0009 - lon/360.0
}

func degToRad(deg float64) float64 {
	return math.Pi * deg / 180.0
}

func radToDeg(rad float64) float64 {
	return 180.0 * rad / math.Pi
}
