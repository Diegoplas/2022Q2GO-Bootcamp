package graph

import (
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/Diegoplas/2022Q2GO-Bootcamp/model"
)

func TestGrapher_MakeGraph(t *testing.T) {

	dateGraph1 := time.Date(2022, 06, 05, 0, 0, 0, 0, time.UTC)
	dateGraph2 := time.Date(2022, 06, 04, 0, 0, 0, 0, time.UTC)
	dateGraph3 := time.Date(2022, 06, 03, 0, 0, 0, 0, time.UTC)

	cryptoRecords := model.CryptoRecordValues{
		Dates:        []time.Time{dateGraph1, dateGraph2, dateGraph3},
		AveragePrice: []float64{39.085, 37.605, 39.585},
		MaxPrice:     39.585,
		MinPrice:     37.605,
	}

	validTestCryptoCode := "test"
	validTestCryptoDays := "182"
	validTestGraphName := fmt.Sprintf("historical-usd-%s-%s-days-graph.png",
		validTestCryptoCode, validTestCryptoDays)

	expectedCreatedGraphPath := "historical-usd-test-182-days-graph.png"

	type args struct {
		records    model.CryptoRecordValues
		cryptoCode string
		days       string
	}
	tests := []struct {
		name              string
		g                 Grapher
		args              args
		wantGraphFilename string
	}{
		{
			name: "Valid test",
			args: args{
				records:    cryptoRecords,
				cryptoCode: "test",
				days:       "182",
			},
			wantGraphFilename: validTestGraphName,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotGraphFilename := tt.g.MakeGraph(tt.args.records, tt.args.cryptoCode, tt.args.days); gotGraphFilename != tt.wantGraphFilename {
				t.Errorf("Grapher.MakeGraph() = %v, want %v", gotGraphFilename, tt.wantGraphFilename)
			}
			if _, err := os.Stat(expectedCreatedGraphPath); errors.Is(err, os.ErrNotExist) {
				t.Errorf("Grapher.MakeGraph() File %v, was not created.", expectedCreatedGraphPath)
			}
			err := os.Remove(expectedCreatedGraphPath)
			if err != nil {
				t.Errorf("Was not possible to delete file %v.", expectedCreatedGraphPath)
			}
		})
	}
}
