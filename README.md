# Go Microservice Boilerplate
A starter project to build multiple microservice in the same repo with Golang, Echo, Mysql and Redis.

## Directory structure

```bash
|
|____internal
|    |
|    |____configs
|    |    |____configs.go
|    |
|    |____api
|    |    |____notes.go
|    |    |____users.go
|    |
|    |____users
|    |    |____store.go
|    |    |____cache.go
|    |    |____users.go
|    |    |____users_test.go
|    |
|    |____notes
|    |    |____notes.go
|    |
|    |____pkg
|    |    |____stringutils
|    |    |____datastore
|    |    |     |____datastore.go
|    |    |____cachestore
|    |    |    |____cachestore.go
|    |    |____logger
|    |         |____logger.go
|    |
|    |____server
|         |____http
|         |    |____web
|         |    |    |____templates
|         |    |         |____index.html
|         |    |____handlers_notes.go
|         |    |____handlers_users.go
|         |    |____http.go
|         |
|         |____grpc
|
|____lib
|    |____notes
|         |____notes.go
|
|____vendor
|
|____docker
|    |____Dockerfile # your 'default' dockerfile
|
|____go.mod
|____go.sum
|
|____ciconfig.yml # depends on the CI/CD system you're using. e.g. .travis.yml
|____README.md
|____main.go
|
```


