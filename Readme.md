To get the gokit-base project up and running you'll need to have a few things installed beforehand:

Install Go
Install Consul
Before building the sample application you will need to bootstrap Consul with environment configuration and export necessary variables to initialize the app. The environment configuration example can be found in the docker/gokit-base/resources folder. The application's KV path where you will need to PUT to can be found in config/config.go.

Finally, run go build and ./gokit-base to start the listening server!

# Intro
Inspired by: https://github.com/bnelz/gokit-base

# Usage

```
go run .
```

# Services
1. String service: https://gokit.io/examples/stringsvc.html

2. Advanced: https://github.com/go-kit/examples/tree/master/shipping