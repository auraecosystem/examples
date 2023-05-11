
all: README.md
	go build ./...

README.md: go-types.md ../cmd/weave/*.go $(wildcard */*.go)
	go build -C .. ./cmd/weave
	../weave go-types.md >README.md
	rm ../weave

# This is for previewing using github.
# $HOME/markdown must be a github client.
test: README.md
	cp README.md $$HOME/markdown/
	(cd $$HOME/markdown/ && git commit -m . README.md && git push)

