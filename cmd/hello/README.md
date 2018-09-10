# ochello

Hello world app for OpenCensus proof of concept.

This app will report metrics to prometheus and a trace to jaeger when run in
docker compose with those components.

# Build

(The following is performed in the ochello directory)

This will build a docker container image for ochello and call it ochello:
```
$ ./dockerbuild.sh
```

# Run

This will run the ochello docker container in isolation:
```
$ docker run ochello
hello
```
