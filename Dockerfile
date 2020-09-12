FROM golang:alpine AS binarybuilder
WORKDIR /helloengineer
COPY . .
RUN cd cmd/api \
    && go build -o api -ldflags="-s -w"
FROM alpine:latest
RUN apk --no-cache --no-progress add \
    ca-certificates \
    tzdata
WORKDIR /helloengineer
COPY dist /helloengineer/dist
COPY --from=binarybuilder /helloengineer/cmd/api/api ./api

VOLUME ["/helloengineer/data"]
EXPOSE 8080
CMD ["/helloengineer/api"]