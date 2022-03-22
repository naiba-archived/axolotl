FROM golang:alpine AS binarybuilder
RUN apk --no-cache --no-progress add \
    gcc \
    git \
    musl-dev
WORKDIR /axolotl
COPY . .
RUN cd cmd/api \
    && go build -o api -ldflags="-s -w"
FROM alpine:latest
RUN apk --no-cache --no-progress add \
    ca-certificates \
    tzdata
WORKDIR /axolotl
COPY dist /axolotl/dist
COPY --from=binarybuilder /axolotl/cmd/api/api ./api

VOLUME ["/axolotl/data"]
EXPOSE 80
CMD ["/axolotl/api"]