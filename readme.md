# Third Law

A small monitoring tool meant to stay out of your way. Specifically, it doesn't have to own a process to be able to work with it.

## Example

```json
{
  "frequency": 2,
  "outputs": {
    "success": [{
      "type": "file",
      "truncate": true,
      "path": "last.log"
    }],
    "failure": [
      {
        "type": "file",
        "truncate": true,
        "path": "last.log"
      },
      {
        "type": "stderr",
        "snooze": 120
      }
    ]
  },
  "checks": [{
    "name": "redis",
    "type": "shell",
    "command": "redis-cli",
    "dir": "/opt/redis/",
    "arguments": ["ping"],
    "out": "PONG\n",
    "recover": ["restart redis"]
  }],
  "actions": {
    "restart redis": {
      "type": "shell",
      "command": "/etc/init.d/redis",
      "arguments": ["start"]
    }
  }
}
```

The `include` directive can be used to specify a directory which contains additional configuration files. These files are limited to the `checks`, `check` and `actions` fields. A check in one file can reference an action in a different file.

# Running Flags

- `-config=config.json`: path to the config file
- `-test`: to test the validity of specified config. An exit code of 0 means the file was valid.

# Configuration Options

## frequency
The time, in seconds, to run the checks. Checks are run in series, so the real "sleep" time is going to be the specified frequency + however long it takes to do all the checks and actions.

Defaults to 10 seconds.

## include
Specifies a directory to load additional configuration files from. All files within the directory are loaded. These child files are limited to defining `checks`, `check` and `actions`

## outputs
Defines the outputs to send the results of an iteration to. Outputs can be sent for the case where all checks pass (`success`) or when at least one check fails (`failure`):

```json
{
  "outputs": {
    "success": [{"type": "stdout"}],
    "failure": [
      {"type": "stderr"},
      {
        "type": "file",
        "path": "/var/log/opt/thirdlaw.log",
        "truncate": true
      }
    ]
  }
}
```

Every output supports:

* `snooze`: time, in seconds, to wait before re-using this output. Snoozing allows you to keep a short frequency without worrying about flooding your output (it's particularly useful with the http output).
* `disabled`: disables the output (defaults to false)


### stdout
Writes the results to stdout

### stderr
Writes the results to stderr

### file
Writes the results to a specified file. The file output accepts the following configuration values:

- `path`: the filepath to write results to (defaults to failures.log)
- `truncate`: whether or not to truncate the file before each write (defaults to false)

### http
Sends the response to an HTTP endpoint. The http output accepts the following configuration values:

- `address`: REQUIRED full http path (scheme, host, port, url) (no default)
- `body`: An optional body to send. This can either be a string, or a nested object. In the case of a nested object / array, the body will be converted to JSON. (defaults to empty)

This will issue a GET if `body` is empty. A POST is made otherwise.

It's possible to use the $FRIENDLY$ placeholder within the body in order to generate a simple but friendly error message.

An example of using the HTTP output with OpsGenie:

```json
"failure": [
  {
    "type": "http",
    "snooze": 300,
    "address": "https://api.opsgenie.com/v1/json/alert",
    "body": {
      "apiKey": "YOUR OPS GENIE API KEY",
      "message": "$FRIENDLY$",
      "recipients": "YOUR_ALERT_EMAIL",
      "alias": "health check"
    }
  }
]
```

## checks and check
`checks` and `check` define the code to execute on each iteration. The two fields are only different in that `checks` is an array of `check`.

All `checks` accept:

- `recover`: an array of `action` names to run in case of failure (defaults to none)
- `snooze`: time in second to wait before running the check again, essentially allowing checks to run at custom frequencies (defaults to 0, meaning the global frequency is used).

### http
Makes an HTTP request. Any error or a response with a status code of 300 or more will result in a failure. The http check accepts the following configuration values:

- `address`: the full address (scheme, host, port, path) to make the request to (defaults to http://127.0.0.1/)
- `timeout`: the timeout, in seconds, to wait before getting a response (defaults to 5).
- `contains`: text that must exist in the response body to consider the check a success (defaults to none)

### shell
Invokes the shell and runs the specified command. Any error running the command, include an exit code not equal to 0, will result in a failure. It's also possible to specify the expected output. The shell check accepts the following configuration values:

- `command`: REQUIRED command to run (no default)
- `arguments`: array of arguments to pass to the command (defaults to none)
- `dir`: the worker directory to use (defaults to thirdlaw's working directory)
- `out`: expected stdout text (defaults to ignoring any output)

## actions
Actions are invoked when a check fails

- `retries`: how many times to retry the action should it fail (defaults to 0)
- `delay`: how long to wait, in seconds, between failed retries (detauls to 1)

Remember, if you're frequence is set to 10 and you have a single check with a single action that has a retry of 5 and a delay of 1, the worst case check frequency is 15 seconds (10 + 1 * 5) as everything happens synchronously.

### shell
Invokes the shell and runs the specified command. The shell action accepts the following configuration values:

- `command`: REQUIRED command to run (no default)
- `arguments`: array of arguments to pass to the command (defaults to none)
- `dir`: the worker directory to use (defaults to thirdlaw's working directory)
