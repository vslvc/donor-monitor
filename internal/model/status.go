package model

import "time"

type DonorStatus struct{
	BloodGroup    string
	Color         string
	CheckedAt    time.Time
}

