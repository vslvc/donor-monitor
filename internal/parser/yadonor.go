package parser

import (
	"fmt"
	"net/http"
	"time"
	"strings"
	"donor-monitor/internal/model"
	"github.com/PuerkitoBio/goquery"

)

type YaDonor struct{}

func (p YaDonor) Parse(url string)([]model.DonorStatus, error){

	client := &http.Client{
		Timeout: 15 * time.Second,
	}
	
	// создание пустого запроса
	request, err := http.NewRequest("GET", url, nil)
	if err != nil{
		return nil, fmt.Errorf("request failed^ %w", err)
	}
	
	// настройка заголовков под браузер
	request.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	request.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8")
	request.Header.Set("Accept-Language", "ru-RU,ru;q=0.9,en-US;q=0.8,en;q=0.7")	

	// выполнение запроса
	response, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("error to get request: %w", err)
	}

	defer response.Body.Close()

	if response.StatusCode != 200{
		return nil, fmt.Errorf("Server returned error: %d", response.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil{
		return nil, fmt.Errorf("error pars html: %w", err)
	}
	var results []model.DonorStatus
	now := time.Now()
		
	columns := doc.Find(".station-dashboard__data-col")

	var parseErr error

	columns.Each(func(index int, column *goquery.Selection){
		groupText := column.Find(".station-dashboard__data-group").Text()
		groupText = strings.TrimSpace(groupText)

		groupParts := strings.Fields(groupText)

		if len(groupParts) < 1{
			if parseErr == nil{
				parseErr = fmt.Errorf("не удалось найти группу крови, текст: '%s'", groupText)
			} 
			return
		} 
		groupLetter := groupParts[0]

		rhesusSpans := column.Find(".station-dashboard__data-values span")

		rhesusSpans.Each(func(j int, span *goquery.Selection) {
			rhesusText := strings.TrimSpace(span.Text())
			rhesusSign := strings.TrimPrefix(rhesusText, "Rh")

			fullGroup := groupLetter + rhesusSign

			classes, _ := span.Attr("class")

			color := "неизвестно"

			if strings.Contains(classes, "red"){
				color = "red"
			} else if strings.Contains(classes, "yellow"){
				color = "yellow"
			} else if strings.Contains(classes, "green"){
				color = "green"
			}

			status := model.DonorStatus{
				BloodGroup: fullGroup,
				Color:      color,
				CheckedAt:  now,
			}

			results = append(results, status)

		})
	})

	if parseErr != nil{
		return nil, parseErr
	}

	return results, nil
}
