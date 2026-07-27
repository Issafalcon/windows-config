.PHONY: run build test windows dist

run:
	go run .

build:
	go build -o windows-config-tui$(shell go env GOEXE) .

windows:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o windows-config-tui.exe .

# Local release-shaped layout (exe + ElevatedHelper), same as CI zip contents.
dist:
	mkdir -p ../dist/windows-amd64/tools
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -trimpath -ldflags="-s -w" \
		-o ../dist/windows-amd64/windows-config-tui.exe .
	cp ../tools/ElevatedHelper.ps1 ../dist/windows-amd64/tools/

test:
	go test ./...
