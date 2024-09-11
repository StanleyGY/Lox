#include "chunk.hpp"
#include "compiler.hpp"
#include "vm.hpp"
#include <iostream>

void compileAndRun(const std::string &test, const std::string &source) {
    Chunk chunk;
    try {
        Compiler compiler{source};
        chunk = compiler.compile();
    } catch (std::exception &e) {
        std::cerr << e.what() << std::endl;
    }
    chunk.disassemble(test);

    VM vm{&chunk};
    vm.interpret();
}

int main(int argc, const char *argv[]) {
    compileAndRun("unary/binary arithmetic", "(-1 + 2) * 3 - -4");
    compileAndRun("logical", "!(5 - 4 > 3 * 2 == !nil)");
    compileAndRun("string", "\"str\" + \"ing\" == \"string\"");
}