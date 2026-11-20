package service

import (
	"donor-monitor/internal/model"
	"donor-monitor/internal/parser"
	"errors"
	"fmt"
)

func BuildCenter(
	name     string,
	city     string,
	url      string,
	p parser.Parser,
)(model.DonorCenter, error){

	statuses, err := p.Parse(url)
	if err != nil{
		return model.DonorCenter{}, err 
	}
	
	if err := ValidateStatuses(statuses); err != nil{
		return model.DonorCenter{}, fmt.Errorf(
			"validation failed for %s: %w", name, err,
		)
	}

	center := model.DonorCenter{
		NameCenter:  name,
		City:        city,
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

func ValidateStatuses(statuses []model.DonorStatus) error{
	if len(statuses) != 8{
		return fmt.Errorf("not 8 status: %d", len(statuses))
	}
	
	validGroups := map[string]bool{
		"0+": true,
		"0-": true,
		"A+": true,
		"A-": true,
		"B+": true,
		"B-": true,
		"AB+": true,
		"AB-": true,
	}
	validColors := map[string]bool{
        	"red":    true,
        	"yellow": true,
        	"green":  true,
    	}

	seenGroups := make(map[string]bool)

	checkedAt := statuses[0].CheckedAt

	for _, status := range statuses{

		if !validGroups[status.BloodGroup]{
			return fmt.Errorf("invalid blood group: %s", status.BloodGroup)
		}

		if seenGroups[status.BloodGroup]{
			return fmt.Errorf("duplicate blood group: %s", status.BloodGroup)
		}

		seenGroups[status.BloodGroup] = true

		if !validColors[status.Color] {
				return fmt.Errorf("invalid color %s for group %s", status.Color, status.BloodGroup)
		}

		if status.CheckedAt.IsZero(){
			return fmt.Errorf("empty time for group %s", status.BloodGroup)
		}

		if !status.CheckedAt.Equal(checkedAt) {
			return fmt.Errorf("different check time for group %s", status.BloodGroup)
        	}
	}
	return nil
}

