CLI=wapb
SRV=$(CLI)-server
SRCS=$(shell find . -type f -name '*.go')

ALL: bin/$(CLI) bin/$(SRV)

bin/$(CLI): bin/ $(SRCS)
	go build -o bin/ ./cmd/$(@F)

bin/$(SRV): bin/ $(SRCS) cmd/$(SRV)/assets.go
	go build -o bin/ ./cmd/$(@F)

bin/:
	mkdir -p $@


cmd/$(SRV)/assets.go: cmd/$(SRV)/assets_gen.go frontend/dist/index.html
	go generate ./cmd/$(SRV)


frontend/dist/index.html: frontend/node_modules $(shell find frontend -type f -name '*.vue') $(shell find frontend -type f -name '*.js')
	cd frontend && npm run build && npm run generate

frontend/node_modules:
	cd frontend && npm install

run:
	cd frontend && npm run dev

clean:
	$(RM) -rf bin/
	$(RM) -rf frontend/dist

.PHONY: clean run