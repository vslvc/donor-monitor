package parser

import (
	"donor-monitor/internal/model"
)

type Parser interface{
	Parse(url string) ([]model.DonorStatus, error)
}
