lint: ## Lint the source files
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.46.2
	golangci-lint run ./...

OUTPUT_DIR=build/test-results
COVERAGE_FILE=cover.out
JUNIT_FILE=test-report.xml
TIMEOUT=10m
test: ## Run tests.
	# Run tests with coverage
	go run github.com/onsi/ginkgo/v2/ginkgo -r \
	--procs=2 \
	--compilers=2 \
	--randomize-all \
	--randomize-suites \
	--fail-on-pending \
	--keep-going \
	--trace \
	--cover \
	--coverprofile=$(COVERAGE_FILE) \
	--junit-report=$(JUNIT_FILE) \
	--output-dir=$(OUTPUT_DIR) \
	--timeout=$(TIMEOUT) $(TEST_PATHS) \
	./...
