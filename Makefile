serve:
	air

test:
	go test -v --cover ./...

PHONY: serve test