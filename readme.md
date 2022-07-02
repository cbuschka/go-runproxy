# go-runproxy - google cloud run app proxy

### Background

When cloud run (container service on google cloud platform) starts a new service revision it grants computation
resources for initialization. These resources will we revoked when the maximum startup time is reached or the service
opens a receiving socket. In case when the service opens the socket before startup has finished (as it is the case with
spring boot) the service is throttled and cannot finish startup anymore.

runproxy is a small tcp proxy that checks the service if it has come up yet (e.g. via a magic line appearing in the 
log output) and subsequently opens a socket to signal cloud run that the container is ready.

```bash
runproxy \ 
  # listen on 0.0.0.0:8080 and forward tcp conns to 127.0.0.1:8181
  -l 0.0.0.0:8080 \
  -d 127.0.0.1:8181 \
  # assume the app has been started if the following pattern matches a line in the stdout output
  -m '^.*Started application in.*seconds.*$' \
  # the command line to launch the downstream app
  -- java -jar app.jar
```

## Hints
* Make sure the output channel of your service is unbuffered. In case of python use -u flag.

## License

Copyright (c) 2022 by [Cornelius Buschka](https://github.com/cbuschka).

[Apache License, Version 2.0](./license.txt)
