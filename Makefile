.SILENT: clean gohtags
.PHONY: clean
SRCS=$(wildcard *.go)

9ccgo: clean
	go build -gcflags '-N -l' -o gohtags $(SRCS)
	
clean:
	rm -f gohtags *.html

