# 🗺️ OSM Tail Generator & Importer

A Go service for generating osm-tail (a subset of OpenStreetMap data for a given area) and importing OSM data into a storage system (e.g., PostgreSQL + PostGIS).

## 📌 Key Features

-   Import .osm.pbf or .osm.xml files into the database by coordinates (bbox)

```JSON
{
    "minLat": 51.05,
    "minLon": 71.30,
    "maxLat": 51.20,
    "maxLon": 71.50
}
```

-   Generate /tiles/:z/:x/:y based on coordinates (bbox)

### 🛠️ Настройки env (App)

| Variable Name              | Description       | Example Value |
| -------------------------- | ----------------- | ------------- |
| `PORT`                     | Port              | 3000          |
| `ENABLE_HEADER_VALIDATION` | Header validation | true          |

#### Пример `.env` файла:

```ini
POSTGRES_DB=my_database
POSTGRES_USER=admin
POSTGRES_PASSWORD=securepassword123
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
```

### 🛠️ Настройки env (PostgreSQL)

| Variable Name       | Description              | Example Value       |
| ------------------- | ------------------------ | ------------------- |
| `POSTGRES_DB`       | Вatabase name            | `my_database`       |
| `POSTGRES_USER`     | PostgreSQL username      | `admin`             |
| `POSTGRES_PASSWORD` | PostgreSQL user password | `securepassword123` |
| `POSTGRES_HOST`     | PostgreSQL server host   | `localhost` / `db`  |
| `POSTGRES_PORT`     | PostgreSQL port          | `5432`              |

#### Example `.env` file:

```ini
POSTGRES_DB=my_database
POSTGRES_USER=admin
POSTGRES_PASSWORD=securepassword123
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
```

## 🚀 Quick Start

### 📥 Running the Service

```bash
go run main.go
```
