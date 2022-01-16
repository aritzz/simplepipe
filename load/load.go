// Copyright (c) 2020 Aritz Olea
// This file is part of Simplepipe <https://github.com/aritzz/simplepipe>
//
// Simplepipe is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// Simplepipe is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with Simplepipe.  If not, see <https://www.gnu.org/licenses/>.

package load

import (
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"regexp"
	"strings"
	"time"

	"github.com/aritzz/simplepipe/data"
)

const RANDOM_LEN = 10

// ParseFile Parses file to a pipeline
// if is not valid, returns an error
func ParseFile(input string) (data.Pipeline, error) {
	var return_pipe data.Pipeline

	// Read file
	fileContent, err := loadFile(input)
	if err != nil {
		return return_pipe, err
	}

	// Get clean pipeline
	pipeline := cleanRawPipeline(fileContent)

	// Load pipeline to a struct
	return loadPipeline(pipeline)
}

// Load a file to a string
func loadFile(input string) (string, error) {
	var retstring string
	filecontent, err := ioutil.ReadFile(input)
	if err != nil {
		return retstring, err
	}
	return string(filecontent), nil
}

// Clean raw pipeline
func cleanRawPipeline(pipe string) []string {
	fileCleaned := []string{}
	fileLine := strings.Split(pipe, "\n")

	for _, line := range fileLine {
		line_clean := strings.TrimSpace(line)
		if len(line) > 0 {
			if line_clean[0] != ';' {
				fileCleaned = append(fileCleaned, line_clean)
			}
		}
	}

	return fileCleaned
}

// Pipeline loader
func loadPipeline(pipeline []string) (data.Pipeline, error) {
	var currentStatus int
	var err error
	pipelineData := data.Pipeline{}
	pipelineData.Declaration = make(map[string]string)
	currentStatus = STATUS_DEFINE

	// for i, line := range pipeline {
	for i := 0; i < len(pipeline); i++ {
		line := pipeline[i]
		switch currentStatus {
		case STATUS_DEFINE:
			if pipelineData.Name, err = getPipelineDefinition(line); err != nil {
				goto retpipe
			}
			currentStatus = STATUS_DECLARATION
		case STATUS_DECLARATION:
			pipelineData, currentStatus, err = getPipelineDeclaration(line, pipelineData)
			if err != nil {
				goto retpipe
			}
		case STATUS_PIPELINE:
			pipelineData, currentStatus, err = getPipelineContent(line, pipelineData)
			if err != nil {
				goto retpipe
			}
			if currentStatus == STATUS_END {
				i -= 1
			}
		case STATUS_END:
			pipelineData, currentStatus, err = getPipelineEnd(line, pipelineData)
			if err != nil {
				goto retpipe
			}
		}
	}

retpipe:
	return pipelineData, err
}

// Get pipeline content
func getPipelineContent(line string, pipeline data.Pipeline) (data.Pipeline, int, error) {

	// Pipeline content ends here
	if isPipelineEnd(line) {
		return pipeline, STATUS_END, nil
	}

	// Assign with execution
	execassign := regexp.MustCompile(`^(\w+)\s*=\s*\((.+)\)$`)
	if len(execassign.FindStringSubmatch(line)) == 3 {
		executionData := data.PipelineExecution{Type: data.TYPE_EXECASSIGN, Command: execassign.FindStringSubmatch(line)[2], Output: execassign.FindStringSubmatch(line)[1]}
		pipeline.Execution = append(pipeline.Execution, executionData)
		return pipeline, STATUS_PIPELINE, nil
	}

	// Simple assign
	onlyassign := regexp.MustCompile(`^(\w+)\s*=\s*(\w+)$`)
	if len(onlyassign.FindStringSubmatch(line)) == 3 {
		executionData := data.PipelineExecution{Type: data.TYPE_ASSIGN, Command: onlyassign.FindStringSubmatch(line)[2], Output: onlyassign.FindStringSubmatch(line)[1]}
		pipeline.Execution = append(pipeline.Execution, executionData)
		return pipeline, STATUS_PIPELINE, nil
	}

	// Only execute
	onlyexec := regexp.MustCompile(`^\((.+)\)$`)
	if len(onlyexec.FindStringSubmatch(line)) == 2 {
		executionData := data.PipelineExecution{Type: data.TYPE_EXEC, Command: onlyexec.FindStringSubmatch(line)[1], Output: ""}
		pipeline.Execution = append(pipeline.Execution, executionData)
		return pipeline, STATUS_PIPELINE, nil
	}

	return pipeline, STATUS_PIPELINE, errors.New("Invalid line: " + line)
}

func getPipelineEnd(line string, pipeline data.Pipeline) (data.Pipeline, int, error) {
	// Output section
	outputsec := regexp.MustCompile(`^end\s*([\w]*)$`)
	if len(outputsec.FindStringSubmatch(line)) == 2 {
		pipeline.Output = data.PipelineOutput{Defined: true, Value: outputsec.FindStringSubmatch(line)[1]}
		return pipeline, STATUS_END, nil
	}

	if len(outputsec.FindStringSubmatch(line)) == 1 {
		fmt.Println(outputsec.FindStringSubmatch(line))
	}

	return pipeline, STATUS_END, errors.New("Pipeline syntax error: " + line)
}

// Get pipeline declaration
func getPipelineDeclaration(line string, pipeline data.Pipeline) (data.Pipeline, int, error) {

	// Declaration end
	if line == "begin" {
		return pipeline, STATUS_PIPELINE, nil
	}

	// Declaration section
	declaration := regexp.MustCompile(`^use ([\w]+)$`)
	if len(declaration.FindStringSubmatch(line)) == 2 {
		pipeline.Declaration[declaration.FindStringSubmatch(line)[1]] = ""
		return pipeline, STATUS_DECLARATION, nil
	}

	// Random declaration
	randomvars := regexp.MustCompile(`^rand ([\w]+)$`)
	if len(randomvars.FindStringSubmatch(line)) == 2 {
		pipeline.Declaration[randomvars.FindStringSubmatch(line)[1]] = getRandomString(RANDOM_LEN)
		return pipeline, STATUS_DECLARATION, nil
	}

	// Input section
	inputsec := regexp.MustCompile(`^read ([\w]+)(\s+"(.*)")?$`)
	if len(inputsec.FindStringSubmatch(line)) == 2 {
		pipeline.Input = append(pipeline.Input, data.PipelineInput{Name: inputsec.FindStringSubmatch(line)[1], Value: ""})
		return pipeline, STATUS_DECLARATION, nil
	}
	if len(inputsec.FindStringSubmatch(line)) == 4 {
		// pipeline.Input[inputsec.FindStringSubmatch(line)[1]] = inputsec.FindStringSubmatch(line)[3]
		pipeline.Input = append(pipeline.Input, data.PipelineInput{Name: inputsec.FindStringSubmatch(line)[1], Value: inputsec.FindStringSubmatch(line)[3]})
		return pipeline, STATUS_DECLARATION, nil
	}

	// Invalid
	return pipeline, STATUS_DECLARATION, errors.New("Invalid line in declaration: " + line)
}

// Get pipeline definition
func getPipelineDefinition(line string) (string, error) {
	var retstring string
	var err error

	re := regexp.MustCompile(`^pipeline ([\w]+)$`)
	if len(re.FindStringSubmatch(line)) != 2 {
		err = errors.New("Invalid pipeline definition: " + line)
	} else {
		retstring = re.FindStringSubmatch(line)[1]
	}

	return retstring, err
}

// Detect ending
func isPipelineEnd(line string) bool {
	outputsec := regexp.MustCompile(`^end\s*([\w]*)$`)
	return outputsec.MatchString(line)
}

// Get random string
func getRandomString(length int) string {
	var b strings.Builder
	rand.Seed(time.Now().UnixNano())
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")

	for i := 0; i < length; i++ {
		b.WriteRune(chars[rand.Intn(len(chars))])
	}
	return b.String()
}
