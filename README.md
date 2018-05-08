# go-yarn
A Go (golang) library to manage client dependencies using yarn registry

It can be used as a library or as a bin command.

## Installation

As a library

```shell
go get github.com/rgobbo/go-yarn
```

or if you want to use it as a bin command
```shell
go get github.com/rgobbo/go-yarn/cmd/go-yarn
```

## Usage

Add a configuration `yarn.json` file to your project:

```json
{
  "dependencies": [
    {
      "lib": "angular",
      "version": "1.6.4"
    },
    {
      "lib": "jquery",
      "version": "3.2.1"
    }
  ]
}
```

Then in your Go app you can do something like

```go
package main

import (
    goyarn "github.com/rgobbo/go-yarn"
    "log"
    "os"
)

func main() {
  configFile := "./yarn.json"
  destPath := "./static/vendor"
  err := 	err := goyarn.YarnInstall(configFile, destPath)
  if err != nil {
    log.Fatal("Error processing yarn :", err)
  }

  log.Println("Yarn processed successfully !!"


}
```

### Command Mode

Assuming you've installed the command as above and you've got `$GOPATH/bin` in your `$PATH`
There is a sample file (yarn.json) in sample folder.

```shell
go-yarn -c yarn.json -f ./static/vendor/

```
Where -c is the configuration file and -f is the path where the libs will be installed.



### Writing Config Files

The file yarn.json is a json file with array of objects containing lib name and version.

```json
{
  "dependencies": [
    {
      "lib": "angular",
      "version": "1.6.4"
    },
    {
      "lib": "jquery",
      "version": "3.2.1"
    }
  ]
}
```

After the file was processed, it will change adding lines like :

```json
{
      "lib": "angular",
      "version": "1.6.4",
      "resolved": "https://registry.npmjs.org/angular/-/angular-1.6.4.tgz"
}

```