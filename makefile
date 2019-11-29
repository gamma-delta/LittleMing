# Note: I have no clue what I'm doing.
# if you wanted to cross-compile this, you probably know what you're doing already
# and don't need a makefile...

windows : 
	GOOS=windows GOARCH=amd64 go build -v -o='little_ming.exe'

mac :
	GOOS=darwin GOARCH=amd64 go build -v -o='little_ming_mac'

linux :
	GOOS=linux GOARCH=amd64 go build -v -o='little_ming_linux'