#include "chunk.hpp"
#include "compiler.hpp"
#include "vm.hpp"
#include <iostream>

int main(int argc, const char *argv[]) {
    std::string source = "3 * (4 + 5)";
    Compiler compiler{source};
    Chunk chunk;
    try {
        chunk = compiler.compile();
    } catch (const std::exception &e) {
        std::cout << e.what() << std::endl;
    }
    chunk.disassemble("testing");
}