package parser

import (
	"donor-monitor/internal/model"
	"context"
)

type Parser interface{
	Parse(ctx context.Context, url string) ([]model.DonorStatus, error)
}
