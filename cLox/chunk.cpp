#include "chunk.hpp"

void Chunk::addCode(uint8_t byte, int lineNo) {
    code_.emplace_back((OpCode)byte);
    lines_.emplace_back(lineNo);
}

auto Chunk::addConstant(Value value) -> int {
    constants_.emplace_back(value);
    return constants_.size() - 1;
}

void Chunk::disassemble(const std::string& name) const {
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

auto Chunk::disassembleInstruction(int offset) const -> int {
    uint8_t instr = code_[offset];

    switch (instr) {
        case OP_CONSTANT:
            return disassembleConstantInstruction("OP_CONSTANT", offset);
        case OP_NEGATE:
            return disassembleSimpleInstruction("OP_NEGATE", offset);
        case OP_NOT:
            return disassembleSimpleInstruction("OP_NOT", offset);
        case OP_RETURN:
            return disassembleSimpleInstruction("OP_RETURN", offset);
        case OP_ADD:
            return disassembleSimpleInstruction("OP_ADD", offset);
        case OP_SUBTRACT:
            return disassembleSimpleInstruction("OP_SUBTRACT", offset);
        case OP_MULTIPLY:
            return disassembleSimpleInstruction("OP_MULTIPLY", offset);
        case OP_DIVIDE:
            return disassembleSimpleInstruction("OP_DIVIDE", offset);
        case OP_EQUAL:
            return disassembleSimpleInstruction("OP_EQUAL", offset);
        case OP_GREATER:
            return disassembleSimpleInstruction("OP_GREATER", offset);
        case OP_LESS:
            return disassembleSimpleInstruction("OP_LESS", offset);
        default:
            printf("unknown opcode: %d", instr);
            return offset + 1;
    }
    printf("\n");
}

auto Chunk::disassembleConstantInstruction(const std::string& name, int offset) const -> int {
    int idx = code_[offset + 1];
    printf("%-16s %4d ", name.c_str(), idx);

    auto c = constants_[idx];
    c.print();
    printf("\n");
    return offset + 2;
}

auto Chunk::disassembleSimpleInstruction(const std::string& name, int offset) const -> int {
    printf("%-16s\n", name.c_str());
    return offset + 1;
}
