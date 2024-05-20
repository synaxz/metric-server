
# Prometheus Metrics Server

Go application provides an HTTP server for storing and serving Prometheus metrics. The server supports four types of metrics: gauges, counters, histograms, and summaries. Metrics can be stored by sending POST requests to specific endpoints, and they can be accessed via the `/metrics` endpoint.

It suppose to be used in a stateless environment where you can't store metrics in memory between requests. This server will store the metrics in memory and expose them in Prometheus format.

## Table of Contents

- [Use Case](#use-case)
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Endpoints](#endpoints)
    - [Store Gauge](#store-gauge)
    - [Store Counter](#store-counter)
    - [Store Histogram](#store-histogram)
    - [Store Summary](#store-summary)
    - [Metrics Endpoint](#metrics-endpoint)
    - [Accessing Metrics](#accessing-metrics)
- [Contributing](#contributing)
- [License](#license)

## Use Case

This Prometheus Metrics Server is ideal for languages like PHP, which lack persistent memory and shared state between requests

##### Example

Send order metrics from a PHP e-commerce platform to the server. The server aggregates these metrics, making them available to Prometheus for real-time monitoring and alerts, ensuring efficient metrics management in stateless environments.

## Prerequisites

- Go 1.22.3 or higher
- Prometheus client library for Go

## Installation

1. Clone the repository:
    ```sh
    git clone https://github.com/synaxz/metrics-server.git
    cd metrics-server
    ```

2. Install dependencies:
    ```sh
    go mod tidy
    ```

3. Build the application:
    ```sh
    go build -o metrics-server
    ```

4. Run the application:
    ```sh
    ./metrics-server
    ```

## Endpoints

### Store Gauge

- **Endpoint**: `/store/gauge`
- **Method**: POST
- **Payload**:
    ```json
    {
        "key": "gauge_metric",
        "value": 10,
        "labels": {"env": "prod"},
        "help": "Gauge metric",
        "action": "set"
    }
    ```
- **Actions**: `set`, `inc`, `dec`, `add`, `sub`
- **Curl Command**:
    ```sh
    curl -X POST -d '{
      "key": "gauge_metric",
      "value": 10,
      "labels": {"env": "prod"},
      "help": "Gauge metric",
      "action": "set"
    }' http://localhost:8080/store/gauge
    ```

### Store Counter

- **Endpoint**: `/store/counter`
- **Method**: POST
- **Payload**:
    ```json
    {
        "key": "counter_metric",
        "value": 1,
        "labels": {"env": "prod"},
        "help": "Counter metric",
        "action": "inc"
    }
    ```
- **Actions**: `inc`, `add`
- **Curl Command**:
    ```sh
    curl -X POST -d '{
      "key": "counter_metric",
      "value": 1,
      "labels": {"env": "prod"},
      "help": "Counter metric",
      "action": "inc"
    }' http://localhost:8080/store/counter
    ```

### Store Histogram

- **Endpoint**: `/store/histogram`
- **Method**: POST
- **Payload**:
    ```json
    {
        "key": "histogram_metric",
        "value": 10,
        "labels": {"env": "prod"},
        "help": "Histogram metric",
        "action": "observe"
    }
    ```
- **Actions**: `observe`
- **Curl Command**:
    ```sh
    curl -X POST -d '{
      "key": "histogram_metric",
      "value": 10,
      "labels": {"env": "prod"},
      "help": "Histogram metric",
      "action": "observe"
    }' http://localhost:8080/store/histogram
    ```

### Store Summary

- **Endpoint**: `/store/summary`
- **Method**: POST
- **Payload**:
    ```json
    {
        "key": "summary_metric",
        "value": 10,
        "labels": {"env": "prod"},
        "help": "Summary metric",
        "action": "observe"
    }
    ```
- **Actions**: `observe`
- **Curl Command**:
    ```sh
    curl -X POST -d '{
      "key": "summary_metric",
      "value": 10,
      "labels": {"env": "prod"},
      "help": "Summary metric",
      "action": "observe"
    }' http://localhost:8080/store/summary
    ```

### Metrics Endpoint

- **Endpoint**: `/metrics`
- **Method**: GET
- **Description**: Returns all stored metrics in Prometheus format.


### Accessing Metrics

Open your browser and navigate to `http://localhost:8080/metrics` to view all the stored metrics.

## Contributing

Contributions are welcome! Please open an issue or submit a pull request.

## License

This project is licensed under the MIT License.
