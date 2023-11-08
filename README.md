[![Go](https://github.com/blablatov/gosheets/actions/workflows/gosheets-action.yml/badge.svg)](https://github.com/blablatov/gosheets/actions/workflows/gosheets-action.yml)
### Ru

Тестовый модуль для копирования данных `Вес`, `План`, `Факт` из ячеек google-таблицы в программу вебсервиса `KPI`.   
Строка запуска модуля:

    gosheets 1A663PCe8LUilZ-tWbImbj4vlSikymqRBPA62gDVVddw

где:
> `1A663PCe8LUilZ-tWbImbj4vlSikymqRBPA62gDVVddw` -- ID google-таблицы, берется из URL таблицы `Google Sheets`

Для аутентификации модуля на Google-сервисе необходимо получить `credentials.json` и `token.json` файлы.  
#### Алгоритм получения:
>1. Зайти на [APIs & Services](https://console.cloud.google.com/apis/credentials)
>2. Выбрать слева `Credentials`.
>3. Выбрать свой проект.
>4. Нажать `CREATE CREDENTIALS`.
>5. Выбрать способ аутентификации `OAuth 2.0 Client IDs`. Desktop app -- `Name` (имя проекта). Получить `client-secret...json` файл. Переименовать `client-secret...json` в `credentials.json`. Удалить старый `token.json` в директории модуля.
>6. Выбрать слева и настроить `OAuth consent screen` (Publishing status, User type --> `External`, Test users --> `мыло@gmail.com`(ящик нужной учетки).
>7. Запустить наш go-модуль, пройти по ссылке в браузере которую вернет модуль. Ввести полученный код авторизации обратно в комстроку модуля и нажать enter. В директории модуля появится файл `token.json`. 

### En

Test module for copy data of `Weight`, `Plan`, `Fact` from table of Google to programm of web-server в `KPI`.  

We can see use `Google Sheets API`.
String of run module:

    gosheets 1A663PCe8LUilZ-tWbImbj4vlSikymqRBPA62gDVVddw

where is:
> `1A663PCe8LUilZ-tWbImbj4vlSikymqRBPA62gDVVddw` -- ID Google sheet, we got it from URL of a sheet `Google Sheets`

For work to service Google Sheets API we should get `credentials.json` and `token.json` files.
