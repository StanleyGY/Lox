#include "chunk.hpp"
#include "compiler.hpp"
#include "vm.hpp"
#include <iostream>

void compileAndRun(const std::string &test, const std::string &source) {
    Compiler compiler{source};
    Chunk chunk;
    try {
        chunk = compiler.compile();
    } catch (const std::exception &e) {
        std::cout << e.what() << std::endl;
    }
    chunk.disassemble(test);

    VM vm{&chunk};
    vm.interpret();
}

int main(int argc, const char *argv[]) {
    compileAndRun("unary/binary arithmetic", "(-1 + 2) * 3 - -4");
    compileAndRun("logical", "!true");
}