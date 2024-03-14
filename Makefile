win:
	GOOS=windows GOARCH=amd64 go build --ldflags "-X main.version=`git describe --tags`" -o webtest.exe .
default:
	go build -ldflags "-X main.version=`git describe --tags`".

clean:
	rm webtest.exe webtest
