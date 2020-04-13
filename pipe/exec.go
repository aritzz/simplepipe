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
  "strings"
  "bytes"
  "os/exec"
  "os"
)


func execCommandOutput(command string) (string, error) {
  commandWithArgs := strings.Split(command, " ")

	var out bytes.Buffer
  cmd := exec.Command(commandWithArgs[0], commandWithArgs[1:]...)
  cmd.Env = os.Environ()
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return out.String(), err
	}

  return strings.TrimSuffix(out.String(), "\n"), nil
}

func execCommand(command string) (error) {
  commandWithArgs := strings.Split(command, " ")

  cmd := exec.Command(commandWithArgs[0], commandWithArgs[1:]...)
  cmd.Env = os.Environ()

	err := cmd.Run()
	if err != nil {
		return err
	}

  return nil
}
