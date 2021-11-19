CC		= go build
EXEC	= rts_server
SRC		= $(wildcard src/*.go)

.PHONY: clean run

$(EXEC): $(SRC)
	$(CC) -o $(EXEC) $^

run: $(EXEC)
	./$(EXEC)

clean:
	rm $(EXEC)
