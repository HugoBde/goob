# build image
FROM golang:1.22 as build-stage

WORKDIR /usr/bin
RUN curl -sLO https://github.com/tailwindlabs/tailwindcss/releases/download/v3.4.3/tailwindcss-linux-x64
RUN mv tailwindcss-linux-x64 tailwindcss
RUN chmod +x tailwindcss

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY cmd cmd
COPY pkg pkg
RUN CGO_ENABLED=0 go build ./cmd/goob.go

COPY tailwind.config.js ./
COPY public/input.css public/input.css
RUN tailwindcss -i public/input.css -o public/index.css

# deploy image
FROM alpine:3.14 as deploy-stage
WORKDIR /app

COPY --from=build-stage /app/goob goob
COPY --from=build-stage /app/public/index.css public/index.css

ENTRYPOINT [ "./goob", "-port", "5000" ]
# ENTRYPOINT ["ls", "-laR"]
