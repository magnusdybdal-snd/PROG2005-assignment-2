package utils

// Paths
const ROOT_PATH = "/"
const REGISTRATION_PATH = "/dashboard/v1/registrations/"
const DASHBOARD_PATH = "/dashboard/v1/dashboards/"
const NOTIFICATION_PATH = "/dashboard/v1/notifications/"
const STATUS_PATH = "/dashboard/v1/status/"

// API's currently poiting to mocked service, make sure its running first.
const RESTCountriesAPI = "localhost:8081/v3.1/alpha/{code}"
const MetroAPI = "localhost:8081/v1/forecast?latitude={lat}&longitude={long}&hourly=temperature_2m,precipitation"
const CurrencyAPI = "localhost:8081/currency/{code}"
