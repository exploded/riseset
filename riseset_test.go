/*
Package riseset calculates the rise and set times for the Sun, Moon and twilight.

# Copyright 2015 James McHugh

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
package riseset

import (
	"testing"
	"time"
)

func TestRiseset(t *testing.T) {

	var rstests = []struct {
		year  int
		month int
		day   int
		zone  float64
		lon   float64
		lat   float64
		mrise string
		mset  string
		srise string
		sset  string
		nrise string
		nset  string
	}{
		{2000, 01, 03, 0, -1.91667, 52.50, "05:01", "14:09", "08:18", "16:06", "06:53", "17:31"},
		{1999, 12, 25, 0, 00.00000, 67.43, "17:47", "12:06", "11:56", "12:04", "07:47", "16:13"},
		{2000, 01, 03, 1, 17.42000, 68.43, "06:29", "11:57", "-", "-", "07:43", "16:07"},
	}

	// Only using the date part so GMT is okay
	mytime, _ := time.LoadLocation("GMT")

	for _, tt := range rstests {

		//Moon
		got := Riseset(1, time.Date(tt.year, time.Month(tt.month), tt.day, 0, 0, 0, 0, mytime), tt.lon, tt.lat, tt.zone)
		if got.Rise != tt.mrise || got.Set != tt.mset {
			t.Errorf("Riseset(1,%v) == %v, want rise=%v set=%v", tt, got, tt.mrise, tt.mset)
		}

		//Sun
		got = Riseset(2, time.Date(tt.year, time.Month(tt.month), tt.day, 0, 0, 0, 0, mytime), tt.lon, tt.lat, tt.zone)
		if got.Rise != tt.srise || got.Set != tt.sset {
			t.Errorf("Riseset(2,%v) == %v, want rise=%v set=%v", tt, got, tt.srise, tt.sset)
		}

		//Twilight
		got = Riseset(3, time.Date(tt.year, time.Month(tt.month), tt.day, 0, 0, 0, 0, mytime), tt.lon, tt.lat, tt.zone)
		if got.Rise != tt.nrise || got.Set != tt.nset {
			t.Errorf("Riseset(3,%v) == %v, want rise=%v set=%v", tt, got, tt.nrise, tt.nset)
		}

	}
}
