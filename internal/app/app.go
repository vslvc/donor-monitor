package app

import (
	"fmt"
	"donor-monitor/internal/model"
	// "donor-monitor/internal/parser"
	// "donor-monitor/internal/service"
	"donor-monitor/internal/config"
	"donor-monitor/internal/repository"
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"donor-monitor/internal/scheduler"
)

func Run() error {

	// var ArrCenters []model.DonorCenter

	fmt.Printf("\n ====== CONNECT DATABASE ====== \n")
	cfg, err := config.Load()
	if err != nil{
		return err
	}
	repo, err := repository.NewPostgresRepository(cfg.DatabaseURL)
	if err != nil{
		return err
	}
	defer repo.Close()
	fmt.Printf("\n ====== POSTGRESQL CONNECTED ====== \n")

	// fmt.Printf("\n ====== PARSING ЦКЛО ====== \n")
	// zklo := parser.YaDonor{}
	// urlzklo := "https://yadonor.ru/reserve/spk1558"
	// centerZKLO, err := service.BuildCenter(
	// 	"ЦКЛО", 
	// 	"Saint-P",
	// 	urlzklo, 
	// 	zklo,
	// )
	// if err != nil{ 
	// 	return err 
	// }
	// // PrintCenter(centerZKLO)
	//
	// ArrCenters = append(ArrCenters, centerZKLO)
	//
	// fmt.Printf("\n ====== PARSING ALMAZOV ====== \n")
	// almazov := parser.YaDonor{}
	// urlAlmazov := "https://yadonor.ru/reserve/spk1810"
	// centerAlmazov, err := service.BuildCenter(
	// 	"ALMAZOV", 
	// 	"Saint-P",
	// 	urlAlmazov, 
	// 	almazov,
	// )
	// if err != nil{ 
	// 	return err 
	// }
	// // PrintCenter(centerAlmazov)
	//
	// ArrCenters = append(ArrCenters, centerAlmazov)
	//
	// fmt.Printf("\n ====== PARSING BLOODSPB26====== \n")
	// blood26 := parser.YaDonor{}
	// urlBloodspb := "https://yadonor.ru/reserve/spk1048"
	// centerBloodSpb26, err := service.BuildCenter(
	// 	"bloodspb26", 
	// 	"Saint-P",
	// 	urlBloodspb, 
	// 	blood26,
	// )
	// if err != nil{ 
	// 	return err 
	// }
	// // PrintCenter(centerBloodSpb26)
	//
	// ArrCenters = append(ArrCenters, centerBloodSpb26 )
	//
	// fmt.Printf("\n ====== PARSING YADONOR ====== \n")
	// yaDonor := parser.YaDonor{}
	// url := "https://yadonor.ru/reserve/yadonorspb"
	// yadonorCenter, err := service.BuildCenter(
	// 	"yadonor", 
	// 	"Saint-P",
	// 	url, 
	// 	yaDonor,
	// )
	// if err != nil{ 
	// 	return err 
	// }
	// PrintCenter(yadonorCenter)


	// ArrCenters = append(ArrCenters, yadonorCenter)


	//
	// fmt.Printf("\n\n ===== PrintCenters =====\n\n")
	// PrintCenters(ArrCenters)
	// fmt.Printf("\n\n ===== saveCenters =====\n\n")
	//
	// centers, err := repo.SaveCenters(ArrCenters)
	// if err != nil{
	// 	return err
	// }
	//
	//
	//
	// fmt.Printf("\n ===== SAVE CENTERES =====\n")
	// PrintCenters(centers)
	//
	//
	//
	// fmt.Println("\n===== SAVE ALL CHECKS =====")
	//
	// for _, center := range centers {
	//     err := repo.SaveCheck(center.ID, center.Statuses)
	//     if err != nil {
	// 	return fmt.Errorf(
	// 	    "save check for %s: %w",
	// 	    center.NameCenter,
	// 	    err,
	// 	)
	//     }
	//
	//     fmt.Printf("Saved: %s (ID: %d)\n", center.NameCenter, center.ID)
	// }	
	


	// fmt.Printf("\n===== GET HISTORY BY GROUP FROM DB =====\n")
	// history, err := repo.GetHistory(1, "0+")
	// if err != nil {
	// 	return err
	// }
	// PrintArrayDonorStatus(history)

	// fmt.Printf("\n===== GET CENTERS FROM DB =====\n")
	// centersFromDB, err := repo.GetCenters()
	// if err != nil{
	// 	return err
	// }
	// PrintCenters(centersFromDB)

	// fmt.Printf("\n===== GET CENTERS AND GET LATEST DATA =====\n")
	//
	// for i := range centersFromDB {
	// 	latest, err := repo.GetLatest(centersFromDB[i].ID)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	centersFromDB[i].Statuses = latest
	// }
	// PrintCenters(centersFromDB)


	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	scheduler.Run(
		ctx,
		10*time.Second,
		func(jobCtx context.Context) error{
			return updateAllCenters(jobCtx, repo)
		})

	fmt.Printf("\n===== GET CENTERS FROM DB =====\n")
	centersFromDB, err := repo.GetCenters()
	if err != nil{
		return err
	}
	PrintCenters(centersFromDB)
	return nil
}

func PrintArrayDonorStatus(s []model.DonorStatus){
	for index, value := range s{
		fmt.Println(index, value)
	}
}

func PrintCenter(c model.DonorCenter){
	fmt.Println("ID: ", c.ID)
	fmt.Println("Название центра:", c.NameCenter)
	fmt.Println("Город: ", c.City)
	fmt.Println("Ссылка:", c.URL)
	PrintArrayDonorStatus(c.Statuses)
}

func PrintCenters(c []model.DonorCenter){
	for _, center := range c{
		PrintCenter(center)
		PrintArrayDonorStatus(center.Statuses)
	}
}
