# Micro-API-Mongo

## Description
This is a inimalistic API written in GO with MongoDB storage

## Usage

```bash
docker run -d --name micro-api -p 8080:8080 -e ADMIN_PASSWORD=admin micro-api
```

## Available endpoints

| URL (relative)  | METHOD | REQUEST          | RESPONSE         | AUTH             | RESPONSE             | COMMENT                          |
| :-------------- | :----- | :--------------- | :--------------- | :--------------- | :------------------- | :------------------------------- |
| /records        | GET    | n/a              | application/json | no               | 200 OK + json list   | Get all records                  |
| /records/$ID    | GET    | n/a              | application/json | no               | 200 OK + json object | Get record with id $ID           |
| /records        | POST   | application/json | n/a              | no               | 200 OK               | Create a new record with payload |
| /records/random | GET    | n/a              | application/json | no               | 200 OK + json list   | Get random record                |
| /admin          | GET    | n/a              | text/plain       | no               | 401 Unauthorized     | Try to access w/o authentication |
| /admin          | GET    | n/a              | text/html        | yes (basic-auth) | 200 OK + html        | Access w/ correct authentication |


## JSON payload format

`GET /records` - Response:
```json
[
	{
  		"name": "Record 1",
  		"desc": "Description of the first record",
		"id": "1647195672916783205"
	},
	...
]
```

`GET /records/1647195672916783205` - Response:
```json
{
	"name": "Record 1",
	"desc": "Description of the first record",
	"id": "1647195672916783205"
}
```
`GET /records/random` - Response:
```json
{
	"name": "Record 1",
	"desc": "Description of the first record",
	"id": "1647195672916783205"
}
```

`POST /records/$ID` - Request:
```json
{
	"name": "Record 1",
	"desc": "Description of the first record"
}
```


