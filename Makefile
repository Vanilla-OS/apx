PREFIX=/usr/
DESTDIR=/
BINARY_NAME=apx

all: build

build:
	go build -o ${BINARY_NAME}

install:
	install -Dm755 ${BINARY_NAME} ${DESTDIR}${PREFIX}/bin/${BINARY_NAME}
	mkdir -p ${DESTDIR}/etc/apx
	sed -i 's|/usr/share/apx/distrobox|${PREFIX}/share/apx/distrobox|g' config/config.json
	install -Dm644 config/config.json ${DESTDIR}/etc/apx/config.json
	mkdir -p ${DESTDIR}${PREFIX}/share/apx
	sh distrobox/install --prefix ${DESTDIR}${PREFIX}/share/apx
	mv ${DESTDIR}${PREFIX}/share/apx/bin/distrobox* ${DESTDIR}${PREFIX}/share/apx/.

clean:
	rm -f ${BINARY_NAME}
	go clean
