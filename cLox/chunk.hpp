#ifndef CHUNK_H
#define CHUNK_H

#include <vector>

// For now, only support double precision for constant value
using Value = double;

/*
Bytecode format:
    OP_CONSTANT const_idx
    OP_NEGATE
    OP_RETURN
    OP_ADD
    OP_SUBTRACT
    OP_MULTIPLY
    OP_DIVIDE
*/
enum OpCode {
    OP_CONSTANT,
    OP_ADD,
    OP_SUBTRACT,
    OP_MULTIPLY,
    OP_DIVIDE,
    OP_NEGATE,
    OP_RETURN,
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