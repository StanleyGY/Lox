#include "chunk.hpp"

#include <iostream>

void Chunk::addCode(uint8_t byte, int lineNo) {
    code_.emplace_back((OpCode)byte);
    lines_.emplace_back(lineNo);
}

auto Chunk::addConstant(Value value) -> int {
    values_.emplace_back(value);
    return values_.size() - 1;
}

void Chunk::disassemble(const std::string& name) {
    printf("== %s ==\n", name.c_str());
    for (int offset = 0; offset < code_.size();) {
        printf("%04d ", offset);
        if (offset > 0 && lines_[offset] == lines_[offset - 1]) {
            printf("   | ");
        } else {
            printf("%4d ", lines_[offset]);
        }
        offset = disassembleInstruction(offset);
    }
}

auto Chunk::disassembleInstruction(int offset) -> int {
    uint8_t instr = code_[offset];

    switch (instr) {
        case OP_CONSTANT:
            return disassembleConstantInstruction("OP_CONSTANT", offset);
        case OP_RETURN:
            return disassembleSimpleInstruction("OP_RETURN", offset);
        default:
            std::cout << "unknown opcode: " << instr << std::endl;
            return offset + 1;
    }
}

auto Chunk::disassembleConstantInstruction(const std::string& name, int offset) -> int {
    int idx = code_[offset + 1];
    printf("%-16s %4d %g\n", name.c_str(), idx, values_[idx]);
    return offset + 2;
}

auto Chunk::disassembleSimpleInstruction(const std::string& name, int offset) -> int {
    printf("%-16s\n", name.c_str());
    return offset + 1;
}
