# build image
FROM golang:1.22 as build-stage

WORKDIR /app

COPY go.mod go.sum ./ 
RUN go mod download && go mod verify

COPY cmd cmd
COPY pkg pkg
RUN CGO_ENABLED=0 go build ./cmd/goob.go

# deploy image
FROM alpine:3.14 as deploy-stage
WORKDIR /app

COPY --from=build-stage /app/goob goob

ENTRYPOINT [ "./goob", "-port", "5000" ]
# ENTRYPOINT ["ls", "-laR"]
