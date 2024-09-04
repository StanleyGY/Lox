#include "chunk.hpp"

int main(int argc, const char *argv[]) {
    Chunk chunk;
    chunk.addCode(OP_RETURN, 1);

    int idx = chunk.addConstant(3.0f);
    chunk.addCode(OP_CONSTANT, 1);
    chunk.addCode(idx, 1);

    chunk.disassemble("testing");
}