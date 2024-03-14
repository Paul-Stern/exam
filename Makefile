win:
	GOOS=windows GOARCH=amd64 go build -o webtest.exe .
default:
	go build -ldflags "-X main.version=`git rev-parse HEAD`".

clean:
	rm webtest.exe webtest
