FROM golang:1.17-alpine AS build
WORKDIR /app
COPY . .
ARG CGO_ENABLED=0
RUN go build -o bin/server ./cmd/server

FROM alpine
WORKDIR /
COPY --from=build /app/bin/server /app-server
CMD [ "/app-server" ]
