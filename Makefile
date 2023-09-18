install:
	go mod tidy
	go mod vendor
	git submodule init && git submodule update
	git submodule add -f https://github.com/alex-miller-0/safe-global-encoder vendor/safe-global-encoder
	cd vendor/safe-global-encoder && git pull origin main && make install && make build
	cd ../..
	make build
	cp vendor/safe-global-encoder/encoder bin
build: 
	go build -o ./bin/manager