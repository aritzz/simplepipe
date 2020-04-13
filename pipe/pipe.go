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


package pipe

import(
  "github.com/aritzz/simplepipe/data"
  "os"
  "log"
  "strings"
  "errors"
  "time"
)

var loggerfile *os.File

// Loads an slice of inputs to a pipeline data object
func LoadInput(pipeline *data.Pipeline, input []string) {

  for i, _ := range pipeline.Input {
    if i > len(input) {
      break
    }

    pipeline.Input[i].Value = input[i]
  }


}

// Executes a pipeline
func ExecutePipeline(pipeline data.Pipeline, logredirect string) (data.PipelineResult, error) {
  var pipeline_ret data.PipelineResult
  var err_ret error

  // Exec time
  start_time := time.Now()

  // Enable logging if needed
  if len(strings.TrimSpace(logredirect)) > 0 {
    if err := setLogRedirection(logredirect); err != nil {
      return pipeline_ret, err
    }
    defer loggerfile.Close()
  }

  log.Println("Starting pipeline "+pipeline.Name)

  // Do execution
  pipeline_ret = initVariables(pipeline)
  for i, execItem := range pipeline.Execution {
    log.Println("Running [", execItem.Command,"]")
    pipeline_ret, err_ret = execStep(execItem, pipeline_ret, i)
    if err_ret == nil {
      log.Println("Finished in ", pipeline_ret.ExecStep[i].ExecTime)
    } else {
      log.Println("Execution error ", err_ret.Error())
      break
    }

  }

  if err_ret == nil {
    pipeline_ret, err_ret = getPipelineOutput(pipeline_ret, pipeline)
  }

  pipeline_ret.Time = time.Since(start_time)
  log.Println("Pipeline execution time", pipeline_ret.Time)
  return pipeline_ret, err_ret
}

func execStep(execstep data.PipelineExecution, prevresult data.PipelineResult, i int) (data.PipelineResult, error) {
  var err_ret error
  pipeline_ret := prevresult

  switch execstep.Type {
  case data.TYPE_ASSIGN:
    return execStepAssign(execstep, prevresult, i)
  case data.TYPE_EXECASSIGN:
    return execStepExecAssign(execstep, prevresult, i)
  case data.TYPE_EXEC:
    return execStepExec(execstep, prevresult, i)
  }


  return pipeline_ret, err_ret
}

func execStepAssign(execstep data.PipelineExecution, prevresult data.PipelineResult, i int) (data.PipelineResult, error) {
  var err error
  retresult := prevresult

  // Exec time
  start_time := time.Now()

  // Get var
  varcontent, err := getVarValue(prevresult, execstep.Command)
  if err != nil {
    goto stepEnd
  }

  // Assign
  err = setVarValue(&retresult, execstep.Output, varcontent)
  if err != nil {
    goto stepEnd
  }

  stepEnd:
  retresult.ExecStep[i].ExecTime = time.Since(start_time)
  retresult.ExecStep[i].Command = varcontent
  if err != nil {
    retresult.ExecStep[i].Error = err.Error()
  }
  return retresult, err
}

func execStepExecAssign(execstep data.PipelineExecution, prevresult data.PipelineResult, i int) (data.PipelineResult, error) {
  var err error
  retresult := prevresult

  // Exec time
  start_time := time.Now()

  // Replace values
  commandexec := cmdReplaceVars(prevresult, execstep.Command)

  // Execute command
  strOut, err :=  execCommandOutput(commandexec)
  if err != nil {
    goto execEnd
  }

  // Assign
  err = setVarValue(&retresult, execstep.Output, strOut)
  if err != nil {
    goto execEnd
  }



  execEnd:
  retresult.ExecStep[i].ExecTime = time.Since(start_time)
  retresult.ExecStep[i].Command = commandexec
  if err != nil {
    retresult.ExecStep[i].Error = err.Error()
  }

  return retresult, err
}

func execStepExec(execstep data.PipelineExecution, prevresult data.PipelineResult, i int) (data.PipelineResult, error) {
  var err error
  retresult := prevresult

  // Exec time
  start_time := time.Now()

  // Replace values
  commandexec := cmdReplaceVars(prevresult, execstep.Command)

  // Execute command
  err = execCommand(commandexec)
  if err != nil {
    goto execEnd
  }


  execEnd:
  retresult.ExecStep[i].ExecTime = time.Since(start_time)
  retresult.ExecStep[i].Command = commandexec
  if err != nil {
    retresult.ExecStep[i].Error = err.Error()
  }

  return retresult, err
}

func getPipelineOutput(piperesult data.PipelineResult, pipeline data.Pipeline) (data.PipelineResult, error) {
  ret_pipe := piperesult
  var err error

  if !pipeline.Output.Defined {
    return ret_pipe, err
  }

  if val, ok := piperesult.Variables[pipeline.Output.Value]; ok {
    ret_pipe.Output = val
    return ret_pipe, err
  }

  return ret_pipe, errors.New("Return variable not found")
}


func initVariables(pipeline data.Pipeline) (data.PipelineResult) {
  pipeline_ret := data.PipelineResult{}
  pipeline_ret.Variables = make(map[string]string)

  for _, val := range pipeline.Input {
    pipeline_ret.Variables[val.Name] = val.Value
  }

  for key, val := range pipeline.Declaration {
    pipeline_ret.Variables[key] = val
  }

  for i := 0; i < len(pipeline.Execution); i++ {
    stepnew := data.PipelineResultExecStep{}
    pipeline_ret.ExecStep = append(pipeline_ret.ExecStep, stepnew)
  }

  return pipeline_ret
}


func setLogRedirection(filename string) (error) {
  var err error
  loggerfile, err = os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0766)
  if err != nil {
    return err
  }

  log.SetOutput(loggerfile)

  return err
}
