.PHONY: build test lint clean

VERSION := $(shell grep "\[VERSION\]" -A 1 METADATA | awk 'NR==2')
WEBSITE := $(shell grep "\[WEBSITE\]" -A 1 METADATA | awk 'NR==2')
REPOSITORY := $(shell grep "\[REPOSITORY\]" -A 1 METADATA | awk 'NR==2')
CURRENT_DATETIME := $(shell date +%Y-%m-%d\ %H:%M:%S)
LICENSE := $(shell head -n 1 LICENSE)

INTERNAL_DIRS := $(shell find ./internal -mindepth 1 -maxdepth 1 -type d -not -name components -not -name menus -printf './internal/%f/... ')

# Public targets
build: clean lint .change-package-json-version .build .build-linux .build-windows .build-debian .build-rpm .build-flatpak

test:
	@echo "\033[33m[Make]\033[0m \033[32mRunning tests...\033[0m"
	@go install gotest.tools/gotestsum@v1.12.3
	@env TEST_ENV=true gotestsum ./cmd/commands/... $(INTERNAL_DIRS) -- -count=1
	@echo "\033[33m[Make]\033[0m \033[32mTests completed.\033[0m"

test_cov:
	@echo "\033[33m[Make]\033[0m \033[32mRunning tests with coverage...\033[0m"
	@go install gotest.tools/gotestsum@v1.12.3
	@env TEST_ENV=true gotestsum ./cmd/commands/... $(INTERNAL_DIRS) -coverprofile=coverage.out
	@echo "\033[33m[Make]\033[0m \033[32mTests with coverage completed.\033[0m"

test_cov_ui:
	@echo "\033[33m[Make]\033[0m \033[32mRunning tests with coverage...\033[0m"
	@go install gotest.tools/gotestsum@v1.12.3
	@env TEST_ENV=true gotestsum ./cmd/commands/... $(INTERNAL_DIRS) -coverprofile=coverage.out && go tool cover -html=coverage.out
	@echo "\033[33m[Make]\033[0m \033[32mTests with coverage completed.\033[0m"

lint:
	@echo "\033[33m[Make]\033[0m \033[32mRunning linter...\033[0m"
	@find . -type f -name '*.go' -not -path './vendor/*' | xargs gofmt -w
	@echo "\033[33m[Make]\033[0m \033[32mLinters finished.\033[0m"

clean: .debian-clean
	@echo "\033[33m[Make]\033[0m \033[32mCleaning up vendor folder...\033[0m"
	@rm -rf ./vendor
	@echo "\033[33m[Make]\033[0m \033[32mCleaning up build folder...\033[0m"
	@rm -rf ./build	
	@rm -rf ./.flatpak-builder
	@echo "\033[33m[Make]\033[0m \033[32mCleaned.\033[0m"

# Internal targets

.change-package-json-version:
	@echo "\033[33m[Make]\033[0m \033[32mUpdating version of package.json...\033[0m"
	@if [ -f ./docs/package.json ]; then \
	sed -i 's/\("version" *: *\)"[^"]*"/\1"$(VERSION)"/' ./docs/package.json; \
	echo "\033[33m[Make]\033[0m \033[32mpackage.json version updated to $(VERSION).\033[0m"; \
	else \
	echo "\033[33m[Make]\033[0m \033[31mpackage.json not found in ./docs.\033[0m"; \
	fi

.build:
	@echo "\033[33m[Make]\033[0m \033[32mBuilding...\033[0m"
	@mkdir -p ./build

.build-linux:
	@echo "\033[33m[Make]\033[0m \033[32mBuilding Linux binary...\033[0m"
	@GOOS=linux GOARCH=amd64 go build \
		-ldflags="-X 'github.com/diogopereiradev/httpzen/cmd/commands/version.Version=$(VERSION)' \
		-X 'github.com/diogopereiradev/httpzen/cmd/commands/version.BuildDate=$(CURRENT_DATETIME)' \
		-X 'github.com/diogopereiradev/httpzen/cmd/commands/version.Website=$(WEBSITE)' \
		-X 'github.com/diogopereiradev/httpzen/cmd/commands/version.Repository=$(REPOSITORY)' \
		-X 'github.com/diogopereiradev/httpzen/cmd/commands/version.License=$(LICENSE)'" \
		-o ./build/httpzen main.go
	@echo "\033[33m[Make]\033[0m \033[32mLinux binary build finished.\033[0m"

.build-windows:
	@echo "\033[33m[Make]\033[0m \033[32mBuilding Windows binary...\033[0m"
	@GOOS=windows GOARCH=amd64 go build \
		-ldflags="-X 'github.com/diogopereiradev/httpzen/cmd/commands/version.Version=$(VERSION)' \
		-X 'github.com/diogopereiradev/httpzen/cmd/commands/version.BuildDate=$(CURRENT_DATETIME)' \
		-X 'github.com/diogopereiradev/httpzen/cmd/commands/version.Website=$(WEBSITE)' \
		-X 'github.com/diogopereiradev/httpzen/cmd/commands/version.Repository=$(REPOSITORY)' \
		-X 'github.com/diogopereiradev/httpzen/cmd/commands/version.License=$(LICENSE)'" \
		-o ./build/httpzen.exe main.go
	@echo "\033[33m[Make]\033[0m \033[32mWindows binary build finished.\033[0m"

.debian-clean:
	@echo "\033[33m[Make]\033[0m \033[32mCleaning up Debian build...\033[0m"
	@rm -rf ./build/debian

.build-debian:
	@echo "\033[33m[Make]\033[0m \033[32mBuilding Debian package...\033[0m"
	@mkdir -p ./build
	@rm -rf ./build/debian
	@mkdir -p ./build/debian/.cache/usr/bin
	@echo "\033[33m[Make]\033[0m \033[32mUpdating Debian control file version to $(VERSION)...\033[0m"
	@sed -i "s/^Version: .*/Version: $(VERSION)/" ./pkgroot/DEBIAN/control
	@cp ./build/httpzen ./build/debian/.cache/usr/bin/httpzen
	@cp -r ./pkgroot/DEBIAN ./build/debian/.cache/DEBIAN
	@dpkg-deb --build ./build/debian/.cache ./build/debian/httpzen.deb
	@echo "\033[33m[Make]\033[0m \033[32mDebian package build ended.\033[0m"

.build-rpm:
	@echo "\033[33m[Make]\033[0m \033[32mBuilding RPM package...\033[0m"
	@mkdir -p ./build/rpmroot/usr/bin
	@cp ./build/httpzen ./build/rpmroot/usr/bin/httpzen
	@fpm -s dir -t rpm -n httpzen -v $(VERSION) -C ./build/rpmroot -p ./build/httpzen.rpm usr/bin/httpzen || echo "FPM not found, install using: gem install --no-document fpm"
	@echo "\033[33m[Make]\033[0m \033[32mRPM package build ended.\033[0m"

.build-flatpak:
	@echo "\033[33m[Make]\033[0m \033[32mBuilding Flatpak package...\033[0m"
	@mkdir -p ./build/flatpak
	flatpak-builder --force-clean --repo=./build/flatpak/repo ./build/flatpak flatpak-manifest.yaml --default-branch=$(VERSION)
	flatpak build-bundle ./build/flatpak/repo ./build/flatpak/httpzen.flatpak github.diogopereiradev.httpzen $(VERSION)
	@echo "\033[33m[Make]\033[0m \033[32mFlatpak package build ended.\033[0m"