## Device Management Server
This service receives and processes data in the form of `readings` from devices out in the field.

### Getting Started
This service should be runnable as a docker image.
Make sure you have [docker desktop](https://www.docker.com/products/docker-desktop/) downloaded and running, then run the following from the project root:
```shell
docker build -t brightwheel-app .
docker run -p 8080:8080 brightwheel-app
```
This should build the docker image and then run that image on port 8080.

To test whether the service is up:
```shell
curl http://localhost:8080/ping
``` 

### Device Upload Endpoint

**Host Endpoint** `http://localhost:8080/devices`

**Method:** POST

**Description:** Send device readings to be processed and persisted in an in-memory storage pattern.

### Request Headers
| Header | Type | Required | Description   |
| --------- | ---- | ------- | ------------- |
| `Authorization` | string  | Yes     | should be `Basic YWRtaW46cGFzc3dvcmQxIQ==`  |
| `Content-Type` | string  | Yes     | should always be `application/json` |

### Example Payload
```json
{"id": "36d5658a-6908-479e-887e-a949ec199271","readings": [{"timestamp": "2021-09-29T16:08:15+01:00","count": 2},{"timestamp":"2021-09-29T16:09:15+01:00","count": 15}]}
```

### Example cURL Request
``` shell
curl -X POST -u admin:password1! -H "Content-Type: application/json" -d '{"id": "36d5658a-6908-479e-887e-a949ec199271","readings": [{"timestamp": "2021-09-29T16:08:15+01:00","count": 2},{"timestamp":"2021-09-29T16:09:15+01:00","count": 15}]}' http://localhost:8080/device
```

### Response Codes

| Status Code | Description      |
| ----------- | ---------------- |
| 200         | OK               |
| 415         | Unsupported Media Type |
| 500         | Internal Server Error |
| 404         | Record Not Found |

### Get Device's Latest Timestamp Endpoint

**Host Endpoint** `http://localhost:8080/latest_timestamp`

**Method:** GET

**Description:** Get the latest timestamp for a device

### Request Headers
| Header | Type | Required | Description   |
| --------- | ---- | ------- | ------------- |
| `Authorization` | string  | Yes     | should be `Basic YWRtaW46cGFzc3dvcmQxIQ==`  |

### Request Parameters
| Parameter | Type | Required | Description   |
| --------- | ---- | ------- | ------------- |
| `deviceId` | string  | Yes     | ID for the device in question  |

### Example cURL Request
``` shell
curl -u admin:password1! http://localhost:8080/latest_timestamp?deviceId=<deviceIdString>
```

### Response Codes

| Status Code | Description      |
| ----------- | ---------------- |
| 200         | OK               |
| 400         | Bad Request |
| 500         | Internal Server Error |
| 404         | Record Not Found |

### Get Device's Cumulative Count Endpoint

**Host Endpoint** `http://localhost:8080/cumulative_count`

**Method:** GET

**Description:** Get the cumulative count of all the readings for a single device

### Request Headers
| Header | Type | Required | Description   |
| --------- | ---- | ------- | ------------- |
| `Authorization` | string  | Yes     | should be `Basic YWRtaW46cGFzc3dvcmQxIQ==`  |

### Request Parameters
| Parameter | Type | Required | Description   |
| --------- | ---- | ------- | ------------- |
| `deviceId` | string  | Yes     | ID for the device in question  |

### Example cURL Request
``` shell
curl -u admin:password1! http://localhost:8080/latest_timestamp?deviceId=<deviceIdString>
```

### Response Codes

| Status Code | Description      |
| ----------- | ---------------- |
| 200         | OK               |
| 400         | Bad Request |
| 500         | Internal Server Error |
| 404         | Record Not Found |

### Next steps
Next steps I would like to implement for this service:
- build out the test suite for the handlers package
- Use a community logger library (like [zap](https://pkg.go.dev/go.uber.org/zap)) to better implement logging in this project
- improve the project structure by separating the api/business/storage layers better
