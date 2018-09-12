package main

import (
	"fmt"
	"github.com/go-graphite/carbonapi/expr"
	"github.com/go-graphite/carbonapi/expr/functions"
	"github.com/go-graphite/carbonapi/expr/rewrite"
	"github.com/go-graphite/carbonapi/expr/types"
	"github.com/go-graphite/carbonapi/pkg/parser"
	pb "github.com/go-graphite/protocol/carbonapi_v3_pb"
	//	ui "github.com/gizak/termui"
)

type graphRequest struct {
	metric   string
	request  string
	values   []float64
	stepTime int64
	from     int64
	until    int64
}

func main() {
	ee := graphRequest{
		metric:   "stdin",
		request:  "derivative(stdin)",
		from:     1437127020,
		until:    1437127140,
		values:   []float64{0, 1, 2, 4, 8, 16, 16, 16},
		stepTime: 60,
	}

	rewrite.New(make(map[string]string))
	functions.New(make(map[string]string))

	exp, e, err := parser.ParseExpr(ee.request)
	if err != nil || e != "" {
		fmt.Errorf("error='%v', leftovers='%v'", err, e)
	}

	metricData := types.MetricData{
		FetchResponse: pb.FetchResponse{
			Name:              ee.metric,
			StartTime:         ee.from,
			StopTime:          ee.until,
			StepTime:          ee.stepTime,
			Values:            ee.values,
			ConsolidationFunc: "average",
			PathExpression:    ee.metric,
		},
	}

	metricMap := make(map[parser.MetricRequest][]*types.MetricData)
	{
		request := parser.MetricRequest{
			Metric: ee.metric,
			From:   ee.from,
			Until:  ee.until,
		}
		metricMap[request] = []*types.MetricData{
			&metricData,
		}
	}

	out, err := expr.EvalExpr(exp, ee.from, ee.until, metricMap)
	if err != nil {
		fmt.Errorf("error='%v' expr='%v'", err, exp)
	}
	fmt.Printf("exp: %v\n", exp)
	fmt.Printf("metricMap: %v\n", out)

	/*

		err := ui.Init()
		if err != nil {
			panic(err)
		}
		defer ui.Close()



			lc0 := ui.NewLineChart()
			lc0.BorderLabel = "braille-mode Line Chart"
			lc0.Data = metricData.Values
			lc0.Width = 600
			lc0.Height = 600
			lc0.X = 0
			lc0.Y = 0
			lc0.AxesColor = ui.ColorWhite
			lc0.LineColor = ui.ColorGreen | ui.AttrBold

			ui.Render(lc0)

			ui.Handle("q", func(ui.Event) {
				ui.StopLoop()
			})

			//	ui.Loop()
	*/

}
