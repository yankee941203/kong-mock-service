# KONG-MOCK-SERVICE

## Quick Start

```
1. Clone this repo
2. CD to kong-mock-service/ path
3. go run main.go
4. Use curl, have fun
```
### Tips

- The mock data will be inserted into db when serive start to running, if you dont want repetitive insertion of data, please clean-up db file under `sources/`

## Desgin

### APIs
Please Check the OpenAPI v3.0 doc in [here](openapi.yml)

#### GET /api/v1/servcies
Return a list of services info, support filtering, sorting, pagination

Example Request: 
```
curl "0.0.0.0:8000/api/v1/services?page=2&per_page=10&sort_by=Name" | jq`
```
Example Response: 
```
{
  "meta": {
    "page": 2,
    "per_page": 10,
    "total_pages": 12,
    "total_items": 120,
    "prev_page": "/api/v1/services?page=1&per_page=10&sort_by=Name",
    "next_page": "/api/v1/services?page=3&per_page=10&sort_by=Name",
    "last_page": "/api/v1/services?page=12&per_page=10&sort_by=Name"
  },
  "data": [
    {
      "id": 101,
      "name": "kong-test-10",
      "description": "a test service",
      "available_version": [
        "pokemon/pika:v1",
        "pokemon/pika:v2"
      ]
    },
    ...
}
```

#### GET /api/v1/servcie/{id}
Return a particular service, including the availabel version of its images.

Example Request:
```
curl "0.0.0.0:8000/api/v1/service/11" | jq
```
Example Response:
```
{
  "id": 11,
  "name": "kong-test-10",
  "description": "a test service",
  "available_version": [
    "pokemon/pika:v1",
    "pokemon/pika:v2"
  ]
}
```

### Project Info

This service is mainly based on gin and gorm

This repo struct is designed based on MVC, the definations of API are in `internal/router`, the logic of controller layer is in `internal/controller`, and the data models are defined in `internal/model`. The CRUD related operations are implemented in `internal/repository`.

The dependencies between the logic layers are all defined by interfaces and implementations, making it easier for future expansion or testing.

The many-to-many relationship between Service and Image is implemented by two junction tables.

## TODO

Due to time and resource constraints, there are still certain parts of the work that have not been completed, or there is room for improvement.
- Unittest&Test Plan, need to add UT cases to improve quality
- Logic of filtering/sorting/pagination in repository layer, the filtering/sorting/pagination logic are implemented in controller layer, in this way, performance issues may arise as the data scale continues to grow in the future. Implementing this feature in the repository layer can significantly improve performance.
