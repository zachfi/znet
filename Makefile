BIN_NAME=znet
PKG_SITE=gopkg.larch.space

build:
	go get -v
	go build -o build/$(BIN_NAME) -v

clean:
	rm -f build/

publish:
	scp -o StrictHostKeyChecking=no $(BIN_NAME) webdeploy@${PKG_SITE}:/usr/local/www/${PKG_SITE}/
