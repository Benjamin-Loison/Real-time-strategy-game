CC		= go build
EXEC	= rts
SRC		= $(wildcard src/*.go)

.PHONY: clean run

run: $(EXEC)
	./$(EXEC)

$(EXEC): $(SRC)
	$(CC) -o $(EXEC) $^

clean:
	rm $(EXEC)
