#ifndef VM_H
#define VM_H

#include <list>
#include <vector>

#include "chunk.hpp"

enum InterpretResult {
    INTERPRET_OK,
    INTERPRET_COMPILE_ERROR,
    INTERPRET_RUNTIME_ERROR
};

class VM {
   public:
    VM() = default;

    auto interpret(const Chunk *chunk) -> InterpretResult;

   private:
    auto readByte() -> int;
    void push(Value val);
    auto pop() -> Value;
    void printStack();

    const Chunk *chunk_;
    int ip_;
    std::list<Value> stack_;  // using list for debugging
};

#endif