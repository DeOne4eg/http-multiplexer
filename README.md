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
+ ✅ tests

## HOWTO
+ run with `make run`
+ build with `make build`
+ test with `make test`