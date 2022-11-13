package weipost

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

func PostTestdb(wg sync.WaitGroup, weight, indicator_to_mo_fact_id, indicator_to_mo_id string, cw chan string) {
	defer wg.Done()
	//Запрос для добавления значения веса методом сервера
	apiUrl := "https://testdb.kpi-drive.ru/_api/indicators/save_indicator_instance_field"
	// Добавить данные в запрос
	params := url.Values{"indicator_to_mo_fact_id": {indicator_to_mo_fact_id}}
	params.Set("indicator_to_mo_id", indicator_to_mo_id)
	params.Set("period_start", "2022-02-01")
	params.Set("period_end", "2022-02-28")
	params.Set("period_key", "month")
	params.Set("auth_user_id", "4")
	params.Set("field_name", "weight")
	params.Set("field_value", weight)

	dreq := strings.NewReader(params.Encode()) // Сформировать строку с данными

	client := &http.Client{} // Сформировать код клиента для выполнения запроса
	req, _ := http.NewRequest(http.MethodPost, apiUrl, dreq)
	// Добавить данные заголовков в запрос
	req.Header.Set("Content-type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "d0f00715-09ad-4808-b7a7-a7208e90bdec")

	resp, err := client.Do(req) // Выполнить запрос
	if err != nil {
		log.Fatalf("Невозможно прочитать код авторизации: %v", err)
	}
	// Отложить выполнения закрытия запроса, до выполения и получения данных
	defer resp.Body.Close()
	//fmt.Printf("%v ", resp.Status)
	body, err := ioutil.ReadAll(resp.Body) // При ошибке получения данных с сервера
	if err != nil {
		log.Println("Ошибка чтения данных ответа: ", err)
		//panic(err)
	}
	//log.Println("\nResponse of server:", string([]byte(body)))
	cw <- string([]byte(body))
}
