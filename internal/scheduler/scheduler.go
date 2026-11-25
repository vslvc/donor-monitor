package scheduler

import (
	"context"
	"fmt"
	"time"
	"log"
)

func Run(
	ctx context.Context, 
	interval time.Duration,
	job func(context.Context) error,
){
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for{
		select{
		case <-ticker.C:
			fmt.Println("время обновлять через функцию епте ")
			if ctx.Err() != nil{
				return
			}
			if err := job(ctx); err != nil{
				log.Printf("scheduler error %v", err)
			} 

		case <-ctx.Done():
			fmt.Println("Scheduler stopped")
			return
		}
	}
}
