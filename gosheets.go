// Main package for copying data of Weight, Plan, Fact from google-table cells to the KPI web service program.
// Основной пакет модуля, для копирования данных Вес, План, Факт из ячеек google-таблицы в программу вебсервиса KPI.
// www.kpi-drive.ru

package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/blablatov/gosheets/cmweipost"
	"github.com/blablatov/gosheets/factpost"
	"github.com/blablatov/gosheets/getmo"
	"github.com/blablatov/gosheets/planpost"
	"github.com/blablatov/gosheets/weipost"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

type SheetData struct {
	mu                      sync.Mutex
	indicator_to_mo_fact_id string `json:"indicator_to_mo_fact_id"`
	indicator_to_mo_id      string `json:"indicator_to_mo_id"`
	weight                  string `json:"weight"`
	plan_default            string `json:"plan_default"`
	value                   string `json:"value"`
}

// Anonymous field. Composition for secure access to types and methods of the `cmweipost` package.
// Анонимное поле. Композиция для безопасного доступа к типам и методам пакета `cmweipost`.
type embtypes struct {
	cmweipost.DataPost
}

// Getting token, saving token, returning generated client.
// Получение токена, сохранение токена, возвращение сгенерированного клиента.
func getClient(config *oauth2.Config) *http.Client {
	// The token.json file stores access tokens and  updates it,
	// created automatically on first login.
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

// Request token via web, return got token. // Запросить токен через веб, вернуть извлеченный токен.
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

// Get token from local file. Получить токен из локального файла.
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

// Save path to token file. Сохранить путь к файлу токена.
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
	// Get and check id-tables, it's key for of command line run the module.
	// Получить и проверить id-таблицы, это ключ коммандной строки запуска модуля.
	var spreadsheetId, sep string
	for i := 1; i < len(os.Args); i++ {
		spreadsheetId += sep + os.Args[i]
		sep = " "
	}
	fmt.Println("\nКлючи коммандной строки: ", spreadsheetId)
	//spreadsheetId = "1A663PCe8LUilZ-tWbImbj4vlSikymqRBPA62gDVVddw" //id-таблицы для отладки
	arsl := len("1A663PCe8LUilZ-tWbImbj4vlSikymqRBPA62gDVVddw")
	if len(spreadsheetId) > arsl || len(spreadsheetId) < arsl {
		log.Fatalf("\nВведите ID таблицы\n")
	}
	sp := &SheetData{}
	sp.indicator_to_mo_fact_id = "0"

	cw := make(chan string) // Channels for functions of post-requests.
	cp := make(chan string) // Каналы для функций post-запросов.
	cf := make(chan string)
	cg := make(chan string)
	ctx := context.Background()
	b, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Не удалось прочитать секретный файл клиента: %v", err)
	}

	// If one change these areas, one must delete previously saved token.json file.
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

	// Getting `mo_id` via request to the `get_mo` of server method.
	// Run getmo function in goroutine.
	// Получение `mo_id` через запрос к методу `get_mo` сервера.
	// Getting data URL from config file. Получение URL из конфига.
	cnurl := make(chan string)
	go func() {
		ptrn := regexp.MustCompile(`https://testdb.kpi-drive.ru/_api/mo/get_mo`)
		gurl := ReadUrlServ()
		gstr := ptrn.FindAllString(gurl, -1)
		url := fmt.Sprint(gstr)
		cnurl <- url
	}()
	apiUrl := <-cnurl

	go getmo.GetMoid(apiUrl, cg)
	sp.indicator_to_mo_id = <-cg // Get data from channel of gorutine. Получение данных из канала горутины.
	//sp.indicator_to_mo_id = "994" //For debug, indicator mo_id. Для отладки.
	fmt.Println("\n Результат запроса get: ", sp.indicator_to_mo_id)
	secs := time.Since(start).Seconds()
	fmt.Printf(" %.2fs время выполнения запроса\n", secs)

	p := recover()
	if len(sp.indicator_to_mo_id) == 0 {
		log.Println("Panic, internal error, data of answed not got. recover()")
		panic(p)
	}

	// Get column `mo_id` from table. Получение столбца с `mo_id` таблицы.
	columnRange := "5. RKPI-Карта!A4:A4"
	resm, err := srv.Spreadsheets.Values.Get(spreadsheetId, columnRange).Do()
	if err != nil {
		log.Fatalf("Не удалось получить данные mo_id из столбца A: %v", err)
	}
	if len(resm.Values) == 0 {
		fmt.Println("Данные столбца не найдены.")
	} else {
		// Get row data for each mo_id of table in loop.
		// Получить в цикле данные строки для каждого mo_id таблицы.
		for _, rowm := range resm.Values {
			if rowm != nil {
				mid := fmt.Sprint(rowm)
				// Компиляция строкового представления регулярного выражения единожды, а не многократно.
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

				// Get data from sheet of table. Получить данные с листа таблицы.
				resp, err := srv.Spreadsheets.Values.Get(spreadsheetId, readRange).Do()
				if err != nil {
					log.Fatalf("Не удалось получить данные с листа таблицы: %v", err)
				}

				if len(resp.Values) == 0 {
					fmt.Println("Данные не найдены.")
				} else {
					fmt.Println(" \nДанные таблицы RKPI-Карта:")
					for _, row := range resp.Values {
						// Set columns by indexes `row[x]`. Вывести столбцы,по индексам `row[x]`
						fmt.Printf(" Вес: %s\n План: %s\n Факт: %s\n", row[0], row[1], row[2])
					}
					spars := fmt.Sprint(resp.Values)  // Get a range data of one string. Диапазон данных одной строкой.
					pars := strings.Trim(spars, "[]") // Del. Удалить []

					////////////////////////////////////////////////////////////////////
					// Parsing rows of data from a sheet of table.
					// Парсинг строк данных с листа таблицы.
					//sd := &DataType{}
					counts := make(map[string]int)
					//Slice for save all keys from mapping. Срез для хранения всех ключей мапы.
					datakeys := make([]string, 0, len(counts))
					for _, line := range strings.Split(string(pars), ":") {
						counts[line]++
						log.Println(line) // Checks data of mapping.
					}
					// Sorts keys to list and to sets values in order.
					// Сортировка ключей для перечисления и присваивания значений по порядку.
					for countkeys := range counts {
						datakeys = append(datakeys, countkeys)
					}
					sort.Strings(datakeys)
					for _, countkeys := range datakeys {
						fmt.Printf("\nCountkeys: %v\nCounts: %v\n", countkeys, counts[countkeys])
						if countkeys != "" {
							if sp.weight == "" {
								sp.mu.Lock()
								sp.weight = countkeys
								sp.mu.Unlock()
							} else {
								if sp.plan_default == "" {
									sp.mu.Lock()
									sp.plan_default = countkeys
									sp.mu.Unlock()
								} else {
									sp.mu.Lock()
									sp.value = countkeys
									sp.mu.Unlock()
								}
							}
						}
					}
					// Output for test. Тестовый вывод данных.
					log.Println("Weight: ", sp.weight)
					log.Println("Plan_default: ", sp.plan_default)
					log.Println("Value : ", sp.value)

					fmt.Println("  Данные для передачи на testdb.kpi-drive.ru\n  weight:", sp.weight, "\n  plan_default:", sp.plan_default,
						"\n  value:", sp.value, "\n  mo_fact_id:", sp.indicator_to_mo_fact_id, "\n  mo_id:", sp.indicator_to_mo_id)

					var wg sync.WaitGroup
					wg.Add(3) // Counter of goroutines. Значение счетчика горутин.
					// Run function to get `Weight` via goroutine. Функция получения `Веса` в горутине.
					go weipost.PostTestdb(wg, sp.weight, sp.indicator_to_mo_fact_id, sp.indicator_to_mo_id, cw)
					log.Println("\n\nРезультат Вес:", <-cw) // Get data fro gorutine. Получить данные из канала горутины.
					secs := time.Since(start).Seconds()
					fmt.Printf("%.2fs время выполнения запроса\n", secs)

					// Run function to get `Plan` via goroutine. Функция получения `Плана` в горутине.
					go planpost.PlanTestdb(wg, sp.plan_default, sp.indicator_to_mo_fact_id, sp.indicator_to_mo_id, cp)
					log.Println("\n\nРезультат План:", <-cp) // Получение данных из канала горутины
					secs2 := time.Since(start).Seconds()
					fmt.Printf("%.2fs время выполнения запроса\n", secs2)

					// Run function to get `Fact` via goroutine. Функция получения `Факта` в горутине.
					go factpost.FactTestdb(wg, sp.value, sp.indicator_to_mo_fact_id, sp.indicator_to_mo_id, cf)
					log.Println("\n\nРезультат Факт:", <-cf) // Получение данных из канала горутины
					secs3 := time.Since(start).Seconds()
					fmt.Printf("%.2fs время выполнения запроса\n", secs3)
					// Wait of counter. Ожидание счетчика.
					go func() {
						wg.Wait()
						close(cw)
						close(cp)
						close(cf)
					}()

					// Option one.
					// Calling an interface method via struct embedding.
					// Вызов метода ReadCoils, через встроенную структуру.
					start4 := time.Now()
					var w embtypes

					// Formating data of structure. Заполнение структуры.
					w.DataPost.Indicator_to_mo_fact_id = sp.indicator_to_mo_fact_id
					w.DataPost.Indicator_to_mo_id = sp.indicator_to_mo_id
					w.DataPost.Weight = sp.weight

					// Getting data URL from config file. Получение URL из конфига.
					cnurlw := make(chan string)
					go func() {
						ptrn := regexp.MustCompile(`https://testdb.kpi-drive.ru/_api/indicators/save_indicator_instance_field`)
						gurl := ReadUrlServ()
						gstr := ptrn.FindAllString(gurl, -1)
						url := fmt.Sprint(gstr)
						cnurlw <- url
					}()
					w.ApiUrlw = <-cnurlw

					res, err := embtypes.WPostSend(w)
					if err != nil {
						log.Fatalf("Error of method: %v", err)
					}
					log.Println("Result of request an interface method via struct embedding: ", res)
					secs4 := time.Since(start4).Seconds()
					fmt.Printf("%.2fs Request execution time via method WPostSend of struct embedding\n", secs4)
				}
			}
		}
	}
}

// Func reads data from file the ./gosheets.conf
func ReadUrlServ() string {
	var url string
	f, err := os.Open("gosheets.conf")
	if err != nil {
		log.Fatalf("Error open a conf-file gosheets: %v", err)
	}
	defer f.Close()
	input := bufio.NewScanner(f)
	for input.Scan() {
		url = input.Text()
	}
	return url
}
