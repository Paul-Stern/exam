win:
	GOOS=windows GOARCH=amd64 go build -o webtest.exe .
default:
	go build .

clean:
	rm webtest.exe webtest
