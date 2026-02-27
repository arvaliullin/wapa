package hyperfine

import (
	"bytes"
	"encoding/json"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"

	"github.com/arvaliullin/wapa/internal/domain"
)

type HyperfineCommand struct {
	Task                    domain.Task
	HyperfineResultFilePath string
	Command                 *exec.Cmd
}

func NewHyperfineCommand(task domain.Task, designPayload domain.DesignPayload) *HyperfineCommand {
	scripts := map[string]string{
		"cpp":        "/opt/wapa/scripts/cpp.js",
		"go":         "/opt/wapa/scripts/go.js",
		"rust":       "/opt/wapa/scripts/rs.js",
		"javascript": "/opt/wapa/scripts/js.js",
	}

	resultDir := filepath.Join(os.TempDir(), designPayload.ID)

	err := os.MkdirAll(resultDir, os.ModePerm)
	if err != nil {
		log.Printf("Не удалось создать каталог %s: %v", resultDir, err)
		return nil
	}

	cmd := HyperfineCommand{
		Task:                    task,
		HyperfineResultFilePath: filepath.Join(resultDir, "hyperfine.json"),
	}

	scriptPath, ok := scripts[designPayload.Lang]
	if !ok {
		return nil
	}

	args := []string{
		"--runs", strconv.Itoa(designPayload.Repeats),
	}

	if designPayload.Warmup {
		args = append(args, "--warmup", "15")
	}

	args = append(args, "bun "+scriptPath)

	args = append(args, "--export-json", cmd.HyperfineResultFilePath)
	args = append(args, "--show-output")

	cmd.Command = exec.Command("hyperfine", args...)
	cmd.Command.Env = append(os.Environ(), makeTaskEnvStr(task))
	return &cmd
}

func makeTaskEnvStr(task domain.Task) (taskEnv string) {
	taskStr, err := json.Marshal(task)
	if err != nil {
		return taskEnv
	}
	taskEnv = "TASK_JSON=" + string(taskStr)
	return taskEnv
}

func (command *HyperfineCommand) Run() (*HyperfineResult, error) {
	var stdoutBuf, stderrBuf bytes.Buffer
	command.Command.Stdout = &stdoutBuf
	command.Command.Stderr = &stderrBuf

	log.Printf("COMMAND: %s", command.Command.String())
	log.Printf("COMMAND ENV : %s", command.Command.Env)
	if err := command.Command.Run(); err != nil {
		log.Printf("COMMAND ERROR: %s", stderrBuf.String())
		return nil, err
	}

	log.Printf("COMMAND OUTPUT: %s", stdoutBuf.String())

	data, err := os.ReadFile(command.HyperfineResultFilePath)
	if err != nil {
		return nil, err
	}

	var result HyperfineResult
	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
