package getmo

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"regexp"
	"strings"
)

func GetMoid(apiUrl string, cg chan string) {
	//TODO config
	params := url.Values{"period_start": {"2022-02-01"}}
	params.Set("period_end", "2022-02-28")
	params.Set("period_key", "month")
	//params.Set("requested_mo_id", "0")
	params.Set("auth_user_id", "4")

	greq := strings.NewReader(params.Encode())

	// Create a new request using http
	req, err := http.NewRequest(http.MethodPost, apiUrl, greq)
	// add authorization header to the req
	req.Header.Set("Content-type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "d0f00715-09ad-4808-b7a7-a7208e90bdec")

	fmt.Println("Данные для запроса get_mo_indicators\n  period_start:", "2022-02-28", "\n  period_end:", "2022-02-28",
		"\n  period_key:", "month", "\n  auth_user_id:", "4")

	// Send req using http Client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error on response.\n[ERROR] -", err)
	}
	defer resp.Body.Close()
	fmt.Printf(" %v ", resp.Status)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error while reading the response bytes:", err)
	}
	//log.Println(string([]byte(body)))

	moid := regexp.MustCompile(`"mo_id":..`)
	smoid := string([]byte(body))
	fmoid := moid.FindAllString(smoid, -1)

	pmoid := fmt.Sprint(fmoid)
	// Делать компиляцию строкового представления регулярного выражения единожды, а не многократно.
	pattern := regexp.MustCompile(`[^mo_id:1234567890, ]`)
	gmoid := pattern.ReplaceAllString(pmoid, "")
	//fmt.Printf("\n Номера mo_id сотрудников: %q\n\n", gmoid)
	cg <- string(gmoid)
}
