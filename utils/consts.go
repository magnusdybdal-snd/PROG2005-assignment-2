package utils

// Paths
const ROOT_PATH = "/"
const REGISTRATION_PATH = "/dashboard/v1/registrations/"
const DASHBOARD_PATH = "/dashboard/v1/dashboards/{id}"
const NOTIFICATION_PATH = "/dashboard/v1/notifications/"
const STATUS_PATH = "/dashboard/v1/status/"

// API's currently poiting to mocked service, make sure its running first.
const RESTCountriesAPI = "http://localhost:8081/v3.1/alpha/"
const MetroAPI = "http://localhost:8081/v1/forecast?latitude=%f&longitude=%f&hourly=temperature_2m,precipitation"
const CurrencyAPI = "http://localhost:8081/currency/"

// collections in firebase
const WebhooksCollection = "webhooks"
const WEBHOOKTESTCOLLECTION = "webhook_test"
const DASHBOARD_COLLECTION = "dashboards"

const INVOKE = "INVOKE"
const REGISTER = "REGISTER"
const CHANGE = "CHANGE"
const DELETE = "DELETE"
const NOTREACHABLE = "NOT_REACHABLE"
const ACCESS_FAILURE = "ACCESS_FAILURE"

// Version of the API
const VERSION = "v1"
