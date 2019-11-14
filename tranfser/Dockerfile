FROM golang AS build-env
ADD .. /
RUN go build -o /app /cmd/main.go

FROM ubuntu:14.04

COPY --from=build-env /app /app
RUN apt-get update && apt-get install -y ca-certificates
ADD agent.yml /
ADD forwarder.yml /
ADD bench.yml /
ENTRYPOINT ["/app"]