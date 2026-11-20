package parser 

import (
	"time"
	"donor-monitor/internal/model"
)

type MockParser struct{}

func (p MockParser) Parse(url string) ([]model.DonorStatus, error){
	now := time.Now()

	statuses := []model.DonorStatus{
		{"0+", "green", now},
		{"0-", "red", now},
		{"A+", "green", now},
		{"A-", "red", now},
		{"B+", "red", now},
		{"B-", "yellow", now},
		{"AB+", "red", now},
		{"AB-", "green", now},
	}
	return statuses, nil
}

