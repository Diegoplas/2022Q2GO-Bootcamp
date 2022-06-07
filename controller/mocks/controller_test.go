package mocks

import (
	"fmt"
	reflect "reflect"
	"testing"

	"github.com/Diegoplas/2022Q2GO-Bootcamp/csvdata"
	gomock "github.com/golang/mock/gomock"
)

type myTestReporter struct {
	t gomock.TestReporter
}

func (m *myTestReporter) TestDataGetter_ExtractRowsFromCSVFile(t *testing.T) {
	testFilename := "testfilename.ext"
	expectedReturnedRows := [][]string{
		{"TC1", "test-crypto-1"},
		{"TC2", "test-crypto-2"},
	}
	ctrl := gomock.NewController(m.t)
	defer ctrl.Finish()

	mockDattaGetter := NewMockDataGetter(ctrl)

	dataHandler := csvdata.NewCSVDataHandler()

	gomock.InAnyOrder(
		mockDattaGetter.
			EXPECT().ExtractRowsFromCSVFile(testFilename).Return(expectedReturnedRows, nil),
	)

	gotRows, _ := dataHandler.ExtractRowsFromCSVFile(testFilename)

	if reflect.DeepEqual(gotRows, expectedReturnedRows) {
		fmt.Println("test completed succesfully")
	} else {
		fmt.Println("NOT succesfully")
	}
}
