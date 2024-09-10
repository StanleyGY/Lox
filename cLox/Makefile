
CXX = g++
CXXFLAGS = -std=c++20 -Wall -pedantic-errors -g

SRCS=$(wildcard *.cpp)
HEADERS=$(wildcard *.hpp)
OBJS=$(SRCS:.cpp=.o)

MAIN = "clox"

all: $(MAIN)

${MAIN}: ${OBJS} ${HEADERS}
	${CXX} ${CXXFLAGS} ${OBJS} -o ${MAIN}

.cpp.o:
	${CXX} ${CXXFLAGS} -c $< -o $@

clean:
	${RM} ${PROGS} ${OBJS} *.o *~.