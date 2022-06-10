package controller

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Diegoplas/2022Q2GO-Bootcamp/model"

	"github.com/gorilla/mux"
	"github.com/unrolled/render"
)

type DataGraphersUseCases interface {
	GraphBTCValues(inputDayStr string) (graphURL string, err error)
	GraphCryptoRecords(requestedCryptoCode, requestedDays string) (cryptoGraphURL string, err error)
}

// CHECK ::::
type DataGrapher struct {
}

// CHECK ::::
func NewDataGrapher() DataGrapher {
	return DataGrapher{}
}

type DataGetterUseCases interface {
	CreateCSVFile() error
	CopyResponseToCSVFile(resp *http.Response) error
	GetDataFromHistoricalValueRows(requestedDays int, historicalValueRows [][]string) (records model.CryptoRecordValues, dataError error)
	ExtractRowsFromCSVFile(csvFileName string) (rows [][]string, err error)
	ExtractDataFromBTCCSVRows(requestedDay int, csvRows [][]string) (records model.CryptoRecordValues, dataError error)
}

type CoverterUseCases interface {
	AverageHighLowCryptoPrices(lowPrice, highPrice string) (float64, error)
	ConvertCSVStrToDate(strDate string) (time.Time, error)
	ConvertCSVStrDataToNumericTypes(idStr, priceStr string) (int, float64, error)
}
type GraphUseCases interface {
	MakeGraph(records model.CryptoRecordValues, cryptoCode, days string) string
}
type WorkerPoolUseCases interface {
	CSVWorkerPoolRowExtractor(oddOrEven string, items, itemsPerWorker int) ([]model.CryptoPricesAndDates, error)
	Worker(readChan chan []string, writeChan chan model.CryptoPricesAndDates, itemsPerWorker int)
}
type UseCasesHandler struct {
	dataGetterUseCases   DataGetterUseCases
	graphUseCases        GraphUseCases
	converterUseCases    CoverterUseCases
	workerPoolUseCases   WorkerPoolUseCases
	dataGraphersUseCases DataGraphersUseCases
}

func NewUseCasesHandler(dataGetterUseCases DataGetterUseCases, graphUseCases GraphUseCases,
	converterUseCases CoverterUseCases, workerPoolUseCases WorkerPoolUseCases, dataGraphersUseCases DataGraphersUseCases) UseCasesHandler {
	return UseCasesHandler{
		dataGetterUseCases:   dataGetterUseCases,
		graphUseCases:        graphUseCases,
		converterUseCases:    converterUseCases,
		workerPoolUseCases:   workerPoolUseCases,
		dataGraphersUseCases: dataGraphersUseCases,
	}
}

// GraphCryptoRecordsHandler - Handles the request for graphing and saving historical crypto prices.
func (uch UseCasesHandler) GraphCryptoRecordsHandler(w http.ResponseWriter, r *http.Request) {
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
	graphName, err := uch.dataGraphersUseCases.GraphCryptoRecords(rawInputCryptoCode, rawInputDays)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(string(err.Error())))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("%s Data successfully graphed in file: %s", rawInputCryptoCode, graphName)))
}

// GraphBTCValuesHandler - Handles the request for graphing historical BTC prices from an existing CSV file.
func (uch UseCasesHandler) GraphBTCValuesHandler(w http.ResponseWriter, r *http.Request) {
	rawInputDay := mux.Vars(r)["day"]
	if rawInputDay == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("missing parameter. Please input the desired number of days"))
		return
	}
	graphName, err := uch.dataGraphersUseCases.GraphBTCValues(rawInputDay)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(string(err.Error())))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("BTC Data successfully graphed in file: %s", graphName)))
}

// WorkerPoolHandler - Handles extracting concurrently information from a CSV file.
func (uch UseCasesHandler) WorkerPoolHandler(w http.ResponseWriter, r *http.Request) {
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
	workerResponse, err := uch.workerPoolUseCases.CSVWorkerPoolRowExtractor(oddOrEven, items, itemsPerWorker)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(string(err.Error())))
		return
	}
	render.New().JSON(w, http.StatusOK, &workerResponse)
}
