
# build stage
FROM golang:alpine AS build-env
RUN apk --no-cache add build-base git bzr mercurial gcc
ADD . /src
RUN cd /src && go mod download && go build -o go-todo-api

# final stage
FROM alpine
WORKDIR /app
COPY --from=build-env /src/go-todo-api /app/
EXPOSE 8080
ENTRYPOINT ./go-todo-api