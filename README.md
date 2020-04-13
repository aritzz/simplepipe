# Simplepipe

Simplepipe is a simple software to write your own pipelines. You can write a sequence of commands that will be executed in order.


## Installation and usage

You must build your binary in order to run this software. Get a copy of this repository and build it yourself.

`git clone github.com/aritzz/simplepipe`

`go build github.com/aritzz/simplepipe`

Please, use `./simplepipe -h` to get more information.

## Writing a pipeline

A pipeline has two sections: declaration (where you can define variables) and command execution (where you provide commands to execute).

### Declaration

Declaration will start with *pipeline* word, followed by the name you want to use for this pipeline. Then, you can declare three types of variables:

- Simple declaration (*use variablename*): Declares a variable.
- Reader declaration (*read variablename*): Declares a variable that will be readed as argument.
- Random declaration (*rand variablename*): Declares random variable.

This declared variables can be used in command execution as *$varname*.

### Command execution

This section will start with the word *begin*. After that, you can use three type of executions:

- Assignation (*variable1 = variable2*): Simple assignation.
- Assignation with execution (*variable1 = (command to execute)*): Executes a command and assigns the output to a variable.
- Execution (*(command to execute)*): Executes a command.

You can finish command execution file with *end*. If you want to return a variable, you can use *end varname*.


## Examples

See the *examples/* directory on this repository. Execution examples:

`./simplepipe -pipeline examples/test.pipe -logfile test.log "John Doe"`

`./simplepipe -pipeline examples/test.pipe -logfile test.log -outputonly "John Doe"`

`./simplepipe -pipeline examples/transcode.pipe -logfile transcode.log "audio.wav" "audio.mp3"`


## License

Simplepipe is released under GNU General Public License. For more details, take a look at the [LICENSE](https://github.com/aritzz/simplepipe/blob/master/LICENSE)
