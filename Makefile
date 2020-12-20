.PHONY: build
build:
	env GOOS=linux go build -ldflags="-s -w" -o bin/handler handler/main.go

.PHONY: clean
clean:
	rm -rf ./bin

.PHONY: deploy
deploy: clean build
	sls deploy --verbose

.PHONY: test
test:
	sls invoke -f getPresence

.PHONY: remove
remove:
	sls remove


GO_FILES = $(shell go list -f '{{.Dir}}' ./...)

.PHONY: fmtcheck
fmtcheck:
	@command -v goimports > /dev/null 2>&1 || go get golang.org/x/tools/cmd/goimports
	@CHANGES="$$(goimports -d $(GO_FILES))"; \
		if [ -n "$${CHANGES}" ]; then \
			echo "Unformatted :\n\n$${CHANGES}\n\n"; \
			exit 1; \
		else  \
			echo "All good! (goimports)"; \
			exit 0;  \
		fi
	@# Annoyingly, goimports does not support the simplify flag.
	@CHANGES="$$(gofmt -s -d $(GO_FILES))"; \
		if [ -n "$${CHANGES}" ]; then \
			echo "Unformatted :\n\n$${CHANGES}\n\n"; \
			exit 1; \
		else \
			echo "All good! (gofmt)"; \
			exit 0;\
		fi

.PHONY: staticcheck
staticcheck:
	@command -v staticcheck > /dev/null 2>&1 || go get honnef.co/go/tools/cmd/staticcheck
	@staticcheck -checks="all" -tests $(GO_FILES)

.PHONY: spellcheck
spellcheck:
	@command -v misspell > /dev/null 2>&1 || go get github.com/client9/misspell/cmd/misspell
	@misspell -locale="UK" -error -source="text" $(GO_FILES)
