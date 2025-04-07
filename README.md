# 🗺️ OSM Tail Generator & Importer

Go-сервис для генерации `osm-tail` (вырезки из OpenStreetMap по заданной области) и импорта OSM-данных в систему хранения (например, PostgreSQL + PostGIS).

## 📌 Основные возможности

-   Импорт `/coordinates/` `.osm.pbf` или `.osm.xml` файлов в БД по координатам (bbox)

```JSON
{
    "minLat": 51.05,
    "minLon": 71.30,
    "maxLat": 51.20,
    "maxLon": 71.50
}
```

-   Генерация `/tiles/:z/:x/:y` по координатам (bbox)

### 🛠️ Настройки env (App)

| Variable Name              | Description          | Example Value |
| -------------------------- | -------------------- | ------------- |
| `PORT`                     | Порт                 | 3000          |
| `ENABLE_HEADER_VALIDATION` | Проверка заголовоков | true          |

#### Пример `.env` файла:

```ini
POSTGRES_DB=my_database
POSTGRES_USER=admin
POSTGRES_PASSWORD=securepassword123
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
```

### 🛠️ Настройки env (PostgreSQL)

| Variable Name       | Description                         | Example Value       |
| ------------------- | ----------------------------------- | ------------------- |
| `POSTGRES_DB`       | Название основной базы данных       | `my_database`       |
| `POSTGRES_USER`     | Имя пользователя PostgreSQL         | `admin`             |
| `POSTGRES_PASSWORD` | Пароль пользователя PostgreSQL      | `securepassword123` |
| `POSTGRES_HOST`     | Хост, где запущен PostgreSQL сервер | `localhost` / `db`  |
| `POSTGRES_PORT`     | Порт PostgreSQL                     | `5432`              |

#### Пример `.env` файла:

```ini
POSTGRES_DB=my_database
POSTGRES_USER=admin
POSTGRES_PASSWORD=securepassword123
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
```

## 🚀 Быстрый старт

### 📥 Запуск

```bash
go run main.go
```
