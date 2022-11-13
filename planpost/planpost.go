package planpost

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

func PlanTestdb(wg sync.WaitGroup, plan_default, indicator_to_mo_fact_id, indicator_to_mo_id string, cp chan string) {

	//Запрос для добавления значения плана методом сервера
	apiUrl := "https://testdb.kpi-drive.ru/_api/facts/save_fact"
	// Добавить данные в запрос
	params := url.Values{"period_start": {"2022-02-01"}}
	params.Set("period_end", "2022-02-28")
	params.Set("period_key", "month")
	params.Set("indicator_to_mo_fact_id", indicator_to_mo_fact_id)
	params.Set("auth_user_id", "4")
	params.Set("value", plan_default)
	params.Set("fact_time", "2022-02-08")
	params.Set("indicator_to_mo_id", indicator_to_mo_id)
	params.Set("comment", "Запрос к методу save_fact, для сохранения Плана, Blablatov")
	params.Set("is_plan", "1")
	params.Set("request_id", "")

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
	}
	//log.Println("\nResponse of server:", string([]byte(body)))
	cp <- string([]byte(body))
}
