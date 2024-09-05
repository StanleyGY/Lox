#include "chunk.hpp"
#include "vm.hpp"

int main(int argc, const char *argv[]) {
    Chunk chunk;

    int idx = chunk.addConstant(3.0f);
    chunk.addCode(OP_CONSTANT, 1);
    chunk.addCode(idx, 1);

    idx = chunk.addConstant(6.0f);
    chunk.addCode(OP_CONSTANT, 1);
    chunk.addCode(idx, 1);

    chunk.addCode(OP_ADD, 1);

    idx = chunk.addConstant(3.0f);
    chunk.addCode(OP_CONSTANT, 1);
    chunk.addCode(idx, 1);

    chunk.addCode(OP_MULTIPLY, 1);

    chunk.addCode(OP_NEGATE, 1);
    chunk.addCode(OP_RETURN, 1);

    VM vm;
    vm.interpret(&chunk);
}