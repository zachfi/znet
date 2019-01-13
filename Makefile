BIN_NAME=znet
PKG_SITE=gopkg.larch.space

build:
	go get -v
	GOOS=linux go build -o bin/$(BIN_NAME) -v
	docker build . -t znet

clean:
	rm -f build/

publish:
	scp -o StrictHostKeyChecking=no $(BIN_NAME) webdeploy@${PKG_SITE}:/usr/local/www/${PKG_SITE}/

