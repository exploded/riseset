DECLARE FUNCTION hm# (ut AS DOUBLE)
DECLARE FUNCTION sinalt# (iobj AS INTEGER, mjd0 AS DOUBLE, hour AS DOUBLE, glong AS DOUBLE, cphi AS DOUBLE, sphi AS DOUBLE)
DECLARE SUB quad (ym AS DOUBLE, y0 AS DOUBLE, yp AS DOUBLE, xe AS DOUBLE, ye AS DOUBLE, z1 AS DOUBLE, z2 AS DOUBLE, nz AS INTEGER)
DECLARE FUNCTION fpart# (x AS DOUBLE)
DECLARE FUNCTION lmst# (mjd AS DOUBLE, lambda AS DOUBLE)
DECLARE FUNCTION calday$ (mjd AS DOUBLE)
DECLARE FUNCTION ipart# (x AS DOUBLE)
DECLARE FUNCTION cn# (x AS DOUBLE)
DECLARE FUNCTION mjd# (y AS INTEGER, m AS INTEGER, d AS INTEGER, h AS DOUBLE)
DECLARE FUNCTION sn# (x AS DOUBLE)
DECLARE SUB moon (t AS DOUBLE, ra AS DOUBLE, dec AS DOUBLE)
DECLARE SUB sun (t AS DOUBLE, ra AS DOUBLE, dec AS DOUBLE)
DECLARE SUB quad (ym AS DOUBLE, y0 AS DOUBLE, yp AS DOUBLE, xe AS DOUBLE, ye AS DOUBLE, z1 AS DOUBLE, z2 AS DOUBLE, nz AS INTEGER)
'       Rise and set times for Sun and Moon
'       Adapted and modified from Montenbruck
'       and Pfleger, 'Astronomy on the personal
'       Computer' 3rd Edition, Springer
'       section 3.8
'       Accuracy of detection of 'always below' and 'always above'
'       situations depends on the approximate routines used for Sun
'       and Moon. For instance, 1999 Dec 25th, at 0 long, 67.43 lat
'       this program will give an 8 minute long day between sunrise
'       and sunset. More accurate programs say the Sun is always below
'       the horizon on this day.
'
p$ = "    ####"
DEFDBL A-Z
pi = 4 * ATN(1)
rads = pi / 180
degs = 180 / pi
DIM sinho(3)
DIM obname$(5)
obname$(1) = "Moon"
obname$(2) = "Sun"
obname$(3) = "Nautical twilight"
CLS
PRINT "   Rise and set for Sun and Moon"
PRINT "   ============================="
PRINT
INPUT "   Year (yyyy) - - - - - - - - :", y%
INPUT "   Month  (mm) - - - - - - - - :", m%
INPUT "   Day    (dd) - - - - - - - - :", d%
INPUT "   Time zone (East +) - - - -  :", zone
INPUT "   Longitude (w neg, decimals) :", glong
INPUT "   Latitude  (n pos, decimals) :", glat
glong = -glong  'routines use east longitude negative convention
zone = zone / 24
date = mjd(y%, m%, d%, 0#) - zone
'define the altitudes for each object
'treat twilight as a separate object 3, so sinalt routine
'falls through to finding Sun altitude again
sl = sn(glat)
cl = cn(glat)
sinho(1) = sn(8! / 60!)         'moonrise - average diameter used
sinho(2) = sn(-50! / 60!)       'sunrise - classic value for refraction
sinho(3) = sn(-12!)             'nautical twilight
xe = 0
ye = 0
z1 = 0
z2 = 0
FOR iobj% = 1 TO 3
    utrise = 0
    utset = 0
    rise = 0
    sett = 0
    hour = 1
    zero2 = 0
    ' See STEP 1 and 2 of Web page description.
    ym = sinalt(iobj%, date, hour - 1, glong, cl, sl) - sinho(iobj%)
    IF ym > 0! THEN above = 1 ELSE above = 0
    'used later to classify non-risings
    DO
        'STEP 1 and STEP 3 of Web page description
        y0 = sinalt(iobj%, date, hour, glong, cl, sl) - sinho(iobj%)
        yp = sinalt(iobj%, date, hour + 1, glong, cl, sl) - sinho(iobj%)
        xe = 0
        ye = 0
        z1 = 0
        z2 = 0
        nz% = 0
        'STEP 4 of web page description
        quad ym, y0, yp, xe, ye, z1, z2, nz%
        SELECT CASE nz%
            'cases depend on values of discriminant - inner part of STEP 4
            CASE 0 'nothing  - go to next time slot
            CASE 1                      ' simple rise / set event
                IF (ym < 0!) THEN       ' must be a rising event
                        utrise = hour + z1
                        rise = 1
                ELSE                    ' must be setting
                        utset = hour + z1
                        sett = 1
                END IF
            CASE 2                      ' rises and sets within interval
                IF (ye < 0!) THEN       ' minimum - so set then rise
                        utrise = hour + z2
                        utset = hour + z1
                ELSE                    ' maximum - so rise then set
                        utrise = hour + z1
                        utset = hour + z2
                END IF
                rise = 1
                sett = 1
                zero2 = 1
            END SELECT
        ym = yp     'reuse the ordinate in the next interval
        hour = hour + 2
    ' STEP 5 of Web page description - have we finished for this object?
    LOOP UNTIL (hour = 25) OR (rise * sett = 1)
    utrise = hm(utrise)
    utset = hm(utset)
    'STEP 6 of Web page description
    PRINT
    PRINT "   "; obname$(iobj%)
    ' logic to sort the various rise and set states
    IF (rise = 1 OR sett = 1) THEN   'current object rises and sets today
        IF rise = 1 THEN
            PRINT USING p$; utrise
        ELSE
            PRINT "    ----"
        END IF
        IF sett = 1 THEN
            PRINT USING p$; utset
        ELSE
            PRINT "    ----"
        END IF
    ELSE              'current object not so simple
        IF above = 1 THEN
            SELECT CASE iobj%
                CASE 1, 2: PRINT "    always above horizon"
                CASE 3: PRINT "    always bright"
            END SELECT
        ELSE
            SELECT CASE iobj%
                CASE 1, 2: PRINT "    always below horizon"
                CASE 3: PRINT "    always dark"
            END SELECT
        END IF
    END IF
'STEP 7 of Web page description
NEXT iobj%
END


DEFSNG A-Z
FUNCTION calday$ (x AS DOUBLE)
'    returns calendar date as a string in international format
'    given the modified julian date
'    BC dates are in calendar format - i.e. no year zero
'    Gregorian dates are returned after 1582 Oct 10th
'    In English colonies and Sweeden, this does not reflect
'    historical dates
jd# = x + 2400000.5#
jd0 = ipart(jd# + .5)
IF jd0 < 2299161# THEN
    c = jd0 + 1524#
ELSE
    b = ipart((jd0 - 1867216.25#) / 36524.25#)
    c = jd0 + (b - ipart(b / 4)) + 1525#
END IF
d = ipart((c - 122.1#) / 365.25#)
e = 365# * d + ipart(d / 4)
F = ipart((c - e) / 30.6001)
day = ipart(c - e + .5) - ipart(30.6001 * F)
month = F - 1 - 12 * ipart(F / 14)
year = d - 4715 - ipart((month + 7) / 10)
calday$ = STR$(year) + STR$(month) + STR$(day)
END FUNCTION

FUNCTION cn# (x AS DOUBLE)
cn = COS(x * .0174532925199433#)
END FUNCTION

DEFDBL A-Z
FUNCTION fpart# (x AS DOUBLE)
'       returns fractional part of a number
x = x - INT(x)
IF x < 0 THEN
   x = x + 1
END IF
fpart = x
END FUNCTION

FUNCTION hm (ut AS DOUBLE)
' returns number containing the time written in hours and minutes
' rounded to the nearest minute
ut = INT(ut * 60! + .5) / 60!   'round ut to nearest minute
h = INT(ut)
m = INT(60! * (ut - h) + .5)
hm = INT(100 * h + m)
END FUNCTION

DEFSNG A-Z
FUNCTION ipart# (x AS DOUBLE)
ipart = SGN(x) * INT(ABS(x))
END FUNCTION

DEFDBL A-Z
FUNCTION lmst# (mjd AS DOUBLE, glong AS DOUBLE)
'    returns the local siderial time for
'    the mjd and longitude specified
mjd0 = ipart(mjd)
ut = (mjd - mjd0) * 24
t = (mjd0 - 51544.5) / 36525
gmst = 6.697374558# + 1.0027379093# * ut
gmst = gmst + (8640184.812866# + (.093104 - .0000062 * t) * t) * t / 3600#
lmst = 24# * fpart((gmst - glong / 15#) / 24#)
END FUNCTION

DEFSNG A-Z
FUNCTION mjd# (y AS INTEGER, m AS INTEGER, d AS INTEGER, h AS DOUBLE)
'   returns modified julian date
'   number of days since 1858 Nov 17 00:00h
'   valid for any date since 4713 BC
'   assumes gregorian calendar after 1582 Oct 15, Julian before
'   Years BC assumed in calendar format, i.e. the year before 1 AD is 1 BC
a# = 10000# * y + 100# * m + d
IF y < 0 THEN y = y + 1
IF m <= 2 THEN
   m = m + 12
   y = y - 1
END IF
IF a# <= 15821004.1# THEN
   b = -2 + (y + 4716) \ 4 - 1179
ELSE
   b = (y \ 400) - (y \ 100) + (y \ 4)
END IF
a# = 365# * y - 679004#
mjd = a# + b + ipart(30.6001# * (m + 1)) + d + h / 24
END FUNCTION

DEFDBL A-Z
SUB moon (t AS DOUBLE, ra AS DOUBLE, dec AS DOUBLE)
' returns ra and dec of Moon to 5 arc min (ra) and 1 arc min (dec)
' for a few centuries either side of J2000.0
' Predicts rise and set times to within minutes for about 500 years
' in past - TDT and UT time diference may become significant for long
' times
p2 = 6.283185307#
ARC = 206264.8062#
COSEPS = .91748
SINEPS = .39778
L0 = fpart(.606433 + 1336.855225# * t)    'mean long Moon in revs
L = p2 * fpart(.374897 + 1325.55241# * t) 'mean anomaly of Moon
LS = p2 * fpart(.993133 + 99.997361# * t) 'mean anomaly of Sun
d = p2 * fpart(.827361 + 1236.853086# * t)'diff longitude sun and moon
F = p2 * fpart(.259086 + 1342.227825# * t)'mean arg latitude
' longitude correction terms
dL = 22640 * SIN(L) - 4586 * SIN(L - 2 * d)
dL = dL + 2370 * SIN(2 * d) + 769 * SIN(2 * L)
dL = dL - 668 * SIN(LS) - 412 * SIN(2 * F)
dL = dL - 212 * SIN(2 * L - 2 * d) - 206 * SIN(L + LS - 2 * d)
dL = dL + 192 * SIN(L + 2 * d) - 165 * SIN(LS - 2 * d)
dL = dL - 125 * SIN(d) - 110 * SIN(L + LS)
dL = dL + 148 * SIN(L - LS) - 55 * SIN(2 * F - 2 * d)
' latitude arguments
S = F + (dL + 412 * SIN(2 * F) + 541 * SIN(LS)) / ARC
h = F - 2 * d
' latitude correction terms
N = -526 * SIN(h) + 44 * SIN(L + h) - 31 * SIN(h - L) - 23 * SIN(LS + h)
N = N + 11 * SIN(h - LS) - 25 * SIN(F - 2 * L) + 21 * SIN(F - L)
lmoon = p2 * fpart(L0 + dL / 1296000#)  'Lat in rads
bmoon = (18520# * SIN(S) + N) / ARC     'long in rads
' convert to equatorial coords using a fixed ecliptic
CB = COS(bmoon)
x = CB * COS(lmoon)
V = CB * SIN(lmoon)
W = SIN(bmoon)
y = COSEPS * V - SINEPS * W
Z = SINEPS * V + COSEPS * W
rho = SQR(1# - Z * Z)
dec = (360# / p2) * ATN(Z / rho)
ra = (48# / p2) * ATN(y / (x + rho))
IF ra < 0 THEN
        ra = ra + 24#
END IF
END SUB

SUB quad (ym AS DOUBLE, y0 AS DOUBLE, yp AS DOUBLE, xe AS DOUBLE, ye AS DOUBLE, z1 AS DOUBLE, z2 AS DOUBLE, nz AS INTEGER)
'  finds a parabola through three points and returns values of
'  coordinates of extreme value (xe, ye) and zeros if any (z1, z2)
'  assumes that the x values are -1, 0, +1
nz = 0
a = .5 * (ym + yp) - y0
b = .5 * (yp - ym)
c = y0
xe = -b / (2! * a)              'x coord of symmetry line
ye = (a * xe + b) * xe + c      'extreme value for y in interval
dis = b * b - 4! * a * c        'discriminant
IF dis > 0 THEN                 'there are zeros
    dx = .5 * SQR(dis) / ABS(a)
    z1 = xe - dx
    z2 = xe + dx
    IF (ABS(z1) <= 1!) THEN nz = nz + 1     'This zero is in interval
    IF (ABS(z2) <= 1!) THEN nz = nz + 1     'This zero is in interval
    IF (z1 < -1!) THEN z1 = z2
END IF
END SUB

FUNCTION sinalt (iobj AS INTEGER, mjd0 AS DOUBLE, hour AS DOUBLE, glong AS DOUBLE, cphi AS DOUBLE, sphi AS DOUBLE)
' returns sine of the altitude of either the sun or the moon given the
' modified julian day number at midnight UT and the hour of the UT day,
' the longitude of the observer, and the sine and cosine of the latitude
' of the observer
ra = 0
dec = 0
instant = mjd0 + hour / 24#
t = (instant - 51544.5#) / 36525#
IF (iobj = 1) THEN
        moon t, ra, dec
ELSE
        sun t, ra, dec
END IF
tau = 15# * (lmst(instant, glong) - ra)   'hour angle of object
sinalt = sphi * sn(dec) + cphi * cn(dec) * cn(tau)
END FUNCTION

DEFSNG A-Z
FUNCTION sn# (x AS DOUBLE)
sn = SIN(x * .0174532925199433#)
END FUNCTION

DEFDBL A-Z
SUB sun (t AS DOUBLE, ra AS DOUBLE, dec AS DOUBLE)
' Returns RA and DEC of Sun to roughly 1 arcmin for few hundred
' years either side of J2000.0
p2 = 6.283185307#
COSEPS = .91748
SINEPS = .39778
m = p2 * fpart(.993133 + 99.997361# * t)        'Mean anomaly
dL = 6893# * SIN(m) + 72# * SIN(2 * m)          'Eq centre
L = p2 * fpart(.7859453# + m / p2 + (6191.2# * t + dL) / 1296000#)
' convert to RA and DEC - ecliptic latitude of Sun taken as zero
sl = SIN(L)
x = COS(L)
y = COSEPS * sl
Z = SINEPS * sl
rho = SQR(1# - Z * Z)
dec = (360# / p2) * ATN(Z / rho)
ra = (48# / p2) * ATN(y / (x + rho))
IF ra < 0 THEN ra = ra + 24
END SUB
