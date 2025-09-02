#
# This file is part of bspent.
#
# bspent is free software: you can redistribute it and/or modify it
# under the terms of the GNU General Public License as published by the
# Free Software Foundation, either version 3 of the License, or (at
# your option) any later version.
#
# bspent is distributed in the hope that it will be useful, but WITHOUT
# ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or
# FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License
# for more details.
#
# You should have received a copy of the GNU General Public License
# along with Foobar. If not, see <https://www.gnu.org/licenses/>.
#

# Define variables
BIN_DIR := bin
APP_NAME := bspent
MAIN_GO := cmd/bspent.go

SRC := \
	bsp/bspfile.go \
	cmd/bspent.go \
	ent/entfile.go

# Default target
.PHONY: all
all: build

# Build the go application
build: $(BIN_DIR)/$(APP_NAME)

windows: $(BIN_DIR)/$(APP_NAME).exe

$(BIN_DIR)/$(APP_NAME): $(SRC)
	@echo "Building $(APP_NAME)..."
	go build -o $@ $(MAIN_GO)

$(BIN_DIR)/$(APP_NAME).exe: $(SRC)
	@echo "Building $(APP_NAME) for Windows..."
	GOOS=windows GOARCH=amd64 go build -o $@ $(MAIN_GO)

# Run the go application
.PHONY: run
run: build
	@echo "Running $(APP_NAME)..."
	./$(BIN_DIR)/$(APP_NAME)

# Run tests
.PHONY: test
test:
	@echo "Running tests..."
	go test -v ./...

# Clean up generated files
.PHONY: clean
clean:
	@echo "Cleaning up..."
	rm -rf $(BIN_DIR)

# Install Go dependencies
.PHONY: deps
deps:
	@echo "Installing dependencies..."
	go mod tidy

# Display help message
.PHONY: help
help:
	@echo "Usage: make <target>"
	@echo ""
	@echo "Targets:"
	@echo "  all     : Default target, builds the application."
	@echo "  build   : Builds the Go application."
	@echo "  run     : Builds and runs the Go application."
	@echo "  test    : Runs all Go tests."
	@echo "  clean   : Removes generated binary files."
	@echo "  deps    : Installs Go module dependencies."
	@echo "  help    : Displays this help message."
