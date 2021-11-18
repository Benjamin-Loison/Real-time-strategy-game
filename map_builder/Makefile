CC		= go build
EXEC	= map_builder
SRC		= $(wildcard src/*.go)

.PHONY: clean

$(EXEC): $(SRC)
	$(CC) -o $(EXEC) $^

clean:
	rm $(EXEC)
