# USD-Crypto Grapher

This API gets historical prices of Bitcoin or other 575 different crypto currencies on USD dollars, consulting external or local data.
It makes a graph of the last N requested days for having a better visualization of the behavior of that currency on the market. 

### Requirements

* Go 1.15 or newer versions.

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

2. BTC Historic Values on CSV - /btc-values/{day}

   This endpoint gets the historical prices from a CSV file of the latest n days. Ex, if you want to know the prices of the quarter of the year, 
   you would write 90 as the input (that would be 90 days).
   If the number is higher than the number of data stored, it will graph all the existing data.

   ```
   Eg. http://localhost:8080/btc-values/90
   ```
  
    Here is the previous request's graph.

  ![alt text](https://github.com/Diegoplas/2022Q2GO-Bootcamp/blob/second-delivery/historical-usd-BTC-90-days-graph.png)

3. Crypto Historic Values - /usd-crypto-conversion/{cryptoCode}/{days}

   This endpoint gets the historical values on USD dollars from a variaty of crypto currencies, currently 576 different currencies, for a requested number of days. Then write those vales into a CSV File and graph them. 
   
   The data is acquired from Alpha Vantage API, here is a link to the documentation:
   https://www.alphavantage.co/documentation/

   This endpoint have as advantages compared to the previous one that you are not fixed to look only for BTC values and also you don't have to manually update the CSV historical values, the endpoint do it for you.

   ```
   Eg. http://localhost:8080/usd-crypto-conversion/sol/360
   ```

    Here is the previous request's graph.

  ![alt text](https://github.com/Diegoplas/2022Q2GO-Bootcamp/blob/second-delivery/historical-usd-SOL-360-days-graph.png)