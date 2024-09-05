
build: chunk.o main.o vm.o
	g++ -std=c++20 -Wall chunk.o main.o vm.o -o clox

clean:
	rm clox *.o

chunk.o: chunk.cpp
	g++ -std=c++20 -Wall -c chunk.cpp
vm.o: vm.cpp
	g++ -std=c++20 -Wall -c vm.cpp
main.o: main.cpp
	g++ -std=c++20 -Wall -c main.cpp