# build image
FROM golang:1.22 AS build-stage

# install go templ
RUN go env -w GOBIN=/usr/bin
RUN go install github.com/a-h/templ/cmd/templ@v0.2.680

# install tailwind
WORKDIR /usr/bin
RUN curl -sLO https://github.com/tailwindlabs/tailwindcss/releases/download/v3.4.3/tailwindcss-linux-x64
RUN mv tailwindcss-linux-x64 tailwindcss
RUN chmod +x tailwindcss

# build goob
WORKDIR /app

# install deps
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# build templates
COPY pkg pkg
RUN templ generate

# build executable
COPY cmd cmd
RUN CGO_ENABLED=0 go build ./cmd/goob.go

# build css 
COPY tailwind.config.js ./
COPY public/input.css public/input.css
RUN tailwindcss -i public/input.css -o public/index.css

# deploy image
FROM alpine:3.14 AS deploy-stage
WORKDIR /app

COPY --from=build-stage /app/goob goob
COPY --from=build-stage /app/public/index.css public/index.css

ENTRYPOINT [ "./goob", "-port", "5000" ]
# ENTRYPOINT ["ls", "-laR"]
