#ifndef CHUNK_H
#define CHUNK_H

#include "value.hpp"
#include <vector>

enum OpCode {
    OP_CONSTANT,  // OP_CONSTANT const_idx
    OP_ADD,
    OP_SUBTRACT,
    OP_MULTIPLY,
    OP_DIVIDE,
    OP_EQUAL,
    OP_GREATER,
    OP_LESS,
    OP_NEGATE,
    OP_NOT,
    OP_RETURN,
    OP_PRINT,
    OP_DEFINE_VAR,
    OP_GET_VAR,
    OP_SET_VAR,
    OP_POP,
};

// A chunk is a sequence of bytecode
class Chunk {
   public:
    Chunk() = default;
    ~Chunk() = default;

    void addCode(uint8_t byte, int lineNo);
    auto addConstant(Value value) -> int;

    void disassemble(const std::string& name) const;
    auto disassembleInstruction(int offset) const -> int;

    std::vector<OpCode> code_;
    std::vector<Value> constants_;

   private:
    // TODO: use run-length encoding of line numbers to save space
    std::vector<int> lines_;

    auto disassembleConstantInstruction(const std::string& name, int offset) const -> int;
    auto disassembleSimpleInstruction(const std::string& name, int offset) const -> int;
};

#endif