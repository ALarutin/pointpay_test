# pointpay_test

## Глобальные переменные для поднятия сервиса:

+ account
    + LOG_LEVEL
    + MONGO_URI (пример mongodb://0.0.0.0:27017/account)
    + MONGO_COLLECTION_NAME
    + GRPC_SERVER_PORT
+ bank
    + LOG_LEVEL
    + HTTP_SERVER_PORT
    + GRPC_SERVER_PORT

## Что я бы добавил, если бы писал в прод
+ метрики Prometheus
+ healthchecks
+ трассировку Jeager
+ RED мидлвар для http сервера
+ тесты
+ больше логирования
+ механизмы проверки соединений между grpc client/server
+ sentry
+ возможность включения pprof
+ возможность конфигурирования сервисов за счёт иных способов помимо переменных окружения
+ комментарий к коду