
build: chunk.o main.o vm.o compiler.o
	g++ -std=c++20 -Wall chunk.o main.o vm.o compiler.o -o clox

clean:
	rm clox *.o

chunk.o: chunk.cpp chunk.hpp
	g++ -std=c++20 -Wall -c chunk.cpp
vm.o: vm.cpp vm.hpp
	g++ -std=c++20 -Wall -c vm.cpp
compiler.o: compiler.cpp compiler.hpp
	g++ -std=c++20 -Wall -c compiler.cpp
main.o: main.cpp
	g++ -std=c++20 -Wall -c main.cpp