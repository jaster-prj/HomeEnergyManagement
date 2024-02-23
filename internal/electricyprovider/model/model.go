package model

import "time"

type ElectricyPrice struct {
	start time.Time
	end   time.Time
	price float64
	unit  string
}
