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
.PHONY: all build clean test test-all test-dashboard test-chart test-coverage test-context test-integration test-vim test-short fix-tests fmt lint run server install vet

all: test-dashboard build-all

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

# Run all tests including integration and vim mode tests
test-all: test-chart test-vim
	@echo ""
	@echo "╭───────────────────────────────────────────────────────────────╮"
	@echo "│             ✅  All Tests Completed  ✅                        │"
	@echo "╰───────────────────────────────────────────────────────────────╯"

# Comprehensive test dashboard showing all test results
test-dashboard:
	@echo "╭───────────────────────────────────────────────────────────────╮"
	@echo "│           🎵  AlgoScales Complete Test Suite  🎵              │"
	@echo "╰───────────────────────────────────────────────────────────────╯"
	@echo ""
	@echo "Running all tests..."
	@echo ""
	@echo "════════════════════════════════════════════════════════════════"
	@echo "                    📦 Package Tests                             "
	@echo "════════════════════════════════════════════════════════════════"
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
		END { \
			print "\n"; \
			printf "Package Summary: %d passed, %d failed, %d no tests\n", passed, failed, notests; \
			printf "Test Cases: %d passed, %d failed\n", testspassed, testsfailed; \
		} \
	'
	@echo ""
	@echo "════════════════════════════════════════════════════════════════"
	@echo "                    🎹 Vim Mode Tests                            "
	@echo "════════════════════════════════════════════════════════════════"
	@$(GOTEST) -v ./cmd -run TestVimModeIntegration -timeout 30s 2>&1 | grep -E "(PASS|FAIL|ok|---)" | awk ' \
		/--- PASS:/ { printf "✅ %s\n", $$0 } \
		/--- FAIL:/ { printf "❌ %s\n", $$0 } \
		/^PASS$$/ { printf "\033[32m✅ Vim Integration Tests: PASSED\033[0m\n" } \
		/^FAIL$$/ { printf "\033[31m❌ Vim Integration Tests: FAILED\033[0m\n" } \
	'
	@echo ""
	@$(GOTEST) -v ./cmd -run TestVimCommands 2>&1 | grep -E "(PASS|FAIL|ok|---)" | awk ' \
		/--- PASS:/ { printf "✅ %s\n", $$0 } \
		/--- FAIL:/ { printf "❌ %s\n", $$0 } \
		/^PASS$$/ { printf "\033[32m✅ Vim Commands Tests: PASSED\033[0m\n" } \
		/^FAIL$$/ { printf "\033[31m❌ Vim Commands Tests: FAILED\033[0m\n" } \
	'
	@echo ""
	@$(GOTEST) -v ./cmd -run TestMultiLevelHints 2>&1 | grep -E "(PASS|FAIL|ok|---)" | awk ' \
		/--- PASS:/ { printf "✅ %s\n", $$0 } \
		/--- FAIL:/ { printf "❌ %s\n", $$0 } \
		/^PASS$$/ { printf "\033[32m✅ Multi-level Hints Tests: PASSED\033[0m\n" } \
		/^FAIL$$/ { printf "\033[31m❌ Multi-level Hints Tests: FAILED\033[0m\n" } \
	'
	@echo ""
	@echo "════════════════════════════════════════════════════════════════"
	@echo "                    🤖 AI Assistant Tests                        "
	@echo "════════════════════════════════════════════════════════════════"
	@$(GOTEST) -v ./internal/ai 2>&1 | grep -E "(PASS|FAIL|ok|---)" | awk ' \
		/--- PASS:/ { printf "✅ %s\n", $$0 } \
		/--- FAIL:/ { printf "❌ %s\n", $$0 } \
		/^ok/ { printf "\033[32m✅ AI Assistant Tests: PASSED\033[0m\n" } \
		/^FAIL/ { printf "\033[31m❌ AI Assistant Tests: FAILED\033[0m\n" } \
	'
	@echo ""
	@echo "════════════════════════════════════════════════════════════════"
	@echo "                   🔄 Context Integration                        "
	@echo "════════════════════════════════════════════════════════════════"
	@$(GOTEST) -v ./internal/problem ./internal/stats ./internal/registry ./internal/services ./internal/session ./internal/ai 2>&1 | \
		grep -E "(ok|FAIL)" | awk ' \
		/^ok/ { printf "✅ \033[32m%-50s\033[0m %s\n", $$2, "PASS" } \
		/^FAIL/ { printf "❌ \033[31m%-50s\033[0m %s\n", $$2, "FAIL" } \
	'
	@echo ""
	@echo "════════════════════════════════════════════════════════════════"
	@echo "                    🔧 Manual Tests                              "
	@echo "════════════════════════════════════════════════════════════════"
	@if [ -f ./tests/manual_vim_test.sh ]; then \
		./tests/manual_vim_test.sh 2>&1 | grep -E "(Testing:|PASS|FAIL)" | awk ' \
			/PASS/ { printf "✅ %s\n", $$0 } \
			/FAIL/ { printf "❌ %s\n", $$0 } \
			/Testing:/ && !/Testing complete/ { printf "🧪 %s\n", $$0 } \
		'; \
	else \
		echo "⚠️  Manual test script not found"; \
	fi
	@echo ""
	@echo "════════════════════════════════════════════════════════════════"
	@echo "                    📊 Final Summary                             "
	@echo "════════════════════════════════════════════════════════════════"
	@cat /tmp/test-output.txt | awk ' \
		BEGIN { packages=0; passed=0; failed=0; tests=0; testspassed=0; testsfailed=0; } \
		/^ok / { packages++; passed++; } \
		/^FAIL/ && !/\[no test files\]/ { packages++; failed++; } \
		/--- PASS:/ { tests++; testspassed++; } \
		/--- FAIL:/ { tests++; testsfailed++; } \
		END { \
			printf "📦 Total Packages Tested: %d\n", packages; \
			printf "   ✅ Passed: %d (%d%%)\n", passed, (packages>0 ? passed*100/packages : 0); \
			printf "   ❌ Failed: %d (%d%%)\n", failed, (packages>0 ? failed*100/packages : 0); \
			printf "\n"; \
			printf "🧪 Total Test Cases: %d\n", tests; \
			printf "   ✅ Passed: %d (%d%%)\n", testspassed, (tests>0 ? testspassed*100/tests : 0); \
			printf "   ❌ Failed: %d (%d%%)\n", testsfailed, (tests>0 ? testsfailed*100/tests : 0); \
			printf "\n"; \
			if (failed == 0 && testsfailed == 0) { \
				printf "\033[32m🎉 All tests passed! Ready to build.\033[0m\n"; \
			} else { \
				printf "\033[31m⚠️  Some tests failed. Please fix before proceeding.\033[0m\n"; \
				exit 1; \
			} \
		} \
	'
	@rm -f /tmp/test-output.txt
	@echo ""
	@echo "╭───────────────────────────────────────────────────────────────╮"
	@echo "│                    Test Suite Complete                         │"
	@echo "╰───────────────────────────────────────────────────────────────╯"

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
	$(GOTEST) -v ./internal/problem ./internal/stats ./internal/registry ./internal/services ./internal/session ./internal/ai
	@echo ""
	@echo "✅ Context integration tests completed!"

test-integration:
	@echo "╭───────────────────────────────────────────────────────────────╮"
	@echo "│             🧪  Running Integration Tests  🧪                  │"
	@echo "╰───────────────────────────────────────────────────────────────╯"
	@echo ""
	@echo "Running standard integration tests..."
	$(GOTEST) -v -tags=integration $(PKG)
	@echo ""
	@echo "Running vim mode integration tests..."
	$(GOTEST) -v ./cmd -run TestVimModeIntegration -timeout 30s

test-vim:
	@echo "╭───────────────────────────────────────────────────────────────╮"
	@echo "│             🎹  Testing Vim Mode Integration  🎹               │"
	@echo "╰───────────────────────────────────────────────────────────────╯"
	@echo ""
	@echo "Building binary..."
	@$(GOBUILD) -o $(BIN_DIR)/$(BINARY_NAME) -v $(MAIN_PATH)
	@echo ""
	@echo "Running vim mode integration tests..."
	@$(GOTEST) -v ./cmd -run TestVimModeIntegration -timeout 30s
	@echo ""
	@echo "Running vim commands tests..."
	@$(GOTEST) -v ./cmd -run TestVimCommands
	@echo ""
	@echo "Running multi-level hints tests..."
	@$(GOTEST) -v ./cmd -run TestMultiLevelHints
	@echo ""
	@if [ -f ./tests/manual_vim_test.sh ]; then \
		echo "Running manual vim tests..."; \
		./tests/manual_vim_test.sh; \
	fi
	@echo ""
	@echo "✅ Vim mode tests completed!"

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
	@mkdir -p $(BIN_DIR)
	@GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BIN_DIR)/$(BINARY_NAME)-linux-amd64 $(MAIN_PATH)
	@GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BIN_DIR)/$(SERVER_BINARY_NAME)-linux-amd64 $(SERVER_PATH)

build-windows:
	@mkdir -p $(BIN_DIR)
	@GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(BIN_DIR)/$(BINARY_NAME)-windows-amd64.exe $(MAIN_PATH)
	@GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(BIN_DIR)/$(SERVER_BINARY_NAME)-windows-amd64.exe $(SERVER_PATH)

build-darwin:
	@mkdir -p $(BIN_DIR)
	@GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(BIN_DIR)/$(BINARY_NAME)-darwin-amd64 $(MAIN_PATH)
	@GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(BIN_DIR)/$(SERVER_BINARY_NAME)-darwin-amd64 $(SERVER_PATH)
	@GOOS=darwin GOARCH=arm64 $(GOBUILD) -o $(BIN_DIR)/$(BINARY_NAME)-darwin-arm64 $(MAIN_PATH)
	@GOOS=darwin GOARCH=arm64 $(GOBUILD) -o $(BIN_DIR)/$(SERVER_BINARY_NAME)-darwin-arm64 $(SERVER_PATH)

# Build all platforms
.PHONY: build-all
build-all: 
	@echo ""
	@echo "╭───────────────────────────────────────────────────────────────╮"
	@echo "│              🔨  Building All Platforms  🔨                    │"
	@echo "╰───────────────────────────────────────────────────────────────╯"
	@echo ""
	@echo "Building Linux binaries..."
	@$(MAKE) -s build-linux
	@echo "✅ Linux builds complete"
	@echo ""
	@echo "Building Windows binaries..."
	@$(MAKE) -s build-windows
	@echo "✅ Windows builds complete"
	@echo ""
	@echo "Building Darwin (macOS) binaries..."
	@$(MAKE) -s build-darwin
	@echo "✅ Darwin builds complete"
	@echo ""
	@echo "╭───────────────────────────────────────────────────────────────╮"
	@echo "│         🎉  All Builds Complete!  🎉                           │"
	@echo "│                                                                │"
	@echo "│  Binaries available in: bin/                                   │"
	@echo "│                                                                │"
	@echo "│  Linux:   algo-scales-linux-amd64                             │"
	@echo "│  Windows: algo-scales-windows-amd64.exe                        │"
	@echo "│  macOS:   algo-scales-darwin-amd64 (Intel)                    │"
	@echo "│           algo-scales-darwin-arm64 (Apple Silicon)            │"
	@echo "╰───────────────────────────────────────────────────────────────╯"
