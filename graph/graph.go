package graph

//go:generate go run main.go

import (
	"fmt"
	"os"

	"github.com/Diegoplas/2022Q2GO-Bootcamp/model"

	chart "github.com/wcharczuk/go-chart/v2"
	"github.com/wcharczuk/go-chart/v2/drawing"
)

type Grapher struct {
}

func NewGrapher() Grapher {
	return Grapher{}
}

// MakeGraph - Generates the graph with the values/dates and saves it as an image.
func (g Grapher) MakeGraph(records model.CryptoRecordValues, cryptoCode, days string) string {
	// Use the 10% of the Max/Min value of the records to have a better view of the graph values.
	GraphTopSpace := records.MaxPrice * .1
	GraphBottomSpace := records.MinPrice * .1
	graphFileName := fmt.Sprintf("historical-usd-%s-%s-days-graph.png", cryptoCode, days)
	priceSeries := chart.TimeSeries{
		Name: "SPY",
		Style: chart.Style{
			StrokeColor: chart.GetDefaultColor(0),
		},
		XValues: records.Dates,
		YValues: records.AveragePrice,
	}

	smaSeries := chart.SMASeries{
		Name: "SPY - SMA",
		Style: chart.Style{
			StrokeColor:     drawing.ColorRed,
			StrokeDashArray: []float64{5.0, 5.0},
		},
		InnerSeries: priceSeries,
	}

	bbSeries := &chart.BollingerBandsSeries{
		Name: "SPY - Bol. Bands",
		Style: chart.Style{
			StrokeColor: drawing.ColorFromHex("efefef"),
			FillColor:   drawing.ColorFromHex("efefef").WithAlpha(64),
		},
		InnerSeries: priceSeries,
	}

	graph := chart.Chart{
		Title: fmt.Sprintf("USD-%s Average Price History of the last %s Days", cryptoCode, days),
		XAxis: chart.XAxis{
			TickPosition: chart.TickPositionBetweenTicks,
		},
		YAxis: chart.YAxis{
			Range: &chart.ContinuousRange{
				//Values added just to have a better display of the graph/title
				Max: records.MaxPrice + GraphTopSpace,
				Min: records.MinPrice - GraphBottomSpace,
			},
		},
		Series: []chart.Series{
			bbSeries,
			priceSeries,
			smaSeries,
		},
	}

	f, _ := os.Create(graphFileName)
	defer f.Close()
	graph.Render(chart.PNG, f)
	return graphFileName
}
