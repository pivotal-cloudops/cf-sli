package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/pivotal-cloudops/cf-sli/http_wrapper"
	"os"

	uuid "github.com/nu7hatch/gouuid"
	"github.com/pivotal-cloudops/cf-sli/cf_wrapper"
	"github.com/pivotal-cloudops/cf-sli/config"
	"github.com/pivotal-cloudops/cf-sli/logger"
	"github.com/pivotal-cloudops/cf-sli/sli_executor"
)

type Output struct {
	StartTime   float64 `json:"app_start_time_in_sec"`
	StopTime    float64 `json:"app_stop_time_in_sec"`
	StartStatus int     `json:"app_start_status"`
	StopStatus  int     `json:"app_stop_status"`
}

func main() {
	var config config.Config
	var cf_cli cf_wrapper.CfWrapper
	http_wrapper := &http_wrapper.HttpWrapper{}

	app_bits_path := flag.String("app-bits", "./assets/ruby_simple", "App bits path")
	configLocation := flag.String("config", "./.config", "Config path")

	flag.Parse()

	err := config.LoadConfig(*configLocation)
	if err != nil {
		fmt.Fprint(os.Stderr, "Failed to load .config :\n")
		fmt.Fprint(os.Stderr, err)
		fmt.Fprint(os.Stderr, "\n")
		os.Exit(1)
	}

	guid, err := uuid.NewV4()
	if err != nil {
		panic(err)
	}

	app_name := "cf-sli-app-" + guid.String()[0:18]

	sli_executor := sli_executor.NewSliExecutor(cf_cli, http_wrapper, logger.NewLogger())
	result, err := sli_executor.RunTest(app_name, *app_bits_path, config)

	output := &Output{
		StartTime:   result.StartTime.Seconds(),
		StopTime:    result.StopTime.Seconds(),
		StartStatus: result.StartStatus,
		StopStatus:  result.StopStatus,
	}

	json_output, err := json.Marshal(output)
	if err != nil {
		panic(err)
	}

	fmt.Fprintf(os.Stdout, string(json_output))
}
