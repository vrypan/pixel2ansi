VERSION := $(shell git describe --tags 2>/dev/null || echo "v0.0.0")

BINS = pixel2ansi
SOURCES := $(wildcard *.go */*.go go.mod)

# Colors for output
GREEN = \033[0;32m
NC = \033[0m

all: $(BINS)

clean:
	@echo -e "$(GREEN)Deleting ${BINS} binary...$(NC)"
	rm -f $(BINS)

.PHONY: all clean local tag tag-minor tag-major releases

$(BINS): $(SOURCES)
	@echo -e "$(GREEN)Building ${BINS} ${VERSION} $(NC)"
	go build -o $@ -ldflags "-w -s -X main.VERSION=${VERSION}"

tag:
	./bin/auto_increment_tag.sh patch

tag-minor:
	./bin/auto_increment_tag.sh minor

tag-major:
	./bin/auto_increment_tag.sh major

releases:
	goreleaser release --clean
