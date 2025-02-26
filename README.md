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

# Reference

1. https://npmdoc.github.io/node-npmdoc-web3/build/apidoc.html
2. https://www.toptal.com/ethereum-smart-contract/time-locked-wallet-truffle-tutorial
3. https://www.toptal.com/ethereum/one-click-login-flows-a-metamask-tutorial
4. https://eth.wiki/json-rpc/API#eth_sign
5. https://medium.com/@angellopozo/ethereum-signing-and-validating-13a2d7cb0ee3
6. https://eips.ethereum.org/EIPS/eip-191
7. https://javascript.hotexamples.com/examples/eth-lib.lib.hash/-/keccak256/javascript-keccak256-function-examples.html
9. https://geth.ethereum.org/docs/clef/apis
9. https://en.wikipedia.org/wiki/Elliptic_Curve_Digital_Signature_Algorithm
10. https://goethereumbook.org/signature-verify/
11. https://medium.com/mycrypto/the-magic-of-digital-signatures-on-ethereum-98fe184dc9c7

## ecrecovery id

since testRPC other protocols will add 27, we need to minus it for Geth go-ethereum.

- https://github.com/openethereum/parity-ethereum/issues/5490
- https://cloud.githubusercontent.com/assets/27297969/25291256/86ecd9ec-268e-11e7-98ec-fc339c4dbafe.png
- https://github.com/ChainSafe/web3.js/blob/0.20.7/DOCUMENTATION.md#web3ethsign
- https://eips.ethereum.org/EIPS/eip-1271
- https://docs.kaleido.io/faqs/why-ecrecover-fails/
- https://blog.ethereum.org/2019/07/10/geth-v1-9-0/
- https://bitcoin.stackexchange.com/questions/38351/ecdsa-v-r-s-what-is-v?newreg=4f3cd8cd47a74643b02cec1825ed9cc1

if protocol is Geth or Quorum => ecrecover id don't need to - 27
