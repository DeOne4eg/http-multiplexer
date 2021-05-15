![Github CI/CD](https://img.shields.io/github/workflow/status/DeOne4eg/http-multiplexer/Go)
![Go Report](https://goreportcard.com/badge/github.com/DeOne4eg/http-multiplexer)
![Repository Top Language](https://img.shields.io/github/languages/top/DeOne4eg/http-multiplexer)
![Scrutinizer Code Quality](https://img.shields.io/scrutinizer/quality/g/DeOne4eg/http-multiplexer/master)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/DeOne4eg/http-multiplexer)
![Codacy Grade](https://img.shields.io/codacy/grade/c9467ed47e064b1981e53862d0286d65)
![Github Repository Size](https://img.shields.io/github/repo-size/DeOne4eg/http-multiplexer)
![Github Open Issues](https://img.shields.io/github/issues/DeOne4eg/http-multiplexer)
![Lines of code](https://img.shields.io/tokei/lines/github/DeOne4eg/http-multiplexer)
![License](https://img.shields.io/badge/license-MIT-green)
![GitHub last commit](https://img.shields.io/github/last-commit/DeOne4eg/http-multiplexer)
![GitHub contributors](https://img.shields.io/github/contributors/DeOne4eg/http-multiplexer)

<img align="right" width="300px" src="./images/go.png">

# HTTP Multiplexer

## Task description

Приложение представляет собой HTTP-сервер с одним хендлером. Хендлер на вход получает POST-запрос со списком URL в JSON-формате. Сервер запрашивает данные по всем этим URL и возвращает результат клиенту в JSON-формате. Если в процессе обработки хотя бы одного из URL получена ошибка, обработка всего списка прекращается и клиенту возвращается текстовая ошибка.

**Ограничения:**
+ для реализации задачи следует использовать Go 1.13 или выше
+ использовать можно только компоненты стандартной библиотеки Go
+ сервер не принимает запрос если количество URL в нем больше 20
+ сервер не обслуживает больше чем 100 одновременных входящих http-запросов
+ для каждого входящего запроса должно быть не больше 4 одновременных исходящих
+ таймаут на запрос одного URL - 1 секунда
+ обработка запроса может быть отменена клиентом в любой момент, это должно повлечь за собой остановку всех операций связанных с этим запросом
+ сервис должен поддерживать 'graceful shutdown'

## Solution notes
+ 🔱 clean architecture
+ 📖 only standard Go components are used
+ 📝 minimum logs
+ ✅ tests

## HOWTO
+ run with `make run`
+ build with `make build`
+ test with `make test`