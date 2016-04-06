# Guestimator

## Development Setup

Examples assume mac env + homebrew.

Install and start postgres.
After installing, `brew info postgres` should get postgres started.

`brew install postgres`

Setup env-specific database and roles.

`psql -f db/setup.sql postgres`

Install [goose](https://bitbucket.org/liamstask/goose/) which is used to manage db migrations.

`go get bitbucket.org/liamstask/goose/cmd/goose`

Apply all migrations.

`goose up`
`goose -env test up`

Use [godep](https://github.com/tools/godep) to manage dependencies.

`go get github.com/tools/godep`

Fetch go dependencies (vendored).

`go get ./...`

Install testing deps and check if your environment works by

`go get github.com/onsi/ginkgo/ginkgo`
`go get github.com/onsi/gomega`
`ginkgo -r`
