package app

import (
	"log"
	"os"
	"runtime"

	"github.com/arvaliullin/wapa/internal/domain"
	"github.com/arvaliullin/wapa/internal/hyperfine"
)

type Command struct {
	DesignPayload      domain.DesignPayload
	HyperfinePath      string
	HyperfineResultDir string
	NodePath           string
	ScriptPath         string
	WasmPath           string
	JsPath             string
}

func (c *Command) Run() domain.Experiment {
	designID := c.DesignPayload.ID

	var funcResults []domain.FunctionResult

	for _, fn := range c.DesignPayload.Functions {
		task := domain.Task{
			Function: fn.Function,
			Args:     fn.Args,
			WasmPath: c.WasmPath,
			JsPath:   c.JsPath,
		}

		cmd := hyperfine.NewHyperfineCommand(task, c.DesignPayload)
		if cmd == nil {
			log.Printf("Ошибка создания команды %v\n", task)
			continue
		}

		result, err := cmd.Run()

		if err != nil {
			continue
		}

		if result == nil {
			continue
		}

		if len(result.Results) < 1 {
			continue
		}

		metrics := domain.Metrics{
			Mean:   result.Results[0].Mean,
			Stddev: result.Results[0].Stddev,
			Median: result.Results[0].Median,
			User:   result.Results[0].User,
			System: result.Results[0].System,
			Min:    result.Results[0].Min,
			Max:    result.Results[0].Max,
		}

		funcResults = append(funcResults, domain.FunctionResult{
			FunctionName: fn.Function,
			Args:         fn.Args,
			Repeats:      c.DesignPayload.Repeats,
			Metrics:      metrics,
		})
	}

	hostname, err := os.Hostname()

	if err != nil {
		log.Printf("Can't get hostname %s", err)
	}

	return domain.Experiment{
		DesignID:        designID,
		FunctionResults: funcResults,
		Hostname:        hostname,
		Arch:            runtime.GOARCH,
	}
}
