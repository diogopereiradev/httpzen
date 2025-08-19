.PHONY: build test lint clean sign-windows verify-signature
.ONESHELL:

VERSION := $(shell grep "\[VERSION\]" -A 1 METADATA | awk 'NR==2')
WEBSITE := $(shell grep "\[WEBSITE\]" -A 1 METADATA | awk 'NR==2')
REPOSITORY := $(shell grep "\[REPOSITORY\]" -A 1 METADATA | awk 'NR==2')
CURRENT_DATETIME := $(shell date +%Y-%m-%d\ %H:%M:%S)
LICENSE := $(shell head -n 1 LICENSE)

INTERNAL_DIRS := $(shell find ./internal -mindepth 1 -maxdepth 1 -type d -not -name components -not -name menus -printf './internal/%f/... ')

# Public targets
build: clean lint .change-package-json-version .build .build-linux .build-windows .build-debian .build-rpm .build-flatpak

# Sign the Windows binary (optional, requires certificate)
# Usage examples:
#  make sign-windows SIGN_PFX=</path/cert.pfx> SIGN_PFX_PASS=<password>
#  make sign-windows SIGN_CERT=</path/cert.pem> SIGN_KEY=</path/key.pem> SIGN_KEY_PASS=<optional>
# Optional metadata overrides: SIGN_NAME and SIGN_URL
sign-windows: .build-windows .sign-windows

# Verify Authenticode signature on the built exe
verify-signature:
	@which osslsigncode >/dev/null 2>&1 || { echo "\033[33m[Make]\033[0m \033[31mosslsigncode not found. Install it (e.g., apt install osslsigncode).\033[0m"; exit 1; }
	@[ -f ./build/httpzen.exe ] || { echo "\033[33m[Make]\033[0m \033[31mFile ./build/httpzen.exe not found. Build first.\033[0m"; exit 1; }
	@echo "\033[33m[Make]\033[0m \033[32mVerifying signature of httpzen.exe...\033[0m"
	@osslsigncode verify -in ./build/httpzen.exe || { echo "\033[33m[Make]\033[0m \033[31mSignature verification failed.\033[0m"; exit 1; }
	@echo "\033[33m[Make]\033[0m \033[32mSignature verified.\033[0m"

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
	@rm -f ./rsrc_windows_*.syso
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
	@$(MAKE) .winres
	@GOOS=windows GOARCH=amd64 go build \
		-ldflags="-X 'github.com/diogopereiradev/httpzen/cmd/commands/version.Version=$(VERSION)' \
		-X 'github.com/diogopereiradev/httpzen/cmd/commands/version.BuildDate=$(CURRENT_DATETIME)' \
		-X 'github.com/diogopereiradev/httpzen/cmd/commands/version.Website=$(WEBSITE)' \
		-X 'github.com/diogopereiradev/httpzen/cmd/commands/version.Repository=$(REPOSITORY)' \
		-X 'github.com/diogopereiradev/httpzen/cmd/commands/version.License=$(LICENSE)'" \
		-o ./build/httpzen.exe main.go
	@echo "\033[33m[Make]\033[0m \033[32mWindows binary build finished.\033[0m"

# Internal: sign Windows binary with osslsigncode (PFX or CERT/KEY)
# Inputs (env vars):
#  SIGN_PFX, SIGN_PFX_PASS  - OR -  SIGN_CERT, SIGN_KEY, SIGN_KEY_PASS
#  SIGN_NAME (default: httpzen), SIGN_URL (default: $(WEBSITE))
#  SIGN_TSA_URL (default: http://timestamp.digicert.com)
.sign-windows:
	@which osslsigncode >/dev/null 2>&1 || { echo "\033[33m[Make]\033[0m \033[31mosslsigncode not found. Install it (e.g., apt install osslsigncode).\033[0m"; exit 1; }
	@[ -f ./build/httpzen.exe ] || { echo "\033[33m[Make]\033[0m \033[31mFile ./build/httpzen.exe not found. Build first.\033[0m"; exit 1; }
	@SIGN_NAME_EFF=$${SIGN_NAME:-httpzen}; \
	SIGN_URL_EFF=$${SIGN_URL:-"$(WEBSITE)"}; \
	SIGN_TSA_EFF=$${SIGN_TSA_URL:-http://timestamp.digicert.com}; \
	OUT=./build/httpzen.exe.signed; \
	set -e; \
	echo "\033[33m[Make]\033[0m \033[32mSigning Windows binary with osslsigncode...\033[0m"; \
	if [ -n "$$SIGN_PFX" ]; then \
		[ -f "$$SIGN_PFX" ] || { echo "\033[33m[Make]\033[0m \033[31mSIGN_PFX not found: $$SIGN_PFX\033[0m"; exit 1; }; \
		[ -n "$$SIGN_PFX_PASS" ] || { echo "\033[33m[Make]\033[0m \033[31mSIGN_PFX_PASS not set.\033[0m"; exit 1; }; \
		osslsigncode sign -h sha256 -pkcs12 "$$SIGN_PFX" -pass "$$SIGN_PFX_PASS" -n "$$SIGN_NAME_EFF" -i "$$SIGN_URL_EFF" -t "$$SIGN_TSA_EFF" -in ./build/httpzen.exe -out "$$OUT"; \
	elif [ -n "$$SIGN_CERT" ] && [ -n "$$SIGN_KEY" ]; then \
		[ -f "$$SIGN_CERT" ] || { echo "\033[33m[Make]\033[0m \033[31mSIGN_CERT not found: $$SIGN_CERT\033[0m"; exit 1; }; \
		[ -f "$$SIGN_KEY" ] || { echo "\033[33m[Make]\033[0m \033[31mSIGN_KEY not found: $$SIGN_KEY\033[0m"; exit 1; }; \
		if [ -n "$$SIGN_KEY_PASS" ]; then \
			osslsigncode sign -h sha256 -certs "$$SIGN_CERT" -key "$$SIGN_KEY" -pass "$$SIGN_KEY_PASS" -n "$$SIGN_NAME_EFF" -i "$$SIGN_URL_EFF" -t "$$SIGN_TSA_EFF" -in ./build/httpzen.exe -out "$$OUT"; \
		else \
			osslsigncode sign -h sha256 -certs "$$SIGN_CERT" -key "$$SIGN_KEY" -n "$$SIGN_NAME_EFF" -i "$$SIGN_URL_EFF" -t "$$SIGN_TSA_EFF" -in ./build/httpzen.exe -out "$$OUT"; \
		fi; \
	else \
		echo "\033[33m[Make]\033[0m \033[31mNo certificate specified. Provide SIGN_PFX (and SIGN_PFX_PASS) or SIGN_CERT and SIGN_KEY.\033[0m"; exit 1; \
	fi; \
	osslsigncode verify -in "$$OUT" >/dev/null 2>&1 && mv -f "$$OUT" ./build/httpzen.exe || { echo "\033[33m[Make]\033[0m \033[31mSignature verification failed.\033[0m"; rm -f "$$OUT"; exit 1; }

.winres:
	@echo "\033[33m[Make]\033[0m \033[32mPreparing Windows resources (version info, manifest, icon)...\033[0m"
	@go install github.com/tc-hib/go-winres@v0.3.3 || go install github.com/tc-hib/go-winres@latest
	@mkdir -p ./winres
	@echo "\033[33m[Make]\033[0m \033[32mGenerating winres.json...\033[0m"
	@cat > ./winres/winres.json <<-'EOF'
	{
		"RT_VERSION": {
			"#1": {
				"0000": {
					"fixed": {
						"file_version": "$(VERSION).0",
						"product_version": "$(VERSION).0"
					},
					"info": {
						"0409": {
							"CompanyName": "diogopereiradev",
							"FileDescription": "httpzen - HTTP client TUI/CLI",
							"FileVersion": "$(VERSION).0",
							"OriginalFilename": "httpzen.exe",
							"ProductName": "httpzen",
							"ProductVersion": "$(VERSION)",
							"Comments": "$(WEBSITE)",
							"LegalCopyright": "$(LICENSE)"
						}
					}
				}
			}
		},
		"RT_MANIFEST": {
			"#1": {
				"0409": {
					"execution-level": "as invoker",
					"dpi-awareness": "per monitor v2",
					"minimum-os": "win7",
					"use-common-controls-v6": true
				}
			}
		},
		"RT_GROUP_ICON": {
			"APP": {
				"0000": "../docs/public/favicons/favicon-96x96.png"
			}
		}
	}
	EOF
	@echo "\033[33m[Make]\033[0m \033[32mEmbedding Windows resources with go-winres...\033[0m"
	@go-winres make
	@echo "\033[33m[Make]\033[0m \033[32mWindows resources ready.\033[0m"

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