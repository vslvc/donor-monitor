package app

import (
	"fmt"
	"donor-monitor/internal/model"
	// "donor-monitor/internal/parser"
	// "donor-monitor/internal/service"
	"donor-monitor/internal/config"
	"donor-monitor/internal/repository"
)

func Run() error {
	
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




	//
	// fmt.Printf("\n ====== TEST PARSING ====== \n")
	//
	// mockpars := parser.MockParser{}	
	//
	// test, err := service.BuildCenter(
	// 	"Mock", 
	// 	"Sain-P",
	// 	"", 
	// 	mockpars,
	// )
	// if err != nil{
	// 	return err
	// }
	// PrintCenter(test)





	// fmt.Printf("\n ====== PARSING ====== \n")
	//
	// yaDonor := parser.YaDonor{}
	//
	// url := "https://yadonor.ru/reserve/yadonorspb"
	//
	// center, err := service.BuildCenter(
	// 	"yadonor", 
	// 	"Saint-P",
	// 	url, 
	// 	yaDonor,
	// )
	// if err != nil{ 
	// 	return err 
	// }
	// PrintCenter(center)
	//
	//
	//
	//
	// // fmt.Printf("\n ===== Find By Color =====\n")
	// // PrintArrayDonorStatus(service.FindByColor(center, "red"))
	// //
	// // fmt.Printf("\n ===== Find By Group =====\n")
	// // arr, boolev := service.FindByGroup(center, "B-")
	// // if boolev != true{
	// // 	fmt.Println("err")
	// // }
	// // fmt.Println(arr)
	//
	// fmt.Printf("\n ===== SAVE TEST CENTER =====\n")
	//
	// id, err := repo.SaveCenter(test)
	// if err != nil{
	// 	return err
	// }
	// test.ID = id
	// fmt.Println("Test center saved. ID:", test.ID)
	//
	//
	//
	//
	// fmt.Printf("\n ===== SAVE CENTER =====\n")
	//
	// id, err = repo.SaveCenter(center)
	// if err != nil{
	// 	return err
	// }
	// center.ID = id
	// fmt.Println("Center saved. ID:", center.ID)
	//
	//
	//
	// fmt.Printf("\n ===== SAVE GROUP + CHECK =====\n")
	//
	// err = repo.SaveCheck(center.ID, center.Statuses)
	// if err != nil{
	// 	return err
	// }
	// fmt.Println("Check saved successfully")
	//
	

	fmt.Printf("\n ===== GET LATEST FROM DB=====\n")
	latest, err := repo.GetLatest(1)
	if err != nil{
		return err
	}
	PrintArrayDonorStatus(latest)

	fmt.Printf("\n===== GET HISTORY BY GROUP FROM DB =====\n")
	history, err := repo.GetHistory(1, "0+")
	if err != nil {
		return err
	}
	PrintArrayDonorStatus(history)

	fmt.Printf("\n===== GET CENTERS FROM DB =====\n")
	centersFromDB, err := repo.GetCenters()
	if err != nil{
		return err
	}
	PrintCenters(centersFromDB)

	fmt.Printf("\n===== GET CENTERS AND GET LATEST DATA =====\n")

	for i := range centersFromDB {

		latest, err = repo.GetLatest(centersFromDB[i].ID)

		if err != nil {
			return err
		}

		centersFromDB[i].Statuses = latest
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
