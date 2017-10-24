
all: game2048

game2048: clean goBuild
	docker build -t game2048 .

goBuild:
	go build game2048.go

clean:
	rm -rf game2048

.PHONY: all clean 
