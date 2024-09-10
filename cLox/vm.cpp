#include "vm.hpp"
#include "value.hpp"
#include <iostream>

auto VM::interpret() -> InterpretResult {
    while (true) {
        printStack();
        chunk_->disassembleInstruction(ip_);

        uint8_t instruction = readByte();

        // Check types
        switch (instruction) {
            case OP_NEGATE: {
                if (!peek(0).isNumber()) {
                    printRuntimeError("operand must be a number");
                    return INTERPRET_RUNTIME_ERROR;
                }
                break;
            }
            case OP_ADD:
            case OP_SUBTRACT:
            case OP_MULTIPLY:
            case OP_DIVIDE:
                if (!peek(0).isNumber() || !peek(1).isNumber()) {
                    printRuntimeError("operand must be a number");
                    return INTERPRET_RUNTIME_ERROR;
                }
                break;
            default:
                break;
        }

        // Execute bytecode
        switch (instruction) {
            case OP_CONSTANT: {
                auto constant = chunk_->constants_[readByte()];
                push(constant);
                break;
            }
            case OP_ADD: {
                auto r = pop();
                auto l = pop();
                push(Value{l.asNumber() + r.asNumber()});
                break;
            }
            case OP_SUBTRACT: {
                auto r = pop();
                auto l = pop();
                push(Value{l.asNumber() - r.asNumber()});
                break;
            }
            case OP_MULTIPLY: {
                auto r = pop();
                auto l = pop();
                push(Value{l.asNumber() * r.asNumber()});
                break;
            }
            case OP_DIVIDE: {
                auto r = pop();
                auto l = pop();
                push(Value{l.asNumber() / r.asNumber()});
                break;
            }
            case OP_NEGATE: {
                push(Value{-pop().asNumber()});
                break;
            }
            case OP_RETURN: {
                pop();
                break;
            }
        }
    }
    return INTERPRET_OK;
}

auto VM::readByte() -> int {
    return chunk_->code_[ip_++];
}

auto VM::peek(int dist) -> Value {
    auto iter = stack_.crbegin();
    std::advance(iter, dist);
    return *iter;
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
        if ((*iter).isNumber()) {
            printf("%g", (*iter).asNumber());
        }
        printf(" ]");
    }
    printf("\n");
}

void VM::printRuntimeError(const std::string &message) {
    std::cerr << message << std::endl;
}