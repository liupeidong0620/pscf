NAME    := pscf
VERSION := 0.1.0

LDFLAGS := -extldflags "-static"
GOBUILD := go build -a -ldflags

default:
	@ echo "no default target for Makefile"

clean:
	@ rm -rf $(NAME) ./releases

build-linux: clean
	@ CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) '$(LDFLAGS)' -o releases/$(NAME)

build-darwin: clean
	@ CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GOBUILD) "$(LDFLAGS)" -o releases/$(NAME)

build-windows: clean
	@ CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GOBUILD) "$(LDFLAGS)" -o releases/$(NAME).exe
