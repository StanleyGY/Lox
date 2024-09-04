
build: chunk.o main.o
	g++ -std=c++20 -Wall chunk.o main.o -o clox

clean:
	rm clox *.o

chunk.o:
	g++ -std=c++20 -Wall -c chunk.cpp
main.o:
	g++ -std=c++20 -Wall -c main.cpp