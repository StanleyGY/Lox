#include "chunk.hpp"
#include "compiler.hpp"
#include "vm.hpp"
#include <iostream>

int main(int argc, const char *argv[]) {
    // Chunk chunk;

    // int idx = chunk.addConstant(3.0f);
    // chunk.addCode(OP_CONSTANT, 1);
    // chunk.addCode(idx, 1);

    // idx = chunk.addConstant(6.0f);
    // chunk.addCode(OP_CONSTANT, 1);
    // chunk.addCode(idx, 1);

    // chunk.addCode(OP_ADD, 1);

    // idx = chunk.addConstant(3.0f);
    // chunk.addCode(OP_CONSTANT, 1);
    // chunk.addCode(idx, 1);

    // chunk.addCode(OP_MULTIPLY, 1);

    // chunk.addCode(OP_NEGATE, 1);
    // chunk.addCode(OP_RETURN, 1);

    // VM vm;
    // vm.interpret(&chunk);

    // std::string source = "var a = a + b; a = 123.456; b = \"type\"; a >= b;";
    // Scanner scanner{source};
    // auto tokens = scanner.scanTokens();
    // for (auto &token : tokens) {
    //     printf("%d %d %s\n", token->start_, token->length_, source.substr(token->start_, token->length_).c_str());
    // }

    std::string source = "3 * 4 + 5";
    Compiler compiler{source};

    try {
        auto chunk = compiler.compile();
    } catch (const std::exception &e) {
        std::cout << e.what() << std::endl;
    }
    // chunk.disassemble();
}