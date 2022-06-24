package main

import (
	"context"
	"encoding/json"
	"factpost"
	"fmt"
	"getmo"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"planpost"
	"regexp"
	"strings"
	"time"
	"weipost"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

type SheetData struct {
	indicator_to_mo_fact_id string `json:"indicator_to_mo_fact_id"`
	indicator_to_mo_id      string `json:"indicator_to_mo_id"`
	weight                  string `json:"weight"`
	plan_default            string `json:"plan_default"`
	value                   string `json:"value"`
}

// Получение токена, сохранение токена, возвращение сгенерированного клиента.
func getClient(config *oauth2.Config) *http.Client {
	// В файле token.json хранятся токены доступа и обновления пользовательского доступа,
	// созданные автоматически при первой авторизации.
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Запросить токен через веб, вернуть извлеченный токен.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Перейдите по ссылке в браузере и введите  "+
		"код авторизации: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Невозможно прочитать код авторизации: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Не удалось получить токен из веб-запроса: %v", err)
	}
	return tok
}

// Получить токен из локального файла
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Сохранить путь к файлу токена.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Сохранение файла учетных данных в: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Не удалось кэшировать токен аутентификаци : %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func main() {
	start := time.Now()
	// Получить и проверить id-таблицы, ключ коммандной строки запуска модуля
	var spreadsheetId, sep string
	for i := 1; i < len(os.Args); i++ {
		spreadsheetId += sep + os.Args[i]
		sep = " "
	}
	fmt.Println("\nКлючи коммандной строки: ", spreadsheetId)
	//spreadsheetId = "1A663PCe8LUilZ-tWbImbj4vlSikymqRBPA62gDVVddw" //id-таблицы для отладки
	arsl := len("1A663PCe8LUilZ-tWbImbj4vlSikymqRBPA62gDVVddw")
	if len(spreadsheetId) > arsl || len(spreadsheetId) < arsl {
		fmt.Fprintf(os.Stderr, "\nВведите ID таблицы\n")
		os.Exit(1)
	}
	sp := &SheetData{}
	sp.indicator_to_mo_fact_id = "0"

	cw := make(chan string) // каналы для функций post-запросов
	cp := make(chan string)
	cf := make(chan string)
	cg := make(chan string)
	ctx := context.Background()
	b, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Не удалось прочитать секретный файл клиента: %v", err)
	}

	// При изменении этих областей надо удалить ранее сохраненный файл token.json.
	config, err := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/spreadsheets.readonly")
	if err != nil {
		log.Fatalf("Не удалось разобрать секретный файл клиента для конфигурации: %v", err)
	}
	client := getClient(config)

	srv, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Не удалось получить Таблицы клиента: %v", err)
	}

	// Получение mo_id через запрос к методу get_mo сервера
	// Запуск функции getmo в горутине
	go getmo.GetMoid(cg)
	sp.indicator_to_mo_id = <-cg // Получение данных из канала горутины
	//sp.indicator_to_mo_id = "994" //mo_id для отладки
	fmt.Println("\n Результат запроса get: ", sp.indicator_to_mo_id)
	secs := time.Since(start).Seconds()
	fmt.Printf(" %.2fs время выполнения запроса\n", secs)
	if len(sp.indicator_to_mo_id) == 0 {
		fmt.Println(" Данные mo_id не найдены.")
	}

	// Получение столбца с mo_id таблицы
	columnRange := "5. RKPI-Карта!A4:A4"
	resm, err := srv.Spreadsheets.Values.Get(spreadsheetId, columnRange).Do()
	if err != nil {
		log.Fatalf("Не удалось получить данные mo_id из столбца A: %v", err)
	}
	if len(resm.Values) == 0 {
		fmt.Println("Данные столбца не найдены.")
	} else {
		// Получаем в цикле данные строки для каждого mo_id таблицы
		for _, rowm := range resm.Values {
			if rowm != nil {
				mid := fmt.Sprint(rowm)
				// Делать компиляцию строкового представления регулярного выражения единожды, а не многократно.
				pattern := regexp.MustCompile(`\D+`)
				sp.indicator_to_mo_id = pattern.ReplaceAllString(mid, "")
				fmt.Println("Данные столбца mo_id: ", sp.indicator_to_mo_id)
				var readRange string
				switch sp.indicator_to_mo_id {
				case "994":
					readRange = "5. RKPI-Карта!J4:L4"
				case "199":
					readRange = "5. RKPI-Карта!J5:L5"
				default:
					fmt.Println("Данные строки не найдены!")
				}

				// Получить данные с листа таблицы
				resp, err := srv.Spreadsheets.Values.Get(spreadsheetId, readRange).Do()
				if err != nil {
					log.Fatalf("Не удалось получить данные с листа таблицы: %v", err)
				}

				if len(resp.Values) == 0 {
					fmt.Println("Данные не найдены.")
				} else {
					fmt.Println(" \nДанные таблицы RKPI-Карта:")
					for _, row := range resp.Values {
						// Вывести столбцы, соответствующие индексам row[x]
						fmt.Printf(" Вес: %s\n План: %s\n Факт: %s\n", row[0], row[1], row[2])
					}
					spars := fmt.Sprint(resp.Values)  // Получить диапазон данных одной строкой
					pars := strings.Trim(spars, "[]") // Удалить скобки

					// Получение данных из среза для функций Post
					i := strings.Index(pars, " ") // Получить индекс первого символа " "
					sp.weight = pars[:i]          // Получить срез до 1-го символа " "
					substr := pars[i+1:]          // Получить срез от 1-го символа " "

					t := strings.Index(substr, " ") // Получить индекс первого символа " " в подстроке substr
					pland := substr[:t]             // Получить срез до 1-го символа " " подстроки substr
					valf := substr[t+1:]            // Получить срез от 1-го символа " " подстроки substr/в substr2

					sp.plan_default = pattern.ReplaceAllString(pland, "")
					sp.value = pattern.ReplaceAllString(valf, "")

					fmt.Println("  Данные для передачи на testdb.kpi-drive.ru\n  weight:", sp.weight, "\n  plan_default:", sp.plan_default,
						"\n  value:", sp.value, "\n  mo_fact_id:", sp.indicator_to_mo_fact_id, "\n  mo_id:", sp.indicator_to_mo_id)

					// Запуск функции получения Веса в горутине
					go weipost.PostTestdb(sp.weight, sp.indicator_to_mo_fact_id, sp.indicator_to_mo_id, cw)
					log.Println("\n\nРезультат Вес:", <-cw) // Получение данных из канала горутины
					secs := time.Since(start).Seconds()
					fmt.Printf("%.2fs время выполнения запроса\n", secs)

					// Запуск функции получения Плана в горутине
					go planpost.PlanTestdb(sp.plan_default, sp.indicator_to_mo_fact_id, sp.indicator_to_mo_id, cp)
					log.Println("\n\nРезультат План:", <-cp) // Получение данных из канала горутины
					secs2 := time.Since(start).Seconds()
					fmt.Printf("%.2fs время выполнения запроса\n", secs2)

					// Запуск функции получения Факта в горутине
					go factpost.FactTestdb(sp.value, sp.indicator_to_mo_fact_id, sp.indicator_to_mo_id, cf)
					log.Println("\n\nРезультат Факт:", <-cf) // Получение данных из канала горутины
					secs3 := time.Since(start).Seconds()
					fmt.Printf("%.2fs время выполнения запроса\n", secs3)
				}
			}
		}
	}
}
