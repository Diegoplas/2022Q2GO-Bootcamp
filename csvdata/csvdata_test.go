package csvdata

import (
	"reflect"
	"testing"
	"time"

	"github.com/Diegoplas/2022Q2GO-Bootcamp/model"
)

func TestCSVDataHandler_GetDataFromHistoricalValueRows(t *testing.T) {
	validRequestRows := [][]string{
		{"timestamp", "open (USD)", "high (USD)", "low (USD)", "close (USD)", "open (USD)", "high (USD)", "low (USD)", "close (USD)", "volume,market cap (USD)"},
		{"2022-06-05", "39.02000000", "39.47000000", "38.70000000", "39.39000000", "39.02000000", "39.47000000", "38.70000000", "39.39000000", "102545.70000000", "102545.70000000"},
		{"2022-06-04", "38.28000000", "39.50000000", "35.71000000", "39.00000000", "38.28000000", "39.50000000", "35.71000000", "39.00000000", "3736176.79000000", "3736176.79000000"},
		{"2022-06-03", "40.86000000", "41.54000000", "37.63000000", "38.28000000", "40.86000000", "41.54000000", "37.63000000", "38.28000000", "4437324.64000000", "4437324.64000000"},
	}

	invalidTimestampExternalRequestRows := [][]string{
		{"timestamp", "open (USD)", "high (USD)", "low (USD)", "close (USD)", "open (USD)", "high (USD)", "low (USD)", "close (USD)", "volume,market cap (USD)"},
		{"invalidType", "39.02000000", "39.47000000", "38.70000000", "39.39000000", "39.02000000", "39.47000000", "38.70000000", "39.39000000", "102545.70000000", "102545.70000000"},
		{"2022-06-04", "38.28000000", "39.50000000", "35.71000000", "39.00000000", "38.28000000", "39.50000000", "35.71000000", "39.00000000", "3736176.79000000", "3736176.79000000"},
		{"2022-06-03", "40.86000000", "41.54000000", "37.63000000", "38.28000000", "40.86000000", "41.54000000", "37.63000000", "38.28000000", "4437324.64000000", "4437324.64000000"},
	}
	invalidHighPriceExternalRequestRows := [][]string{
		{"timestamp", "open (USD)", "high (USD)", "low (USD)", "close (USD)", "open (USD)", "high (USD)", "low (USD)", "close (USD)", "volume,market cap (USD)"},
		{"2022-06-05", "39.02000000", "39.47000000", "38.70000000", "39.39000000", "39.02000000", "39.47000000", "38.70000000", "39.39000000", "102545.70000000", "102545.70000000"},
		{"2022-06-04", "38.28000000", "invalidType", "35.71000000", "39.00000000", "38.28000000", "39.50000000", "35.71000000", "39.00000000", "3736176.79000000", "3736176.79000000"},
		{"2022-06-03", "40.86000000", "41.54000000", "37.63000000", "38.28000000", "40.86000000", "41.54000000", "37.63000000", "38.28000000", "4437324.64000000", "4437324.64000000"},
	}
	invalidLowPriceExternalRequestRows := [][]string{
		{"timestamp", "open (USD)", "high (USD)", "low (USD)", "close (USD)", "open (USD)", "high (USD)", "low (USD)", "close (USD)", "volume,market cap (USD)"},
		{"2022-06-05", "39.02000000", "39.47000000", "38.70000000", "39.39000000", "39.02000000", "39.47000000", "38.70000000", "39.39000000", "102545.70000000", "102545.70000000"},
		{"2022-06-04", "38.28000000", "39.50000000", "35.71000000", "39.00000000", "38.28000000", "39.50000000", "35.71000000", "39.00000000", "3736176.79000000", "3736176.79000000"},
		{"2022-06-03", "40.86000000", "41.54000000", "invalidType", "38.28000000", "40.86000000", "41.54000000", "37.63000000", "38.28000000", "4437324.64000000", "4437324.64000000"},
	}
	wantedDate1 := time.Date(2022, 06, 05, 0, 0, 0, 0, time.UTC)
	wantedDate2 := time.Date(2022, 06, 04, 0, 0, 0, 0, time.UTC)
	wantedDate3 := time.Date(2022, 06, 03, 0, 0, 0, 0, time.UTC)

	wantedCryptoRecords := model.CryptoRecordValues{
		Dates:        []time.Time{wantedDate1, wantedDate2, wantedDate3},
		AveragePrice: []float64{39.085, 37.605000000000004, 39.585},
		MaxPrice:     39.585,
		MinPrice:     37.605000000000004,
	}
	type args struct {
		requestedDays       int
		historicalValueRows [][]string
	}
	tests := []struct {
		name        string
		csvdh       CSVDataHandler
		args        args
		wantRecords model.CryptoRecordValues
		wantErr     bool
	}{
		{
			name: "Valid test",
			args: args{
				requestedDays:       3,
				historicalValueRows: validRequestRows,
			},
			wantRecords: wantedCryptoRecords,
			wantErr:     false,
		},
		{
			name: "Invalid test - invalid Timestamp type",
			args: args{
				requestedDays:       3,
				historicalValueRows: invalidTimestampExternalRequestRows,
			},
			wantRecords: model.CryptoRecordValues{},
			wantErr:     true,
		},
		{
			name: "Invalid test - invalid High Price type",
			args: args{
				requestedDays:       3,
				historicalValueRows: invalidHighPriceExternalRequestRows,
			},
			wantRecords: model.CryptoRecordValues{},
			wantErr:     true,
		},
		{
			name: "Invalid test - invalid Low Price type",
			args: args{
				requestedDays:       3,
				historicalValueRows: invalidLowPriceExternalRequestRows,
			},
			wantRecords: model.CryptoRecordValues{},
			wantErr:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRecords, err := tt.csvdh.GetDataFromHistoricalValueRows(tt.args.requestedDays, tt.args.historicalValueRows)
			if (err != nil) != tt.wantErr {
				t.Errorf("CSVDataHandler.GetDataFromHistoricalValueRows() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRecords, tt.wantRecords) {
				t.Errorf("CSVDataHandler.GetDataFromHistoricalValueRows() = %v, want %v", gotRecords, tt.wantRecords)
			}
		})
	}
}

func TestCSVDataHandler_ExtractDataFromBTCCSVRows(t *testing.T) {
	//   id    timestamp		     price
	validBTCRequestRows := [][]string{
		{"Id", "Timestamp", "market-price"},
		{"0", "2022-05-06 0:00:00", "36549.9"},
		{"1", "2022-05-05 0:00:00", "39674.89"},
		{"2", "2022-05-04 0:00:00", "37727.19"},
	}

	invalidIDBTCRequestRows := [][]string{
		{"Id", "Timestamp", "market-price"},
		{"invalidtype", "2022-05-06 0:00:00", "36549.9"},
		{"1", "2022-05-05 0:00:00", "39674.89"},
		{"2", "2022-05-04 0:00:00", "37727.19"},
	}

	invalidTimestampBTCRequestRows := [][]string{
		{"Id", "Timestamp", "market-price"},
		{"0", "2022-05-06 0:00:00", "36549.9"},
		{"1", "2022-05-05 0:00:00", "39674.89"},
		{"2", "invalidtype", "37727.19"},
	}

	invalidPriceBTCRequestRows := [][]string{
		{"Id", "Timestamp", "market-price"},
		{"0", "2022-05-06 0:00:00", "36549.9"},
		{"1", "2022-05-05 0:00:00", "invalidtype"},
		{"2", "2022-05-04 0:00:00", "37727.19"},
	}
	wantedBTCDate1 := time.Date(2022, 05, 06, 0, 0, 0, 0, time.UTC)
	wantedBTCDate2 := time.Date(2022, 05, 05, 0, 0, 0, 0, time.UTC)
	wantedBTCDate3 := time.Date(2022, 05, 04, 0, 0, 0, 0, time.UTC)

	wantedValidBTCRecords := model.CryptoRecordValues{
		Ids:          []int{0, 1, 2},
		Dates:        []time.Time{wantedBTCDate1, wantedBTCDate2, wantedBTCDate3},
		AveragePrice: []float64{36549.9, 39674.89, 37727.19},
		MaxPrice:     39674.89,
		MinPrice:     36549.9,
	}

	type args struct {
		requestedDay int
		csvRows      [][]string
	}
	tests := []struct {
		name           string
		csvdh          CSVDataHandler
		args           args
		wantBTCRecords model.CryptoRecordValues
		wantErr        bool
	}{
		{
			name: "Valid test",
			args: args{
				requestedDay: 2,
				csvRows:      validBTCRequestRows},
			wantBTCRecords: wantedValidBTCRecords,
			wantErr:        false,
		},
		{
			name: "Valid test - requested more days than existing CSV records",
			args: args{
				requestedDay: 8,
				csvRows:      validBTCRequestRows},
			wantBTCRecords: wantedValidBTCRecords,
			wantErr:        false,
		},
		{
			name: "Invalid test - invalid ID type",
			args: args{
				requestedDay: 2,
				csvRows:      invalidIDBTCRequestRows},
			wantBTCRecords: model.CryptoRecordValues{},
			wantErr:        true,
		},
		{
			name: "Invalid test - invalid Timestamp type",
			args: args{
				requestedDay: 2,
				csvRows:      invalidTimestampBTCRequestRows},
			wantBTCRecords: model.CryptoRecordValues{},
			wantErr:        true,
		},
		{
			name: "Invalid test - invalid Price type",
			args: args{
				requestedDay: 2,
				csvRows:      invalidPriceBTCRequestRows},
			wantBTCRecords: model.CryptoRecordValues{},
			wantErr:        true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRecords, err := tt.csvdh.ExtractDataFromBTCCSVRows(tt.args.requestedDay, tt.args.csvRows)
			if (err != nil) != tt.wantErr {
				t.Errorf("CSVDataHandler.ExtractDataFromBTCCSVRows() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRecords, tt.wantBTCRecords) {
				t.Errorf("CSVDataHandler.ExtractDataFromBTCCSVRows() = %v, want %v", gotRecords, tt.wantBTCRecords)
			}
		})
	}
}
func Test_convertCSVStrDataToNumericTypes(t *testing.T) {
	type args struct {
		idStr    string
		priceStr string
	}
	tests := []struct {
		name      string
		args      args
		wantId    int
		wantPrice float64
		wantErr   bool
	}{
		{
			name: "Valid test",
			args: args{
				idStr:    "6",
				priceStr: "123.45"},
			wantId:    6,
			wantPrice: 123.45,
			wantErr:   false,
		},
		{
			name: "Data error - ID",
			args: args{
				idStr:    "wrongTipe",
				priceStr: "123.45"},
			wantId:    0,
			wantPrice: 0.0,
			wantErr:   true,
		},
		{
			name: "Data error - Price",
			args: args{
				idStr:    "6",
				priceStr: "wrongTipe"},
			wantId:    0,
			wantPrice: 0.0,
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := NewCSVDataConverter().ConvertCSVStrDataToNumericTypes(tt.args.idStr, tt.args.priceStr)
			if (err != nil) != tt.wantErr {
				t.Errorf("convertCSVStrDataToNumericTypes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.wantId {
				t.Errorf("convertCSVStrDataToNumericTypes() got = %v, want %v", got, tt.wantId)
			}
			if got1 != tt.wantPrice {
				t.Errorf("convertCSVStrDataToNumericTypes() got1 = %v, want %v", got1, tt.wantPrice)
			}
		})
	}
}

func Test_averageHighLowCryptoPrices(t *testing.T) {
	type args struct {
		lowPrice  string
		highPrice string
	}
	tests := []struct {
		name        string
		args        args
		wantAverage float64
		wantErr     bool
	}{
		{
			name: "Valid test",
			args: args{
				lowPrice:  "28.6",
				highPrice: "54.2"},
			wantAverage: 41.400000000000006,
			wantErr:     false,
		},
		{
			name: "Invalid test- Invalid Low Price",
			args: args{
				lowPrice:  "invalid type",
				highPrice: "54.2"},
			wantAverage: 0.0,
			wantErr:     true,
		},
		{
			name: "Invalid test- Invalid High Price",
			args: args{
				lowPrice:  "25.6",
				highPrice: "invalid type"},
			wantAverage: 0.0,
			wantErr:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotAverage, err := NewCSVDataConverter().AverageHighLowCryptoPrices(tt.args.lowPrice, tt.args.highPrice)
			if (err != nil) != tt.wantErr {
				t.Errorf("averageHighLowCryptoPrices() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotAverage != tt.wantAverage {
				t.Errorf("averageHighLowCryptoPrices() = %v, want %v", gotAverage, tt.wantAverage)
			}
		})
	}
}

func Test_convertCSVStrToDate(t *testing.T) {
	type args struct {
		strDate string
	}

	wantedDate := time.Date(2021, 11, 25, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name     string
		args     args
		wantTime time.Time
		wantErr  bool
	}{
		{
			name: "Valid test",
			args: args{
				strDate: "2021-11-25"},
			wantTime: wantedDate,
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTime, err := NewCSVDataConverter().ConvertCSVStrToDate(tt.args.strDate)
			if (err != nil) != tt.wantErr {
				t.Errorf("convertCSVStrToDate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotTime, tt.wantTime) {
				t.Errorf("convertCSVStrToDate() = %v, want %v", gotTime, tt.wantTime)
			}
		})
	}
}
