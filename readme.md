# go-runproxy - google cloud run app proxy

### Background

When cloud run (container service on google cloud platform) starts a new service revision it grants computation
resources for initialization. These resource will we revoked when the maximum startup time is reached or the service
opens a receiving socket. In case when the service opens the socket before startup has finished (as it is the case with
spring boot) the service is throttled and cannot finish startup anymore.

runproxy is a small http proxy that checks the service if it has come up yet (i.e. via the actuator health check
endpoint) and subsequently opens a socket to signal cloud run that the container is ready.

## License

Copyright (c) 2022 by [Cornelius Buschka](https://github.com/cbuschka).

[Apache License, Version 2.0](./license.txt)
