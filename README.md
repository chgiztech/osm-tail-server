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

## 🚀 Быстрый старт

### 📥 Запуск

```bash
go run main.go
```
