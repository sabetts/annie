package handlers

import (
	"fmt"
	"goirc/bot"
	"math"
	"time"
)

func TimeLeft(params bot.HandlerParams) error {
	left := time.Unix(2<<30, 0).Sub(time.Now())

	years := int(math.Round(left.Hours() / 24 / 365))
	months := int(math.Round(left.Hours() / 24 / 365 * 12))
	weeks := int(math.Round(left.Hours() / 24 / 7))
	days := int(math.Round(left.Hours() / 24))
	hours := int(math.Round(left.Hours()))
	minutes := int(math.Round(left.Minutes()))
	seconds := int(math.Round(left.Seconds()))

	str := fmt.Sprintf("%d years / %d months / %d weeks / %d days / %d hours / %d minutes / %d seconds\n", years, months, weeks, days, hours, minutes, seconds)

	params.Privmsgf(params.Target, str)

	return nil
}
