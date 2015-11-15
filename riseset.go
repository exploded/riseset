/*
Package riseset calculates the rise and set times for the Sun, Moon and twilight.

Copyright 2015 James McHugh

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.


GO Program by James McHugh, converted from QBASIC version developed by
Keith Burnett keith@xylem.demon.co.uk
http://www.stargazing.net/kepler/moonrise.html

Original QBASIC program adapted and modified from Montenbruck and Pfleger,
'Astronomy on the personal Computer' 3rd Edition, Springer section 3.8

Accuracy of detection of 'always below' and 'always above' situations depends
on the approximate routines used for Sun and Moon. For instance, 1999 Dec 25th,
at 0 long, 67.43 lat this program will give an 8 minute long day between sunrise
and sunset. More accurate programs say the Sun is always below the horizon on
this day.

*/
package github.com/exploded/riseset

import (
	"math"
	"strconv"
	"time"
)

var t, ra, dec, ym, y0, yp, xe, ye, z1, z2 float64
var nz int

/*
RiseSet holds the rise and set times as strings in the form hh:mm
*/
type RiseSet struct {
	Rise string
	Set  string
}

// Object specifies the astronomical object to calculate i.e. Sun, Moon or twilight
type Object int

const (
	Moon Object = 1 + iota
	Sun
	Twilight
)

/*
Riseset calculates the rise and set time for a given object, date, location and timezone

Example

object : 1 = Moon
		 2 = Sun
		 2 = Nautical twilight
year   : 2015
month  : 10
day    : 21
zone   : 11.         Time zone in decimal, East is +ve, West is -ve
glong  : 144.966944  Longitude in decimal, East is +ve, West is -ve
glat   : -37.816944  Latitude  in decimal, North is +ve, South is -ve

*/
func Riseset(object Object, eventdate time.Time, glong float64, glat float64, zone float64) (results RiseSet) {
	sinho := make([]float64, 4)
	
	day := eventdate.Day()
	month := int(eventdate.Month())
	year := eventdate.Year()

	glong = -glong // Routines use east longitude negative convention
	zone = zone / 24
	date := mjd(year, month, day, 0) - zone
	//define the altitudes for each object
	//treat twilight a separate object #3, so sinalt routine
	//falls through to finding Sun altitude again
	sl := sn(glat)
	cl := cn(glat)
	sinho[1] = sn(8. / 60.)   //moonrise - average diameter used
	sinho[2] = sn(-50. / 60.) //sunrise - classic value for refraction
	sinho[3] = sn(-12)        //nautical twilight
	xe = 0.
	ye = 0.
	z1 = 0.
	z2 = 0.

	iobj := object
	utrise := 0.
	utset := 0.
	rise := 0
	sett := 0
	hour := 1.
	ym = sinalt(iobj, date, hour-1, glong, cl, sl) - sinho[iobj]

	for (hour != 25) && (rise*sett != 1) {
		y0 = sinalt(iobj, date, hour, glong, cl, sl) - sinho[iobj]
		yp = sinalt(iobj, date, hour+1, glong, cl, sl) - sinho[iobj]
		xe = 0
		ye = 0
		z1 = 0
		z2 = 0
		nz = 0

		quad() // Note uses and updates package variables

		switch nz {
		//cases depend on values of discriminant
		case 0: //nothing  - go to next time slot
		case 1: // simple rise / set event
			if ym < 0 { // must be a rising event
				utrise = float64(hour) + z1
				rise = 1
			} else { // must be setting
				utset = float64(hour) + z1
				sett = 1
			}
		case 2: // rises and sets within interval
			if ye < 0 { // minimum - so set then rise
				utrise = float64(hour) + z2
				utset = float64(hour) + z1
			} else { // maximum - so rise then set
				utrise = float64(hour) + z1
				utset = float64(hour) + z2
			}
			rise = 1
			sett = 1
		}
		ym = yp //reuse the ordinate in the next interval
		hour = hour + 2
	}

	if rise == 1 {
		results.Rise = hm(utrise)
	} else {
		results.Rise = "-"
	}

	if sett == 1 {
		results.Set = hm(utset)
	} else {
		results.Set = "-"
	}

	return
}

/*
Returns calendar date a string in international format given the modified julian
date.
BC dates are in calendar format - i.e. no year zero
Gregorian dates are returned after 1582 Oct 10th
In English colonies and Sweeden, this does not reflect historical dates.
*/
func calday(x float64) string {
	var b, c, d, e, F, jd, jd0 float64
	var monthx, dayx, yearx int
	jd = x + 2400000.5
	jd0 = ipart(jd + .5)
	if jd0 < 2299161 {
		c = jd0 + 1524
	} else {
		b = ipart((jd0 - 1867216.25) / 36524.25)
		c = jd0 + (b - ipart(b/4)) + 1525
	}
	d = ipart((c - 122.1) / 365.25)
	e = 365*d + ipart(d/4.)
	F = float64(ipart((c - e) / 30.6001))
	dayx = int(ipart(c-e+.5) - ipart(30.6001*F))
	monthx = int(F - 1 - 12*ipart(F/14))
	yearx = int(d - 4715. - ipart((float64(monthx)+7.)/10.))
	return strconv.Itoa(yearx) + strconv.Itoa(monthx) + strconv.Itoa(dayx)
}

/*
Returns string containing the time written in hours and minutes rounded to
the nearest minute
*/
func hm(ut float64) (hhmm string) {
	var ut2 float64
	ut2 = float64(int(ut*60+0.5)) / 60. //round ut to nearest minute
	h := int(ut2)
	month := int(60.*(float64(ut2)-float64(h)) + 0.5)

	hhmm = strconv.Itoa(month)
	if month < 10 {
		hhmm = "0" + hhmm
	}
	hhmm = strconv.Itoa(h) + ":" + hhmm
	if h < 10 {
		hhmm = "0" + hhmm
	}
	return
}

/*
Returns modified julian date number of days since 1858 Nov 17 00:00h
Valid for any date since 4713 BC
Assumes gregorian calendar after 1582 Oct 15, Julian before
Years BC assumed in calendar format, i.e. the year before 1 AD is 1 BC
*/
func mjd(year int, month int, day int, h float64) float64 {
	// Note: the original code used the QBASIC "\" operator for integer division
	var a float64
	var b int
	a = float64(10000*year + 100*month + day)
	if year < 0 {
		year = year + 1
	}
	if month <= 2 {
		month = month + 12
		year = year - 1
	}
	if a <= 15821004.1 {
		b = -2 + (year+4716)/4 - 1179 // Integer division is intentional
	} else {
		b = (year / 400) - (year / 100) + (year / 4) // Integer division is intentional
	}
	a = 365*float64(year) - 679004
	return a + float64(b) + ipart(30.6001*(float64(month)+1)) + float64(day) + float64(h)/24
}

/*
Returns the local siderial time for the mjd and longitude specified
*/
func lmst(mjd float64, glong float64) float64 {
	mjd0 := ipart(mjd)
	ut := (mjd - mjd0) * 24
	t := (mjd0 - 51544.5) / 36525
	gmst := 6.697374558 + 1.0027379093*ut
	gmst = gmst + (8640184.812866+(.093104-.0000062*t)*t)*t/3600
	return 24 * fpart((gmst-glong/15)/24)
}

// Returns fractional part of a number.
func fpart(x float64) float64 {
	_, x = math.Modf(x) // ignore the integer part
	return x
}

// Returns the integer part of a number as a float
func ipart(x float64) float64 {
	return float64(int(x))
}

/*
Finds a parabola through three points and returns values of coordinates of
extreme value (xe, ye) and zeros if any (z1, z2)
Assumes that the x values are -1, 0, +1
*/
func quad() {
	var a, b, c float64
	nz = 0
	a = 0.5*(ym+yp) - y0
	b = 0.5 * (yp - ym)
	c = y0
	xe = -b / (2 * a)    //x coord of symmetry line
	ye = (a*xe+b)*xe + c //extreme value for y in interval
	dis := b*b - 4*a*c   //discriminant
	if dis > 0 {         //there are zeros
		dx := 0.5 * math.Sqrt(dis) / math.Abs(a)
		z1 = xe - dx
		z2 = xe + dx
		if math.Abs(z1) <= 1 {
			nz = nz + 1 //This zero is in interval
		}
		if math.Abs(z2) <= 1 {
			nz = nz + 1 //This zero is in interval
		}
		if z1 < -1 {
			z1 = z2
		}
	}
	return
}

//Returns SIN of x degrees
func cn(x float64) float64 {
	return math.Cos(x * .0174532925199433)
}

//Returns COS of x degrees
func sn(x float64) float64 {
	return math.Sin(x * .0174532925199433)
}

/*
Returns sine of the altitude of either the sun or the moon given the modified
julian day number at midnight UT and the hour of the UT day, the longitude of
the observer, and the sine and cosine of the latitude of the observer.
*/
func sinalt(iobj Object, mjd0 float64, hour float64, glong float64, cphi float64, sphi float64) float64 {
	instant := mjd0 + hour/24.
	t = (instant - 51544.5) / 36525
	if iobj == 1 {
		moonsub()
	} else {
		sun()
	}
	tau := 15 * (lmst(instant, glong) - ra) //hour angle of object
	return sphi*sn(dec) + cphi*cn(dec)*cn(tau)
}

/*
Returns RA and DEC of Sun to roughly 1 arcmin for few hundred years either side
of J2000.0
*/
func sun() {
	p2 := 6.283185307
	COSEPS := .91748
	SINEPS := .39778
	m := p2 * fpart(.993133+99.997361*t)      //Mean anomaly
	dL := 6893*math.Sin(m) + 72*math.Sin(2*m) //Eq centre
	L := p2 * fpart(.7859453+m/p2+(6191.2*t+dL)/1296000)
	// convert to RA and DEC - ecliptic latitude of Sun taken zero
	sl := math.Sin(L)
	x := math.Cos(L)
	y := COSEPS * sl
	Z := SINEPS * sl
	rho := math.Sqrt(1 - Z*Z)
	dec = (360 / p2) * math.Atan(Z/rho)
	ra = (48 / p2) * math.Atan(y/(x+rho))
	if ra < 0 {
		ra = ra + 24
	}
	return
}

/*
Returns ra and dec of Moon to 5 arc min (ra) and 1 arc min (dec) for a few
centuries either side of J2000.0
Predicts rise and set times to within minutes for about 500 years in past
TDT and UT time diference may become significant for long times
*/
func moonsub() {
	p2 := 6.283185307
	ARC := 206264.8062
	COSEPS := .91748
	SINEPS := .39778
	L0 := fpart(.606433 + 1336.855225*t)   //mean long Moon in revs
	L := p2 * fpart(.374897+1325.55241*t)  //mean anomaly of Moon
	LS := p2 * fpart(.993133+99.997361*t)  //mean anomaly of Sun
	d := p2 * fpart(.827361+1236.853086*t) //diff longitude sun and moon
	F := p2 * fpart(.259086+1342.227825*t) //mean arg latitude
	// longitude correction terms
	dL := 22640*math.Sin(L) - 4586*math.Sin(L-2*d)
	dL = dL + 2370*math.Sin(2*d) + 769*math.Sin(2*L)
	dL = dL - 668*math.Sin(LS) - 412*math.Sin(2*F)
	dL = dL - 212*math.Sin(2*L-2*d) - 206*math.Sin(L+LS-2*d)
	dL = dL + 192*math.Sin(L+2*d) - 165*math.Sin(LS-2*d)
	dL = dL - 125*math.Sin(d) - 110*math.Sin(L+LS)
	dL = dL + 148*math.Sin(L-LS) - 55*math.Sin(2*F-2*d)
	// latitude arguments
	S := F + (dL+412*math.Sin(2*F)+541*math.Sin(LS))/ARC
	h := F - 2*d
	// latitude correction terms
	N := -526*math.Sin(h) + 44*math.Sin(L+h) - 31*math.Sin(h-L) - 23*math.Sin(LS+h)
	N = N + 11*math.Sin(h-LS) - 25*math.Sin(F-2*L) + 21*math.Sin(F-L)
	lmoon := p2 * fpart(L0+dL/1296000)     //Lat in rads
	bmoon := (18520*math.Sin(S) + N) / ARC //long in rads
	// convert to equatorial coords using a fixed ecliptic
	CB := math.Cos(bmoon)
	x := CB * math.Cos(lmoon)
	V := CB * math.Sin(lmoon)
	W := math.Sin(bmoon)
	y := COSEPS*V - SINEPS*W
	Z := SINEPS*V + COSEPS*W
	rho := math.Sqrt(1 - Z*Z)
	dec = (360 / p2) * math.Atan(Z/rho)
	ra = (48 / p2) * math.Atan(y/(x+rho))
	if ra < 0 {
		ra = ra + 24
	}
	return
}
