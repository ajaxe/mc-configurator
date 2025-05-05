SHELL := pwsh.exe

clean:
	@echo "Cleaning up..."
	@pwsh -Command "if(Test-Path -Path ./tmp) { rm -recurse -force ./tmp }"

build:
	@echo "Building CLI"
	@pwsh -Command "Set-Item Env:GOARCH amd64; Set-Item Env:GOOS windows; go build -o ./tmp/mc-config.exe ."
	@pwsh -Command "copy ./server.properties ./tmp/"

linux:
	@echo "Building CLI"
	@pwsh -Command "Set-Item Env:GOARCH amd64; Set-Item Env:GOOS linux; go build -o ./tmp/linux/mc-config ."
	@pwsh -Command "copy ./server.properties ./tmp/"