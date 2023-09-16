PREFIX=/usr/
DESTDIR=/
BINARY_NAME=apx

all: build

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -tags netgo -ldflags '-w -extldflags "-static"' -o ${BINARY_NAME}

install:
	install -Dm755 ${BINARY_NAME} ${DESTDIR}${PREFIX}/bin/${BINARY_NAME}
	sudo mkdir -p ${DESTDIR}/etc/apx
	sed -i 's|/usr/share/apx/distrobox|${PREFIX}/share/apx/distrobox|g' config/apx.json
	sudo install -Dm644 config/apx.json ${DESTDIR}/etc/apx/apx.json
	mkdir -p ${DESTDIR}${PREFIX}/share/apx/distrobox
	sh distrobox/install --prefix ${DESTDIR}${PREFIX}/share/apx/distrobox
	mv ${DESTDIR}${PREFIX}/share/apx/distrobox/bin/distrobox* ${DESTDIR}${PREFIX}/share/apx/distrobox/.

install-manpages:
	mkdir -p ${DESTDIR}${PREFIX}/share/man/man1
	cp -r man/* ${DESTDIR}${PREFIX}/share/man/.
	chmod 644 ${DESTDIR}${PREFIX}/share/man/man1/apx*

uninstall:
	sudo rm ${DESTDIR}${PREFIX}/bin/apx
	sudo rm -rf ${DESTDIR}/etc/apx
	sudo rm -rf ${DESTDIR}${PREFIX}/share/apx

uninstall-manpages:
	sudo rm -rf ${DESTDIR}${PREFIX}/share/man/man1/apx*

clean:
	rm -f ${BINARY_NAME}
	go clean
