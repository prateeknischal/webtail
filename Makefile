VERSION = DEV
DIST = webtail-$(VERSION)

package: clean build
	mkdir $(DIST)
	cp webtail $(DIST)/
	cp -R static/ $(DIST)/
	cp -R templates/ $(DIST)/
	tar -cvzf $(DIST).tar.gz $(DIST)
	rm -rf $(DIST)

build:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -o webtail

clean:
	rm -rf webtail
	rm -rf $(DIST) $(DIST).tar.gz
