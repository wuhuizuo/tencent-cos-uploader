TARGET 		= tencent-cos-uploader
VERSION 	= `git describe --tags`
BUILD_DATE 	= `date +%F_%T_%z`

# Setup the -ldflags option for go build here, interpolate the variable values
LDFLAGS=-ldflags "-w -s -X main.version=${VERSION} -X main.buildDate=${BUILD_DATE}"

.PHONY: default
default: clean ${TARGET}

${TARGET}:
	go build ${LDFLAGS} -o ${TARGET} -v

clean: 
	rm -vf $(TARGET)