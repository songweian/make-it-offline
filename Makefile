BINARY := make-it-offline
PKG := ./cmd/cli
OUT := dist

# Default platforms to build. Format: os_arch (use underscore, not slash)
PLATFORMS ?= \
  linux_amd64 \
  linux_arm64 \
  darwin_amd64 \
  darwin_arm64 \
  windows_amd64 \
  windows_arm64

# Disable CGO for portable cross compilation
export CGO_ENABLED := 0

.PHONY: help build cross clean list

help:
	@echo "Targets:"
	@echo "  build        Build for current platform into $(OUT)/"
	@echo "  cross        Cross-compile for multiple platforms into $(OUT)/"
	@echo "  list         Print the target platforms"
	@echo "  clean        Remove $(OUT)/"

$(OUT):
	@mkdir -p $(OUT)

build: $(OUT)
	GO111MODULE=on go build -o $(OUT)/$(BINARY) $(PKG)

list:
	@printf "%s\n" $(PLATFORMS)

# Cross-compile matrix build
cross: $(OUT)
	@set -e; \
	for p in $(PLATFORMS); do \
	  OS=$$(echo $$p | cut -d _ -f1); \
	  ARCH=$$(echo $$p | cut -d _ -f2); \
	  EXT=""; [ "$$OS" = "windows" ] && EXT=".exe"; \
	  OUTFILE="$(OUT)/$(BINARY)-$$OS-$$ARCH$$EXT"; \
	  echo "Building $$OUTFILE"; \
	  GOOS=$$OS GOARCH=$$ARCH GO111MODULE=on go build -ldflags "-s -w" -o "$$OUTFILE" $(PKG); \
	done

clean:
	rm -rf $(OUT)
