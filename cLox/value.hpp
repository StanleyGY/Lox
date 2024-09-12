#ifndef VALUE_H
#define VALUE_H

#include <string>
#include <variant>

enum ValueType {
    VAL_NUMBER,
    VAL_BOOL,
    VAL_NIL,
    VAL_STRING,
};

class Value {
   public:
    Value();
    Value(double v);
    Value(bool v);
    Value(std::string v);

    auto operator==(const Value& other) const -> bool;
    auto isNumber() const -> bool;
    auto isBool() const -> bool;
    auto isNil() const -> bool;
    auto isString() const -> bool;
    auto asNumber() const -> double;
    auto asBool() const -> bool;
    auto asString() const -> std::string;

    friend auto operator<<(std::ostream& oss, const Value& v) -> std::ostream&;

    ValueType type_;
    std::variant<bool, double, std::string> as_;
};

#endif