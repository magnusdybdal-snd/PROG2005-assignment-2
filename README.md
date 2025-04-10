# Assignment2

## Table of contents
- [Description](#description)
- [Installation](#installation)
- [Usage](#usage)
- [Contributing](#contributing)
- [Testing](#testing)

## Description

This project was created to show country information based on registrations in a dashboard. The dashboard is created using the Firebase console and the API is implemented using Go. The API is then used to create a webhook that is invoked when a new registration, deletion, change, invoke or not reachable event is registered.

## Installation
### Ubuntu/Debian
```bash
sudo apt install golang-go
git clone https://git.gvk.idi.ntnu.no/course/prog2005/prog2005-2025-workspace/olemgl/assignment2.git
cd assignment2
go mod download
go run main.go
```
### Windows
https://go.dev/doc/install
```bash
git clone https://git.gvk.idi.ntnu.no/course/prog2005/prog2005-2025-workspace/olemgl/assignment2.git
cd assignment2
go mod download
go run main.go
```
### API Key
You will need a firebase API key to run the project. In order to use the API kek with the docker-compose file, create a folder caller api-keys and add a file called serviceAccountKey.json.

## Usage
This API offers the following functionality:
### Registration endpoint
### Register a new dashboard
Make a POST request to the /dashboards/v1/registrations/ endpoint with the following JSON schema (all fields are required):
```json
{
    "country": "Norway",
    "isoCode": "NO",
    "features": {
        "capital": true,
        "coordinates": true,
        "population": true,
        "area": true,
        "targetCurrencies": ["USD", "EUR", "SEK"]
    }
}
```
It will return a 201 status code if the dashboard was registered successfully.
### Get all dashboards
Make a GET request to the /dashboards/v1/registrations/ endpoint. The response will have all dashboards registered. It will return a 200 status code if the dashboards were found.

### Get a specific dashboard
Make a GET request to the /dashboards/v1/registrations/{id} endpoint. The response will have the dashboard with a specified id. It will return a 200 status code if the dashboard was found.

### Delete a specific dashboard
Make a DELETE request to the /dashboards/v1/registrations/{id} endpoint. The dashboard with a specified id will be deleted. It will return a 204 status code if the dashboard was deleted successfully. This will send a notification to all webhooks registered for the dashboard and or the country.

### Update a specific dashboard
Make a PUT request to the /dashboards/v1/registrations/{id} endpoint with the following JSON schema (all fields are required):
```json
{
    "country": "Norway",
    "isoCode": "NO",
    "features": {
        "capital": true,
        "coordinates": true,
        "population": true,
        "area": true,
        "targetCurrencies": ["USD", "EUR", "SEK"]
    }
}
```
It will return a 200 status code if the dashboard was updated successfully.
This will send a notification to all webhooks registered for the dashboard and or the country.
### Dashboards endpoint
Make a GET request to the /dashboards/v1/dashboards/{id} endpoint. The response will have the dashboard with a specified id. It will return a 200 status code if the dashboard was found.
### Notification endpoint
This endpoint gives you controll over webhooks in the API. The events for notifications are:
- REGISTER - when a new dashboard is registered
- CHANGE - when a dashboard is updated
- DELETE - when a dashboard is deleted
- INVOKE - when a dashboard is invoked
- NOT_REACHABLE - when an external API is not reachable
### Register a new webhook
Make a POST request to the /dashboards/v1/notifications/ endpoint with the following JSON schema (all fields are required):
```json
{
    "country": "Norway",
    "event": "INVOKE",
    "url": "https://example.com/webhook"
}
```
It will return a 201 status code if the webhook was registered successfully.
### Get all webhooks
Make a GET request to the /dashboards/v1/notifications/ endpoint. The response will have all webhooks registered. It will return a 200 status code if the webhooks were found.
### Get a specific webhook
Make a GET request to the /dashboards/v1/notifications/{id} endpoint. The response will have the webhook with a specified id. It will return a 200 status code if the webhook was found.
### Delete a specific webhook
Make a DELETE request to the /dashboards/v1/notifications/{id} endpoint. The webhook with a specified id will be deleted. It will return a 204 status code if the webhook was deleted successfully.
### Update a specific webhook
Make a PATCH request to the /dashboards/v1/notifications/{id} endpoint with the following JSON schema to update one or more fields:
```json
{
    "country": "Norway",
    "event": "INVOKE",
    "url": "https://example.com/webhook"
}
```
It will return a 204 status code if the webhook was updated successfully.
### Invoke of a webhook
It will make a POST request to the webhook URL specified in the request. It will send a JSON payload with the following schema:
```json
{
    "webhookId": "webhook-id",
    "country": "",
    "event": "",
    "time": "YYYYMMDD HH:MM"
}  
```
It will log the http status code of the response and return a 200 status code if the webhook was invoked successfully.

## Contributing
The endpoints were created by: <br>
Registration endpoint/handler - Marius Eilertsen (mailert) <br>
Dashboards endpoint/handler - Magnus Dybdal (magndy) <br>
Notification endpoint/handler/Webhook invokation - Ole Marius Glomsrud (olemgl) <br>
Status endpoint - Ole Marius Glomsrud (olemgl) <br>
All tests were written by the respected authors.

## Testing
All tests are written in the handlers package. The tests are run using the go test command. To run the tests, run the following command in the root directory of the project:
```bash
cd handlers
go test
```
