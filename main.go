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

package main

import (
	"flag"
	"fmt"

	"github.com/aritzz/simplepipe/load"
	"github.com/aritzz/simplepipe/pipe"
)

const VERSION = "1.0.1"

// main Main function for Simplepipe software
func main() {

	pipelineFile := flag.String("pipeline", "", "pipeline file")
	showArgs := flag.Bool("args", false, "get pipeline argument list")
	timeExec := flag.Bool("time", false, "get global execution time")
	timeExecCmd := flag.Bool("timecmd", false, "get execution time for each command")
	fileLogger := flag.String("logfile", "", "redirect logging to a file")
	onlyOutput := flag.Bool("outputonly", false, "get only output information")
	flag.Parse()

	// Parse pipeline file
	if len(*pipelineFile) == 0 {
		fmt.Println("Pipeline not provided. Use -h to get help.")
		return
	}
	data, err := load.ParseFile(*pipelineFile)
	if err != nil {
		fmt.Println("Error parsing pipeline file: ", err)
		return
	}

	if !*onlyOutput {
		fmt.Println("[] Simplepipe v" + VERSION)
		fmt.Println("Pipeline '" + data.Name + "' loaded")
	}

	argCount := len(data.Input)
	userArgCount := len(flag.Args())

	// See arguments to be provided
	if *showArgs {
		fmt.Print("You must provide ", argCount, " argument(s):")
		for _, value := range data.Input {
			if len(value.Value) == 0 {
				fmt.Print(" <", value.Name, ">")
			} else {
				fmt.Print(" <", value.Value, ">")
			}
		}
		fmt.Print("\n")
		return
	}

	if argCount != userArgCount {
		fmt.Println("Invalid argument count. ", userArgCount, "provided,", argCount, "needed.")
		return
	}

	// Load input data
	pipe.LoadInput(&data, flag.Args())

	if !*onlyOutput {
		fmt.Println("Executing pipeline")
	}

	// Execute pipeline
	pipelineOutput, err := pipe.ExecutePipeline(data, *fileLogger)

	if err != nil {
		fmt.Println(err)
	}

	if !*onlyOutput {
		fmt.Println("Pipeline output: " + pipelineOutput.Output)
	} else {
		fmt.Println(pipelineOutput.Output)
	}

	// Print execution times if needed
	if *timeExec {
		printExectimeGlobal(pipelineOutput)
	}

	if *timeExecCmd {
		printExectimeFunction(pipelineOutput)
	}
}
