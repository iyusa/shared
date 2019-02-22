lib:
	go get -u github.com/gotestyourself/gotest.tools
	go get -u github.com/nu7hatch/gouuid
	go get -u gopkg.in/ini.v1

build:
	@echo "Tidak ada"	

test:
	go test ./...	