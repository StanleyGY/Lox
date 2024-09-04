#ifndef CHUNK_H
#define CHUNK_H

#include <vector>

// For now, only support double precision for constant value
typedef double Value;

/*
Bytecode format:
    OP_CONSTANT const_idx
    OP_RETURN
*/
enum OpCode {
    OP_CONSTANT,
    OP_RETURN,
};

// A chunk is a sequence of bytecode
class Chunk {
   public:
    Chunk() = default;
    ~Chunk() = default;

    void addCode(uint8_t byte, int lineNo);
    auto addConstant(Value value) -> int;
    void disassemble(const std::string& name);

   private:
    auto disassembleInstruction(int offset) -> int;
    auto disassembleConstantInstruction(const std::string& name, int offset) -> int;
    auto disassembleSimpleInstruction(const std::string& name, int offset) -> int;

    std::vector<OpCode> code_;
    std::vector<int> lines_;  // TODO: use run-length encoding of line numbers to save space
    std::vector<Value> values_;
};

#endif