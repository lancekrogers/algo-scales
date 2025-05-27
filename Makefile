# Build automation and commands

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOFMT=$(GOCMD) fmt

# Default test package path
PKG?=./...

# Binary names and paths
BINARY_NAME=algo-scales
SERVER_BINARY_NAME=algo-scales-server
BIN_DIR=bin

# Main package paths
MAIN_PATH=./
SERVER_PATH=./server

# Targets
.PHONY: all build clean test test-chart test-coverage test-context test-integration test-short fix-tests fmt lint run server install vet

all: test-chart build

build:
	mkdir -p $(BIN_DIR)
	$(GOBUILD) -o $(BIN_DIR)/$(BINARY_NAME) -v $(MAIN_PATH)
	$(GOBUILD) -o $(BIN_DIR)/$(SERVER_BINARY_NAME) -v $(SERVER_PATH)

clean:
	$(GOCLEAN)
	rm -f $(BIN_DIR)/$(BINARY_NAME)
	rm -f $(BIN_DIR)/$(SERVER_BINARY_NAME)
	rm -f $(BIN_DIR)/$(BINARY_NAME).test
	rm -f $(BIN_DIR)/$(SERVER_BINARY_NAME).test
	rm -f coverage.out

test:
	$(GOTEST) -v $(PKG)

test-chart:
	@echo "╭───────────────────────────────────────────────────────────────╮"
	@echo "│                 🎵  AlgoScales Test Results  🎵                │"
	@echo "╰───────────────────────────────────────────────────────────────╯"
	@echo ""
	@$(GOTEST) -v $(PKG) 2>&1 | tee /tmp/test-output.txt | awk ' \
		BEGIN { passed=0; failed=0; buildfail=0; notests=0; testspassed=0; testsfailed=0; } \
		/^ok / { passed += 1; printf "✅ \033[32m%-50s\033[0m %s\n", $$2, "PASS" } \
		/^FAIL.*/ { if (match($$0, /\[build failed\]/)) { \
						buildfail += 1; printf "🔨 \033[33m%-50s\033[0m %s\n", $$2, "BUILD FAILED" \
					} else { \
						failed += 1; printf "❌ \033[31m%-50s\033[0m %s\n", $$2, "FAIL" \
					} \
		} \
		/\?\s+/ { notests += 1; printf "🔍 \033[36m%-50s\033[0m %s\n", $$2, "NO TESTS" } \
		/--- PASS:/ { testspassed += 1 } \
		/--- FAIL:/ { testsfailed += 1 } \
		END { total = passed + failed + buildfail + notests; \
			printf "\n📊 \033[1mTest Summary:\033[0m\n"; \
			printf "   Total Packages: %d\n", total; \
			printf "   ✅ Passed:      %d (%d%%)\n", passed, (total>0 ? passed*100/total : 0); \
			printf "   ❌ Failed:      %d (%d%%)\n", failed, (total>0 ? failed*100/total : 0); \
			printf "   🔨 Build Failed: %d (%d%%)\n", buildfail, (total>0 ? buildfail*100/total : 0); \
			printf "   🔍 No Tests:     %d (%d%%)\n", notests, (total>0 ? notests*100/total : 0); \
			printf "\n🧪 \033[1mTest Cases:\033[0m\n"; \
			printf "   Total Tests:   %d\n", testspassed + testsfailed; \
			printf "   ✅ Passed:      %d\n", testspassed; \
			if (testsfailed > 0) { \
				printf "   ❌ Failed:      %d\n", testsfailed; \
			} \
		} \
	'
	@echo ""
	@echo "╭───────────────────────────────────────────────────────────────╮"
	@echo "│ Legend: ✅ Passed  ❌ Failed  🔨 Build Failed  🔍 No Tests     │"
	@echo "╰───────────────────────────────────────────────────────────────╯"
	@rm -f /tmp/test-output.txt
	@if [ $$? -ne 0 ]; then \
		echo ""; \
		echo "⚠️  Some tests are failing! Run 'make test' for details."; \
		exit 1; \
	fi

test-coverage:
	$(GOTEST) -coverprofile=coverage.out $(PKG)
	$(GOCMD) tool cover -html=coverage.out

test-context:
	@echo "╭───────────────────────────────────────────────────────────────╮"
	@echo "│             🔄  Testing Context Integration  🔄                │"
	@echo "╰───────────────────────────────────────────────────────────────╯"
	@echo ""
	@echo "Testing packages with context.Context integration..."
	@echo ""
	$(GOTEST) -v ./internal/problem ./internal/stats ./internal/registry ./internal/services ./internal/session
	@echo ""
	@echo "✅ Context integration tests completed!"

test-integration:
	@echo "╭───────────────────────────────────────────────────────────────╮"
	@echo "│             🧪  Running Integration Tests  🧪                  │"
	@echo "╰───────────────────────────────────────────────────────────────╯"
	@echo ""
	$(GOTEST) -v -tags=integration $(PKG)

test-short:
	@echo "╭───────────────────────────────────────────────────────────────╮"
	@echo "│             ⚡  Running Quick Tests  ⚡                        │"
	@echo "╰───────────────────────────────────────────────────────────────╯"
	@echo ""
	$(GOTEST) -short $(PKG)

fix-tests:
	@echo "╭───────────────────────────────────────────────────────────────╮"
	@echo "│             🔧  AlgoScales Test Fixing Guide  🔧               │"
	@echo "╰───────────────────────────────────────────────────────────────╯"
	@echo ""
	@echo "Common issues and how to fix them:"
	@echo ""
	@echo "1️⃣  \033[1mRedeclared functions in session package\033[0m"
	@echo "   ➡️  Make helper functions like selectProblem mockable variables:"
	@echo "   var selectProblem = func(pattern, difficulty string) (*problem.Problem, error) { ... }"
	@echo ""
	@echo "2️⃣  \033[1mUndefined functions in problem package\033[0m"
	@echo "   ➡️  Import problem package correctly and use mockable variables:"
	@echo "   import \"github.com/lancekrogers/algo-scales/internal/problem\""
	@echo ""
	@echo "3️⃣  \033[1mField and method conflicts in session package\033[0m"
	@echo "   ➡️  Rename fields to avoid conflicts with method names:"
	@echo "   Change ShowHints field to hintsShown"
	@echo ""
	@echo "4️⃣  \033[1mMissing pattern styles in view tests\033[0m"
	@echo "   ➡️  Add expected pattern styles in view_test.go"
	@echo ""
	@echo "5️⃣  \033[1mFailed stats tests\033[0m"
	@echo "   ➡️  Ensure mock data is properly initialized in stats_test.go"
	@echo ""
	@echo "6️⃣  \033[1mContext parameter missing errors\033[0m"
	@echo "   ➡️  All service/repository methods now require context.Context as first parameter"
	@echo "   ➡️  Update calls to pass context.Background() or appropriate context"
	@echo ""
	@echo "7️⃣  \033[1mMock implementations don't match interfaces\033[0m"
	@echo "   ➡️  Update mock methods to accept context.Context parameter"
	@echo "   ➡️  Check MockProblemRepository, MockStatsService, MockStorage"
	@echo ""
	@echo "Run specific test packages:"
	@echo "   make test PKG=./internal/session"
	@echo "   make test PKG=./internal/problem"
	@echo ""
	@echo "Test context integration:"
	@echo "   make test-context"
	@echo ""
	@echo "For more details on test failures:"
	@echo "   make test"

fmt:
	$(GOFMT) ./...

vet:
	$(GOCMD) vet ./...

lint:
	golangci-lint run

run:
	mkdir -p $(BIN_DIR)
	$(GOBUILD) -o $(BIN_DIR)/$(BINARY_NAME) -v $(MAIN_PATH)
	$(BIN_DIR)/$(BINARY_NAME)

server:
	mkdir -p $(BIN_DIR)
	$(GOBUILD) -o $(BIN_DIR)/$(SERVER_BINARY_NAME) -v $(SERVER_PATH)
	$(BIN_DIR)/$(SERVER_BINARY_NAME)

install:
	mkdir -p $(BIN_DIR)
	$(GOBUILD) -o $(BIN_DIR)/$(BINARY_NAME) -v $(MAIN_PATH)
	sudo mv $(BIN_DIR)/$(BINARY_NAME) /usr/local/bin/$(BINARY_NAME)

# Cross-compile targets
.PHONY: build-linux build-windows build-darwin

build-linux:
	mkdir -p $(BIN_DIR)
	GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BIN_DIR)/$(BINARY_NAME)-linux-amd64 -v $(MAIN_PATH)
	GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BIN_DIR)/$(SERVER_BINARY_NAME)-linux-amd64 -v $(SERVER_PATH)

build-windows:
	mkdir -p $(BIN_DIR)
	GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(BIN_DIR)/$(BINARY_NAME)-windows-amd64.exe -v $(MAIN_PATH)
	GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(BIN_DIR)/$(SERVER_BINARY_NAME)-windows-amd64.exe -v $(SERVER_PATH)

build-darwin:
	mkdir -p $(BIN_DIR)
	GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(BIN_DIR)/$(BINARY_NAME)-darwin-amd64 -v $(MAIN_PATH)
	GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(BIN_DIR)/$(SERVER_BINARY_NAME)-darwin-amd64 -v $(SERVER_PATH)
	GOOS=darwin GOARCH=arm64 $(GOBUILD) -o $(BIN_DIR)/$(BINARY_NAME)-darwin-arm64 -v $(MAIN_PATH)
	GOOS=darwin GOARCH=arm64 $(GOBUILD) -o $(BIN_DIR)/$(SERVER_BINARY_NAME)-darwin-arm64 -v $(SERVER_PATH)

# Build all platforms
.PHONY: build-all
build-all: build-linux build-windows build-darwin
