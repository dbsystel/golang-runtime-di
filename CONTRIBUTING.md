# Contribution Guidelines

## General

Before you push a PR, please create an [issue](../../issues).

That way we can manage and track the changes being made to the project.

Make sure your code has been linted successfully, and none of the tests fail.

We would kindly ask you to not decrease the code coverage, so please write/adapt tests accordingly.

## Communication

We do not tolerate violent/racist/sexist or any other behavior aiming to harm anyone, 
so please respect each other as human beings (see [code of conduct](./CODE_OF_CONDUCT.MD)). Thank you.

## Linting

To lint the project locally, please run `make lint`.
This will use the linter settings from [here](.golangci.yml).  

## Test coverage

To test with coverage, please run `make test`.
A coverage report will be generated and the total coverage will be shown.

## Releases

Releases are being made by tagging a commit on `main` with a version.
