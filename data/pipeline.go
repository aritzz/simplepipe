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

package data

import (
	"time"
)

const (
	TYPE_ASSIGN ExecutionType = iota
	TYPE_EXECASSIGN
	TYPE_EXEC
)

type ExecutionType int

//
// Pipeline related (before processing)
//

type Pipeline struct {
	Name        string
	Input       []PipelineInput
	Declaration map[string]string
	Output      PipelineOutput
	Execution   []PipelineExecution
}

type PipelineInput struct {
	Name  string
	Value string
}

type PipelineExecution struct {
	Type    ExecutionType
	Command string
	Output  string
}

type PipelineOutput struct {
	Defined bool
	Value   string
}

//
// Pipeline results (after processing)
//

type PipelineResult struct {
	Variables map[string]string
	Time      time.Duration
	ExecStep  []PipelineResultExecStep
	Output    string
}

type PipelineResultExecStep struct {
	Command  string
	Error    string
	ExecTime time.Duration
}
