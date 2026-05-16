package main

import (
	"math"
	"time"

	"github.com/notwithering/kdialog"
)

var signs = []string{
	"Aries", "Taurus", "Gemini", "Cancer", "Leo", "Virgo",
	"Libra", "Scorpio", "Sagittarius", "Capricorn", "Aquarius", "Pisces",
}

func main() {
	t := kdialog.DialogBox{
		Form:  kdialog.FormCalendar,
		Text:  "When were you born?",
		Title: "Zodiac Sign",
	}.MustRun().(time.Time)

	// begin math i dont understand
	jd := float64(t.Unix())/86400.0 + 2440587.5
	n := jd - 2451545.0

	L := math.Mod(280.460+0.9856474*n, 360)
	g := math.Mod(357.528+0.9856003*n, 360) * math.Pi / 180

	lambda := L + 1.915*math.Sin(g) + 0.020*math.Sin(2*g)
	longitude := math.Mod(lambda+360, 360)
	// end math i dont understand

	sign := signs[int(longitude/30)]

	kdialog.DialogBox{
		Form:  kdialog.FormMsgBox,
		Text:  "Your zodiac sign is: " + sign,
		Title: "Zodiac Sign",
	}.MustRun()
}
