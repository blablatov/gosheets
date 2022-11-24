package cmweipost

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type DataWPoster interface {
	WPostSend() (string, error)
}

type DataPost struct {
	Indicator_to_mo_fact_id string `json:"indicator_to_mo_fact_id"`
	Indicator_to_mo_id      string `json:"indicator_to_mo_id"`
	Weight                  string `json:"weight"`
	ApiUrlw                 string
}

func (dp DataPost) WPostSend() (string, error) {
	// Request to add of weight value via API of server.
	// Запрос на добавление значения веса через API сервера.
	// apiUrl := "https://testdb.kpi-drive.ru/_api/indicators/save_indicator_instance_field"

	// Add data to request. Добавить данные в запрос.
	params := url.Values{"indicator_to_mo_fact_id": {dp.Indicator_to_mo_fact_id}}
	params.Set("indicator_to_mo_id", dp.Indicator_to_mo_id)
	params.Set("period_start", "2022-02-01")
	params.Set("period_end", "2022-02-28")
	params.Set("period_key", "month")
	params.Set("auth_user_id", "4")
	params.Set("field_name", "weight")
	params.Set("field_value", dp.Weight)

	dreq := strings.NewReader(params.Encode()) // String of data. Сформировать строку с данными.

	client := &http.Client{} // Create code of request. Сформировать код запроса.
	req, _ := http.NewRequest(http.MethodPost, dp.ApiUrlw, dreq)
	// Request headers. Заголовки запроса.
	req.Header.Set("Content-type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "d0f00715-09ad-4808-b7a7-a7208e90bdec")

	resp, err := client.Do(req) // Run request. Выполнить запрос.
	if err != nil {
		log.Fatalf("Error of response: %v", err)
	}

	defer resp.Body.Close()
	//fmt.Printf("%v ", resp.Status)
	body, err := ioutil.ReadAll(resp.Body) // Read of server data. Чтение данных сервера.
	p := recover()
	if err != nil {
		log.Println("Panic, internal error of server data. recover(): ", err)
		panic(p)
		//panic(err)
	}
	//log.Println("\nResponse of server:", string([]byte(body)))
	return string([]byte(body)), err
}
