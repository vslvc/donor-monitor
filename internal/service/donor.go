package service

import (
	"donor-monitor/internal/model"
	"donor-monitor/internal/parser"
	"errors"
	"fmt"
)

func BuildCenter(
	name     string,
	url      string,
	p parser.Parser,
)(model.DonorCenter, error){

	statuses, err := p.Parse(url)
	if err != nil{
		return model.DonorCenter{}, fmt.Errorf("err pars", err)
	}
	center := model.DonorCenter{
		NameCenter:  name,
		URL:         url,
		Statuses:    statuses,
	} 
	return center, nil
}

func ReadStructStatus(s model.DonorStatus)(model.DonorStatus, error){
	var supColor string
	switch(s.Color){
	case "red":
		supColor = "krasny penis"
	case "yellow":
		supColor = "jeltiy penis"
	case "green": 	
		supColor = "zzeleny penis"
	default: 
		return s, errors.New("anone color" + s.Color)
	}
	s.Color = supColor
	return s, nil
}

func FindByGroup(c model.DonorCenter, group string)(model.DonorStatus, bool){
	for _, value := range c.Statuses{
		if value.BloodGroup == group{
			return value, true
		}
	}
	return model.DonorStatus{}, false
}

func FindByColor(c model.DonorCenter, color string) []model.DonorStatus{
	findColorArr := []model.DonorStatus{}
	for _, value := range c.Statuses{
		if value.Color == color{
			findColorArr = append(findColorArr, value)
		}
	}
	return findColorArr
	
}

