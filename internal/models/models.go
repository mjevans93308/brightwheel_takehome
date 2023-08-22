package models

import "time"

type Reading struct {
	Timestamp time.Time `json:"timestamp"`
	Count     int       `json:"count"`
}

type Device struct {
	ID            string            `json:"id"`
	ReadingsMap   map[time.Time]int `json:"readings"`
	TotalCount    int
	LatestReading Reading
}
