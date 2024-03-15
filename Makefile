prod-win:
	WEBTEST_ENV=prod GOOS=windows GOARCH=amd64 go build --ldflags "-X main.version=`git describe --tags`" -o ./bin/webtest.exe .
win:
	GOOS=windows GOARCH=amd64 go build --ldflags "-X main.version=`git describe --tags`" -o ./bin/webtest.exe .
default:
	go build -ldflags "-X main.version=`git describe --tags`" -o ./bin/webtest .

clean:
	rm webtest.exe webtest
