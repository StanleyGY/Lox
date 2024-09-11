#include "value.hpp"
#include <stdio.h>

Value::Value() : type_(VAL_NIL) {}

Value::Value(double v) : type_(VAL_NUMBER), as_(v) {};

Value::Value(bool v) : type_(VAL_BOOL), as_(v) {};

Value::Value(std::string v) : type_(VAL_STRING), as_(v) {}

auto Value::isNumber() const -> bool {
    return type_ == VAL_NUMBER;
}

auto Value::isBool() const -> bool {
    return type_ == VAL_BOOL;
}

auto Value::isNil() const -> bool {
    return type_ == VAL_NIL;
}

auto Value::isString() const -> bool {
    return type_ == VAL_STRING;
}

auto Value::asNumber() const -> double {
    return std::get<double>(as_);
}

auto Value::asBool() const -> bool {
    return std::get<bool>(as_);
}

auto Value::asString() const -> std::string {
    return std::get<std::string>(as_);
}

auto Value::operator==(const Value &other) const -> bool {
    return type_ == other.type_ && as_ == other.as_;
}

void Value::print() const {
    switch (type_) {
        case VAL_NUMBER:
            printf("%g", asNumber());
            break;
        case VAL_BOOL:
            printf(asBool() ? "true" : "false");
            break;
        case VAL_STRING:
            printf("%s", asString().c_str());
            break;
        case VAL_NIL:
            printf("nil");
            break;
    }
}