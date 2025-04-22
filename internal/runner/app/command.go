package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"

	"github.com/arvaliullin/wapa/internal/domain"
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

		taskJSON, _ := json.Marshal(task)
		cmdStr := fmt.Sprintf("%s '%s %s'", c.NodePath, c.ScriptPath, string(taskJSON))
		hfPath := path.Join(c.HyperfineResultDir, "hf.json")
		hfFlag := fmt.Sprintf("--runs %d --export-json=%s", c.DesignPayload.Repeats, hfPath)

		hyperfineCmd := exec.Command(
			c.HyperfinePath,
			hfFlag,
			cmdStr,
		)

		var hfOut, hfErr bytes.Buffer
		hyperfineCmd.Stdout = &hfOut
		hyperfineCmd.Stderr = &hfErr

		err := hyperfineCmd.Run()

		log.Printf("COMMAND: %s", hyperfineCmd.String())

		var metrics domain.Metrics
		if err == nil {
			hfjson, _ := os.ReadFile(hfPath)
			type HfEntry struct {
				Mean   float64 `json:"mean"`
				Stddev float64 `json:"stddev"`
				Median float64 `json:"median"`
				User   float64 `json:"user"`
				System float64 `json:"system"`
				Min    float64 `json:"min"`
				Max    float64 `json:"max"`
			}
			type HfRes struct {
				Results []HfEntry `json:"results"`
			}
			var hf HfRes
			_ = json.Unmarshal(hfjson, &hf)
			if len(hf.Results) > 0 {
				v := hf.Results[0]
				metrics = domain.Metrics{
					Mean:   v.Mean,
					Stddev: v.Stddev,
					Median: v.Median,
					User:   v.User,
					System: v.System,
					Min:    v.Min,
					Max:    v.Max,
				}
			}
		} else {
			metrics = domain.Metrics{}
		}

		funcResults = append(funcResults, domain.FunctionResult{
			FunctionName: fn.Function,
			Args:         fn.Args,
			Repeats:      c.DesignPayload.Repeats,
			Metrics:      metrics,
		})
	}
	return domain.Experiment{
		DesignID:        designID,
		FunctionResults: funcResults,
	}
}
