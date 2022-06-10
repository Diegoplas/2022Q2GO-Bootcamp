package controller

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/Diegoplas/2022Q2GO-Bootcamp/usecase"

	"github.com/gorilla/mux"
	"github.com/unrolled/render"
)

type DataGrapher struct {
	workerPoolUseCases usecase.WorkPool
}

func NewDataGrapher(workerPoolUseCases usecase.WorkPool) DataGrapher {
	return DataGrapher{
		workerPoolUseCases: workerPoolUseCases,
	}
}

// GraphCryptoRecordsHandler - Handles the request for graphing and saving historical crypto prices.
func (dg DataGrapher) GraphCryptoRecordsHandler(w http.ResponseWriter, r *http.Request) {
	// Load request parameters.
	rawInputCryptoCode := mux.Vars(r)["cryptoCode"]
	if rawInputCryptoCode == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("missing parameter. Please input the desired Crypto Code"))
		return
	}
	rawInputDays := mux.Vars(r)["days"]
	if rawInputDays == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("missing parameter. Please input the desired number of days"))
		return
	}
	graphName, err := dg.workerPoolUseCases.GraphCryptoRecords(rawInputCryptoCode, rawInputDays)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(string(err.Error())))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("%s Data successfully graphed in file: %s", rawInputCryptoCode, graphName)))
}

// GraphBTCValuesHandler - Handles the request for graphing historical BTC prices from an existing CSV file.
func (dg DataGrapher) GraphBTCValuesHandler(w http.ResponseWriter, r *http.Request) {
	rawInputDay := mux.Vars(r)["day"]
	if rawInputDay == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("missing parameter. Please input the desired number of days"))
		return
	}
	graphName, err := dg.workerPoolUseCases.GraphBTCValues(rawInputDay)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(string(err.Error())))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("BTC Data successfully graphed in file: %s", graphName)))
}

// WorkerPoolHandler - Handles extracting concurrently information from a CSV file.
func (dg DataGrapher) WorkerPoolHandler(w http.ResponseWriter, r *http.Request) {
	// Validating input parameters.
	oddOrEven := strings.ToLower(mux.Vars(r)["odd_or_even"])
	if (oddOrEven != "odd" && oddOrEven != "even") || oddOrEven == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("please use odd or even as your first parameter"))
		return
	}
	items, err := strconv.Atoi(mux.Vars(r)["items"])
	if err != nil || items <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("items parameter should be a positive number"))
		return
	}
	itemsPerWorker, err := strconv.Atoi(mux.Vars(r)["items_per_worker"])
	if err != nil || itemsPerWorker <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("items per worker parameter should be a positive number"))
		return
	}
	workerResponse, err := dg.workerPoolUseCases.CSVWorkerPoolRowExtractor(oddOrEven, items, itemsPerWorker)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(string(err.Error())))
		return
	}
	render.New().JSON(w, http.StatusOK, &workerResponse)
}
