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
  "errors"
  "strings"
)


func getVarValue(pipeline data.PipelineResult, variable string) (string, error) {
  var err error
  var_ret, recv := pipeline.Variables[variable]

  if !recv {
    err = errors.New("Error getting variable "+variable)
  }

  return var_ret, err
}

func setVarValue(pipeline *data.PipelineResult, variable string, value string) (error) {
  _, exists := pipeline.Variables[variable]
  if !exists {
    return errors.New("Variable "+variable+" is not declared")
  }

  pipeline.Variables[variable] = value

  return nil
}

func cmdReplaceVars(pipeline data.PipelineResult, command string) (string) {
  replaced := command

  for key, value := range pipeline.Variables {
    currentVar := "$"+key
    replaced = strings.ReplaceAll(replaced, currentVar, value)
  }

  return replaced
}
