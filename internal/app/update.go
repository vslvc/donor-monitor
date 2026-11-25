package app

import (
	"fmt"
	"sync"
	"log"
	"errors"
	"context"

	"donor-monitor/internal/repository"
	"donor-monitor/internal/parser"
	"donor-monitor/internal/model"
	"donor-monitor/internal/service"
)
type centerConfig struct {
	Name string 
	City string 
	URL string
}

type parseResult struct {
	center model.DonorCenter
	err error
} 

func updateAllCenters(
	ctx context.Context,
	repo *repository.PostgresRepository,
) error{

	if err := ctx.Err(); err != nil{  // проверка контекста на отмену
		return err
	}
	
	configs := []centerConfig{
		{"ЦКЛО", "Saint-P", "https://yadonor.ru/reserve/spk1558"},
		{"ALMAZOV", "Saint-P", "https://yadonor.ru/reserve/spk1810"},
		{"bloodspb26", "Saint-P", "https://yadonor.ru/reserve/spk1048"},
		{"yadonor", "Saint-P", "https://yadonor.ru/reserve/yadonorspb"},	
	}
	
	results := make(chan parseResult, len(configs))
	var wg sync.WaitGroup
	for _, cfg := range configs{
		wg.Add(1)

		go func(cfg centerConfig){
			defer wg.Done()

			center, err := service.BuildCenter(
				ctx,
				cfg.Name,
				cfg.City,
				cfg.URL,
				parser.YaDonor{},
			)

			if err != nil{
				results <- parseResult{
					err: fmt.Errorf("parse %s %w ", cfg.Name, err),
				}
				return
			}
			results	<- parseResult{
				center: center,
			}
		}(cfg)
	}
	wg.Wait()
	close(results)

	if err := ctx.Err(); err != nil{
		return fmt.Errorf("update canelled: %w ", err)
	}

	centers := make([]model.DonorCenter, 0, len(configs))
	var problems []error

	for result := range results{
		if result.err != nil{
			problems = append(problems, result.err)
			continue
		}
		centers = append(centers, result.center)
	}

	if len(centers) == 0 {
		return fmt.Errorf("no centers parsed: %w", errors.Join(problems...))
	}

	centers, err := repo.SaveCenters(centers)
	if err != nil{
		problems = append(problems, fmt.Errorf("save centers: %w", err))
		return errors.Join(problems...)
	}

	for _, center := range centers{
		err := repo.SaveCheck(center.ID, center.Statuses)

		if err != nil{
			problems = append(problems,
				fmt.Errorf("save check %s %w", center.NameCenter, err),
			)
			continue
		}
	        log.Printf("Updated: %s (ID: %d)", center.NameCenter, center.ID)
	}
	return errors.Join(problems...)
}

