package main

import (
	"bufio"
	"fmt"
	"github.com/go-graphite/carbonapi/expr"
	"github.com/go-graphite/carbonapi/expr/functions"
	"github.com/go-graphite/carbonapi/expr/rewrite"
	"github.com/go-graphite/carbonapi/expr/types"
	"github.com/go-graphite/carbonapi/pkg/parser"
	pb "github.com/go-graphite/protocol/carbonapi_v3_pb"
	"math"
	"os"
	"strconv"

	"github.com/guptarohit/asciigraph"
)

type graphRequest struct {
	metric   string
	request  string
	values   []float64
	stepTime int64
	from     int64
	until    int64
}

func readLoop(do func(string)) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		do(scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}
}

func main() {
	userExpr := "stdin"
	if len(os.Args) > 1 {
		userExpr = os.Args[1]
	}

	parsedExpr, e, err := parser.ParseExpr(userExpr)
	if err != nil || e != "" {
		fmt.Errorf("error='%v', leftovers='%v'", err, e)
	}

	var ts []float64
	readLoop(
		func(s string) {
			x, err := strconv.ParseFloat(s, 64)
			if err == nil {
				ts = append(ts, x)
				if len(ts) > 60*60 {
					ts = ts[len(ts)>>1:]
				}
				plot(parsedExpr, ts)
			}
		},
	)
}

func plot(exp parser.Expr, ts []float64) {
	metric := "stdin"
	from := int64(0)
	until := int64(len(ts))

	metricData := types.MetricData{
		FetchResponse: pb.FetchResponse{
			Name:              metric,
			StartTime:         from,
			StopTime:          until,
			StepTime:          1,
			Values:            ts,
			ConsolidationFunc: "average",
			PathExpression:    metric,
		},
	}

	metricMap := make(map[parser.MetricRequest][]*types.MetricData)
	{
		request := parser.MetricRequest{
			Metric: metric,
			From:   from,
			Until:  until,
		}
		metricMap[request] = []*types.MetricData{
			&metricData,
		}
	}

	rewrite.New(make(map[string]string))
	functions.New(make(map[string]string))
	out, err := expr.EvalExpr(exp, from, until, metricMap)
	if err != nil {
		fmt.Errorf("error='%v' expr='%v'", err, exp)
	}
	fmt.Println(metricMap)
	fmt.Println(out)
	if len(out) > 0 {
		values := NaNToZero(out[0].Values)
		if len(values) > 1 {
			fmt.Println(values)
			graph := asciigraph.Plot(values)
			fmt.Println(graph)
		}
	}
}

func NaNToZero(xs []float64) []float64 {
	var ys []float64
	for _, x := range xs {
		if math.IsNaN(x) {
			x = 0
		}
		ys = append(ys, x)
	}
	return ys
}
