# Flower Shop

![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white)
![Postgres](https://img.shields.io/badge/postgres-%23316192.svg?style=for-the-badge&logo=postgresql&logoColor=white)
![Docker](https://img.shields.io/badge/docker-%230db7ed.svg?style=for-the-badge&logo=docker&logoColor=white)
![Redis](https://img.shields.io/badge/redis-%23DD0031.svg?style=for-the-badge&logo=redis&logoColor=white)
![Apache Kafka](https://img.shields.io/badge/Apache%20Kafka-000?style=for-the-badge&logo=apachekafka)
![Postman](https://img.shields.io/badge/Postman-FF6C37?style=for-the-badge&logo=postman&logoColor=white)

## Описание

Проект представляет собой разработанный на Golang микросервисный RESTful API backend для цветочного магазина с использованием Kafka, Redis, PostgreSQL, Сhi и развёрнутый на Docker. Весь процесс логируется.
В проекте использовал горутины и каналы для параллельной обработки задач и Postman для тестирования и отладки API.

## Архитектура приложения

![Схема архитектуры](https://github.com/RoiDuNord/flower_shop/blob/master/architecture-diagram.svg)

## Предварительные требования

- **Docker**: Убедитесь, что Docker установлен и работает на вашем компьютере. Команда выведет установленную версию Docker. Если Docker не установлен, вы получите сообщение об ошибке

```
docker --version
```

- Проверка работы Docker

```
docker info
```

- Если Docker работает, вы увидите информацию о конфигурации и состоянии вашего Docker-демона. Если он не запущен, вы получите сообщение об ошибке

## Алгоритм сборки Docker-образа и запуска приложения

1. Клонируйте проект на ваш компьютер из Github с помощью команды

```
git clone https://github.com/RoiDuNord/flower_shop.git
```

2. Перейдите в директорию сервера и запустите приложение с помощью Docker Compose

```
cd server_side
docker-compose up -d --build
```

3. Проверьте работу приложения с помощью Postman или другого HTTP-клиента для отправки запросов к API

- Создайте заказы
```
POST http://localhost:8080/orders/create
```
- Просмотр готовых заказов
```
GET http://localhost:8080/orders/get
```