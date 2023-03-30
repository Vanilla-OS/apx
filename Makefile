PREFIX=/usr/
DESTDIR=/
BINARY_NAME=apx

all: build

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -tags netgo -ldflags '-w -extldflags "-static"' -o ${BINARY_NAME}

install:
	install -Dm755 ${BINARY_NAME} ${DESTDIR}${PREFIX}/bin/${BINARY_NAME}
	sudo mkdir -p ${DESTDIR}/etc/apx
	sed -i 's|/usr/share/apx/distrobox|${PREFIX}/share/apx/distrobox|g' config/config.json
	sudo install -Dm644 config/config.json ${DESTDIR}/etc/apx/config.json
	mkdir -p ${DESTDIR}${PREFIX}/share/apx
	sh distrobox/install --prefix ${DESTDIR}${PREFIX}/share/apx
	mv ${DESTDIR}${PREFIX}/share/apx/bin/distrobox* ${DESTDIR}${PREFIX}/share/apx/.

install-manpages:
	mkdir -p ${DESTDIR}${PREFIX}/share/man/man1
	cp -r man/* ${DESTDIR}${PREFIX}/share/man/.
	chmod 644 ${DESTDIR}${PREFIX}/share/man/man1/apx*
	chmod 644 ${DESTDIR}${PREFIX}/share/man/*/man1/apx*

uninstall:
	sudo rm ${DESTDIR}${PREFIX}/bin/apx
	sudo rm -rf ${DESTDIR}/etc/apx
	sudo rm -rf ${DESTDIR}${PREFIX}/share/apx

uninstall-manpages:
	sudo rm -rf ${DESTDIR}${PREFIX}/share/man/man1/apx*
	sudo rm -rf ${DESTDIR}${PREFIX}/share/man/*/man1/apx*

clean:
	rm -f ${BINARY_NAME}
	go clean
