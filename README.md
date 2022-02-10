## Local Run

```sh
go run ./app.go

# generate swagger
go get github.com/swaggo/swag/cmd/swag
Swag init -g app.go
```

## Test with ginkgo

```sh
cd handlers
# generate test suite for handlers directory
# package name is handlers_test
ginkgo boostrap
ginkgo generate health_check.go

# under that folder
# run test suite
ginkgo --v
# run only regexp matched "Metadata" in any described clause ex: Describe/It
ginkgo -v --focus=Metadata ./handlers
# run only regexp matched "without input payload" in any described clause ex: Describe/It
ginkgo -v --focus="without input payload" ./handlers

ginkgo -v ./handlers
```

## Local Docker Compose Run

```sh
# rebuild docker image for fiber app
docker-compose build

docker-compose up

docker-compose down
```

# Introduction

APIs:

```sh
# bump version when release
/version

/health

# get templating or other app info, ex sign-in templating
/api/ethereum-auth/v1/metadata

# generate and cache nonce on redis then respond to user
/api/ethereum-auth/v1/nonce

# verify signature and auth system account binding etc.
/api/ethereum-auth/v1/login
```

# Getting Started

TODO: Guide users through getting your code up and running on their own system. In this section you can talk about:
1.	Installation process
2.	Software dependencies
3.	Latest releases
4.	API references

# Build and Test
TODO: Describe and show how to build your code and run the tests.

# Contribute
TODO: Explain how other users and developers can contribute to make your code better.

If you want to learn more about creating good readme files then refer the following [guidelines](https://docs.microsoft.com/en-us/azure/devops/repos/git/create-a-readme?view=azure-devops). You can also seek inspiration from the below readme files:
- [ASP.NET Core](https://github.com/aspnet/Home)
- [Visual Studio Code](https://github.com/Microsoft/vscode)
- [Chakra Core](https://github.com/Microsoft/ChakraCore)