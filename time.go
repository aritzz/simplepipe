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
	"fmt"

	"github.com/aritzz/simplepipe/data"
)

// printExectimeGlobal Print global execution time from pipeline
func printExectimeGlobal(pipeline data.PipelineResult) {
	fmt.Println("Execution time:", pipeline.Time)
}

// printExectimeFunction Print function execution time from pipeline
func printExectimeFunction(pipeline data.PipelineResult) {
	for _, el := range pipeline.ExecStep {
		fmt.Println("Command [", el.Command, "] - Time [", el.ExecTime, "]")
	}
}
