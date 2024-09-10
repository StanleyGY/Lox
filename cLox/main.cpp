#include "chunk.hpp"
#include "compiler.hpp"
#include "vm.hpp"
#include <iostream>

int main(int argc, const char *argv[]) {
    std::string source = "(-1 + 2) * 3 - -4 + true";
    Compiler compiler{source};
    Chunk chunk;
    try {
        chunk = compiler.compile();
    } catch (const std::exception &e) {
        std::cout << e.what() << std::endl;
    }
    chunk.disassemble("testing");
}