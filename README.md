# fibonacci_service

Сервис для вычисления среза последовательности Фибоначчи с возможностью подключения
с использованием технологий gRPC и HTTP

## How to start

1. Загрузить исходники

    git clone https://github.com/nik-zaitsev/fibonacci_service.git

2. Перейти в корневую директорию проекта и запустить

    cd fibonacci_service

    go run .

3. В консоли для каждого сервиса ввести номер порта

## How to use
### gRPC

1. Сгенерировать код с использованием утилиты protoc (protobuf схема: schema/fibonacci.proto)

2. Использовать gRPC клиент github.com/ktr0731/evans

### HTTP

    curl http://IP_ADDRESS:PORT?from=FROM\&to=TO 

где IP_ADDRESS - IP адрес, на котором запущен сервис; PORT - порт;
    FROM - левая граница среза; TO - правая граница среза
