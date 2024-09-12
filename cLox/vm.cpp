#include "vm.hpp"
#include "value.hpp"
#include <iostream>

static auto isFalsey(Value v) -> bool {
    return v.isNil() || (v.isBool() && !v.asBool()) || (v.isNumber() && v.asNumber() == 0);
}

auto VM::interpret() -> InterpretResult {
    while (ip_ < chunk_->code_.size()) {
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
                if (
                    !(peek(0).isNumber() && peek(1).isNumber()) && !(peek(0).isString() && peek(1).isString())) {
                    printRuntimeError("operand must be a number or a string");
                    return INTERPRET_RUNTIME_ERROR;
                }
                break;
            case OP_SUBTRACT:
            case OP_MULTIPLY:
            case OP_DIVIDE:
            case OP_GREATER:
            case OP_LESS:
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
                if (l.isNumber()) {
                    push(l.asNumber() + r.asNumber());
                } else {
                    push(l.asString() + r.asString());
                }
                break;
            }
            case OP_SUBTRACT: {
                auto r = pop();
                auto l = pop();
                push(l.asNumber() - r.asNumber());
                break;
            }
            case OP_MULTIPLY: {
                auto r = pop();
                auto l = pop();
                push(l.asNumber() * r.asNumber());
                break;
            }
            case OP_DIVIDE: {
                auto r = pop();
                auto l = pop();
                push(l.asNumber() / r.asNumber());
                break;
            }
            case OP_GREATER: {
                auto r = pop();
                auto l = pop();
                push(l.asNumber() > r.asNumber());
                break;
            }
            case OP_LESS: {
                auto r = pop();
                auto l = pop();
                push(l.asNumber() < r.asNumber());
                break;
            }
            case OP_EQUAL: {
                auto r = pop();
                auto l = pop();
                push(l == r);
                break;
            }
            case OP_NEGATE: {
                push(-pop().asNumber());
                break;
            }
            case OP_NOT: {
                push(isFalsey(pop()));
                break;
            }
            case OP_PRINT: {
                std::cout << pop() << std::endl;
                break;
            }
            case OP_POP:
            case OP_RETURN: {
                pop();
                break;
            }
        }

        printStack();
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
    std::cout << "          ";
    for (auto iter = stack_.begin(); iter != stack_.end(); iter++) {
        std::cout << "[ " << (*iter) << " ]";
    }
    std::cout << std::endl;
}

void VM::printRuntimeError(const std::string &message) {
    std::cerr << message << std::endl;
}