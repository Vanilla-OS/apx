PREFIX := /usr
DESTDIR := /
BINARY_NAME := apx

GO := go

all: clean build

build: ${BINARY_NAME}

${BINARY_NAME}:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 ${GO} build -a -tags netgo -ldflags '-w -extldflags "-static"' -o $@

install: build
	install -Dm755 ${BINARY_NAME} ${DESTDIR}${PREFIX}/bin/${BINARY_NAME}
	mkdir -p ${DESTDIR}/etc/apx
	sed -i 's|/usr/share/apx/distrobox|${PREFIX}/share/apx/distrobox|g' config/apx.json
	install -Dm644 config/apx.json ${DESTDIR}/etc/apx/apx.json
	mkdir -p ${DESTDIR}${PREFIX}/share/apx/distrobox
	sh distrobox/install --prefix ${DESTDIR}${PREFIX}/share/apx/distrobox
	mv ${DESTDIR}${PREFIX}/share/apx/distrobox/bin/distrobox* ${DESTDIR}${PREFIX}/share/apx/distrobox/.

install-manpages:
	mkdir -p ${DESTDIR}${PREFIX}/share/man/man1
	cp -r man/* ${DESTDIR}${PREFIX}/share/man/.
	chmod 644 ${DESTDIR}${PREFIX}/share/man/man1/apx*

uninstall: uninstall-manpages
	rm ${DESTDIR}${PREFIX}/bin/${BINARY_NAME}
	rm -rf ${DESTDIR}/etc/apx
	rm -rf ${DESTDIR}${PREFIX}/share/apx

uninstall-manpages:
	rm -rf ${DESTDIR}${PREFIX}/share/man/man1/apx*

clean:
	rm -f ${BINARY_NAME}
	${GO} clean
