# BTC-USD Grapher

This API consumes historical prices of Bitcoin on USD dollars.
It makes a graph of them all or from the latest value to a specific date. 

### Requirements

* Go 1.15

### Framework

This project utilizes Gorilla web toolkit.

### 1. Install dependencies

As this project utilizes go modules, the dependencies can be easily downloaded executing the following line:
```
go mod download
```

### 2. Usage

1. The main program should be excecuted from root for the paths to match correcly. To excecute it we can use:
   ```
   go run ./main/.
   ```

2. BTC Historic Values on CSV - /btc-values/{date}

   This endpoint gets the historical prices from a CSV file of the latest n days. Ex, last quarter of the year would be 90 (90 days). 
   If the number is higher than the data stored, it will graph all the data. 
   ```
   Eg. http://localhost:8080/btc-values/90
   ```