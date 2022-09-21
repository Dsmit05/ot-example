# ot-example - presentation of OpenTelemetry work

Данный проект показывает работу с openTelemetry.

Все зависимые приложения такие как signoz, кафка и др. находятся в папке env-apps
Для запуска, пожалуйста используйте подготовленные docker-compose файлы.

Проект выполнен в двух вариациях зависимых от использованного брокера сообщений.

![диаграмма компонентов](http://www.plantuml.com/plantuml/proxy?cache=no&fmt=svg&src=https://raw.githubusercontent.com/Dsmit05/ot-example/master/components.plantuml)

### Пример запуска:

1. Поднимите инструмент мониторинга SigNoz:
2. Поднимите postgres
3. RMQ или Kafka
4. После можете запустить го сервисы:
5. И перейти к свагеру `http://localhost:8080/swagger/index.html` и SigNoz UI `http://localhost:3301`

<details>
  <summary>endpoint this project</summary>

1. http://localhost:8080/swagger/index.html - service main swagger
2. http://localhost:9081/swagger/ - service read swagger
3. http://localhost:9080 - service read grpc route
4. http://localhost:3301 - SigNoz ui
5. jdbc:postgresql://localhost:5432/example - database
6. localhost:4317 - Collector URL
7. localhost:9092 - Kafka
8. localhost:5672 - RMQ
9. http://localhost:15672 - RMQ UI
</details>


<details>
  <summary>Справочная информация</summary>

1. https://opentelemetry.io/ - OpenTelemetry is a collection of tools, APIs, and SDKs
2. https://habr.com/ru/company/ru_mts/blog/537892/ - OpenTelemetry на практике
3. https://signoz.io/blog/monitoring-your-go-application-with-signoz/ - Golang application performance with SigNoz
</details>
