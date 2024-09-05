#include "vm.hpp"

auto VM::interpret(const Chunk* chunk) -> InterpretResult {
#define CONSTANT() chunk_->constants_[readByte()];
#define BINARY_OP(op)     \
    {                     \
        double r = pop(); \
        double l = pop(); \
        push(l op r);     \
    }

    // Reset all fields
    chunk_ = chunk;
    ip_ = 0;
    stack_ = std::list<Value>{};

    while (true) {
        printStack();
        chunk_->disassembleInstruction(ip_);

        uint8_t instruction;
        switch (instruction = readByte()) {
            case OP_CONSTANT: {
                auto constant = CONSTANT();
                push(constant);
                break;
            }
            case OP_ADD: {
                BINARY_OP(+);
                break;
            }
            case OP_SUBTRACT: {
                BINARY_OP(-);
                break;
            }
            case OP_MULTIPLY: {
                BINARY_OP(*);
                break;
            }
            case OP_DIVIDE: {
                BINARY_OP(/);
                break;
            }
            case OP_NEGATE: {
                push(-pop());
                break;
            }
            case OP_RETURN: {
                pop();
                return INTERPRET_OK;
            }
        }
    }
}

auto VM::readByte() -> int {
    return chunk_->code_[ip_++];
}

void VM::push(Value val) {
    stack_.push_back(val);
}

auto VM::pop() -> Value {
    auto val = stack_.back();
    stack_.pop_back();
    return val;
}

void VM::printStack() {
    printf("          ");
    for (auto iter = stack_.begin(); iter != stack_.end(); iter++) {
        printf("[ ");
        printf("%g", (*iter));
        printf(" ]");
    }
    printf("\n");
}