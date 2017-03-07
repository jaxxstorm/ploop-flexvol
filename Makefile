SOURCES := $(shell find . 2>&1 | grep -E '.*\.(c|h|go)$$')

.DEFAULT: ploop

ploop: $(SOURCES)
	go build -o ploop .
