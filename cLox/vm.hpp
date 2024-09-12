#ifndef VM_H
#define VM_H

#include <list>
#include <vector>
#include <map>
#include <string>
#include "chunk.hpp"

enum InterpretResult {
    INTERPRET_OK,
    INTERPRET_COMPILE_ERROR,
    INTERPRET_RUNTIME_ERROR
};

class VM {
   public:
    VM(const Chunk *chunk);
    auto interpret() -> InterpretResult;

   private:
    auto readByte() -> int;
    void push(Value val);
    auto pop() -> Value;
    auto peek(int dist) -> Value;

    void printStack();
    void printRuntimeError(const std::string &message);

    const Chunk *chunk_;
    // Points to the index of chunk_.code_
    int ip_;
    // Stores the values after executing chunk_. Remember it stores code in a post-order
    // traversal of the AST tree
    std::list<Value> stack_;

    std::map<std::string, Value> globals_;
};

#endif