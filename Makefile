# JOJO Discord Bot - An advanced multi-purpose discord bot
# Copyright (C) 2022 Lazy Bytez (Elias Knodel, Pascal Zarrad)
#
# This program is free software: you can redistribute it and/or modify
# it under the terms of the GNU Affero General Public License as published
# by the Free Software Foundation, either version 3 of the License, or
# (at your option) any later version.
#
# This program is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU Affero General Public License for more details.
#
# You should have received a copy of the GNU Affero General Public License
# along with this program.  If not, see <https://www.gnu.org/licenses/>.

# Get current directory
CURRENT_DIR:=$(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))

# === Project Setup ===
# =====================
# Copy the example env file
.PHONY: env
env:
	cp .env.example .env

# Install go dependencies
.PHONY: install
install: openapi/dependencies openapi/generate
	go get

.PHONY: setup
setup: env install

# === Running Project ===
# =======================
# Runs the code for development and test usage
.PHONY: run
run: openapi/generate
	go run .

# Builds an executable for production usage
# Build will be tagged as edge and include commit SHA
.PHONY: build
build: openapi/generate
	go build -ldflags "-X github.com/lazybytez/jojo-discord-bot/build.Version=edge -X github.com/lazybytez/jojo-discord-bot/build.CommitSHA=`git rev-parse --short HEAD`" .

# === Quality Assurance ===
# =========================
# Runs tests with specified arguments
.PHONY: test
test:
	go test -race -covermode=atomic -coverpkg=all ./...

# Lints the code
.PHONY: lint
lint:
	docker run --rm -v $(CURRENT_DIR):/app -w /app golangci/golangci-lint:v1.49.0 golangci-lint run -v

# === OpenAPI ===
# =========================
# Install CLI for documentation generation
.PHONY: openapi/dependencies
openapi/dependencies:
	go install github.com/swaggo/swag/cmd/swag@latest

# Generate OpenAPI documentation
.PHONY: openapi/generate
openapi/generate:
	swag init
