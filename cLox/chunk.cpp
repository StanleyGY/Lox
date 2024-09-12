#include "chunk.hpp"
#include <iostream>
#include <map>
#include <format>

void Chunk::addCode(uint8_t byte, int lineNo) {
    code_.emplace_back((OpCode)byte);
    lines_.emplace_back(lineNo);
}

auto Chunk::addConstant(Value value) -> int {
    // TODO: can de-duplicate the constant
    constants_.emplace_back(value);
    return constants_.size() - 1;
}

void Chunk::disassemble(const std::string& name) const {
    std::cout << "== " << name << " ==" << std::endl;

    for (int offset = 0; offset < code_.size();) {
        std::cout << std::format("{0:04}", offset);
        if (offset > 0 && lines_[offset] == lines_[offset - 1]) {
            std::cout << "   | ";
        } else {
            std::cout << std::format("{0:4} ", lines_[offset]);
        }
        offset = disassembleInstruction(offset);
    }
}

auto Chunk::disassembleInstruction(int offset) const -> int {
    uint8_t instr = code_[offset];

    switch (instr) {
        case OP_CONSTANT:
            return disassembleConstantInstruction("OP_CONSTANT", offset);
        case OP_DEFINE_VAR:
            return disassembleConstantInstruction("OP_DEFINE_VAR", offset);
        case OP_GET_VAR:
            return disassembleConstantInstruction("OP_GET_VAR", offset);
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
        case OP_PRINT:
            return disassembleSimpleInstruction("OP_PRINT", offset);
        case OP_POP:
            return disassembleSimpleInstruction("OP_POP", offset);
        default:
            std::cout << "unknown opcode: " << instr;
            return offset + 1;
    }
    std::cout << std::endl;
}

auto Chunk::disassembleConstantInstruction(const std::string& name, int offset) const -> int {
    int idx = code_[offset + 1];
    auto c = constants_[idx];
    std::cout << std::format("{0:<16} {1:4} ", name, idx) << c << std::endl;
    return offset + 2;
}

auto Chunk::disassembleSimpleInstruction(const std::string& name, int offset) const -> int {
    std::cout << std::format("{0:<16}", name) << std::endl;
    return offset + 1;
}
