#ifndef VALUE_H
#define VALUE_H

#include <string>
#include <stdio.h>

enum ValueType {
    VAL_NUMBER,
    VAL_BOOL,
    VAL_NIL,
};

class Value {
   public:
    Value() : type_{VAL_NIL} {}
    Value(double v) : type_{VAL_NUMBER} { as_.number = v; };
    Value(bool v) : type_{VAL_BOOL} { as_.boolean = v; };

    auto isNumber() const -> bool {
        return type_ == VAL_NUMBER;
    }
    auto isBool() const -> bool {
        return type_ == VAL_BOOL;
    }
    auto isNil() const -> bool {
        return type_ == VAL_NIL;
    }
    auto asNumber() const -> double {
        return as_.number;
    }
    auto asBool() const -> bool {
        return as_.boolean;
    }
    void print() const {
        if (isNumber()) {
            printf("%g", asNumber());
        } else if (isBool()) {
            printf(asBool() ? "true" : "false");
        } else {
            printf("nil");
        }
    }

    ValueType type_;
    union {
        bool boolean;
        double number;
    } as_;
};

#endif