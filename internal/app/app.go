package app

import (
	"fmt"
	"donor-monitor/internal/model"
	"donor-monitor/internal/parser"
	"donor-monitor/internal/service"
)

func Run() error {
	

	fmt.Printf("\n ====== TESTING ====== \n")
	


	fmt.Printf("\n ====== TESTING ====== \n")


	fmt.Printf("\n ====== PARSING ====== \n")

	yaDonor := parser.YaDonor{}

	url := "https://yadonor.ru/reserve/yadonorspb"
	
	center, err := service.BuildCenter("yadonor", url, yaDonor)
	if err != nil{
		return err
	}
	PrintCenter(center)

	fmt.Printf("\n ===== Find By Color =====\n")
	PrintArrayDonorStatus(service.FindByColor(center, "green"))

	fmt.Printf("\n ===== Find By Group =====\n")
	arr, boolev := service.FindByGroup(center, "B-")
	if boolev != true{
		fmt.Println("err")
	}
	fmt.Println(arr)

	return nil
}

func PrintArrayDonorStatus(s []model.DonorStatus){
	for index, value := range s{
		fmt.Println(index, value)
	}
}

func PrintCenter(c model.DonorCenter){
	fmt.Println("Название центра:", c.NameCenter)	
	fmt.Println("Ссылка:", c.URL)
	PrintArrayDonorStatus(c.Statuses)
}
