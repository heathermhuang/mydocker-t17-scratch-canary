FROM golang@sha256:bff16a94c0b4b6d9642581e2a3aa81beb3edb940b2f5a3aac1d22c6a68cb7747 AS build
WORKDIR /src
COPY go.mod main.go ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GOPROXY=off GOSUMDB=off go build -buildvcs=false -trimpath -ldflags="-s -w" -o /probe .

FROM scratch
COPY --from=build /probe /probe
USER 65532:65532
EXPOSE 8080
ENTRYPOINT ["/probe"]
