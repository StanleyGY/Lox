#include "compiler.hpp"
#include <map>
#include <string>
#include <iostream>
#include <sstream>

Compiler::Compiler(const std::string &source) : source_(source), scanner_(Scanner{source}) {
    parserRules_ = std::map<TokenType, Rule>{
        {TOKEN_LEFT_PAREN, {&Compiler::grouping, nullptr, PREC_NONE}},
        {TOKEN_RIGHT_PAREN, {nullptr, nullptr, PREC_NONE}},
        {TOKEN_MINUS, {&Compiler::unary, &Compiler::binary, PREC_TERM}},
        {TOKEN_PLUS, {nullptr, &Compiler::binary, PREC_TERM}},
        {TOKEN_STAR, {nullptr, &Compiler::binary, PREC_FACTOR}},
        {TOKEN_SLASH, {nullptr, &Compiler::binary, PREC_FACTOR}},
        {TOKEN_BANG_EQUAL, {nullptr, &Compiler::binary, PREC_EQUALITY}},
        {TOKEN_EQUAL_EQUAL, {nullptr, &Compiler::binary, PREC_EQUALITY}},
        {TOKEN_LESS, {nullptr, &Compiler::binary, PREC_COMPARISON}},
        {TOKEN_LESS_EQUAL, {nullptr, &Compiler::binary, PREC_COMPARISON}},
        {TOKEN_GREATER, {nullptr, &Compiler::binary, PREC_COMPARISON}},
        {TOKEN_GREATER_EQUAL, {nullptr, &Compiler::binary, PREC_COMPARISON}},
        {TOKEN_NUMBER, {&Compiler::number, nullptr, PREC_NONE}},
        {TOKEN_STRING, {&Compiler::string, nullptr, PREC_NONE}},
        {TOKEN_TRUE, {&Compiler::literal, nullptr, PREC_NONE}},
        {TOKEN_FALSE, {&Compiler::literal, nullptr, PREC_NONE}},
        {TOKEN_NIL, {&Compiler::literal, nullptr, PREC_NONE}},
        {TOKEN_BANG, {&Compiler::unary, nullptr, PREC_UNARY}},
        {TOKEN_EOF, {nullptr, nullptr, PREC_NONE}},
    };
}

auto Compiler::compile() -> Chunk {
    // This causes the first token to be stored in `currToken_`
    advance();

    while (currToken_->type_ != TOKEN_EOF) {
        declaration();
    }

    consume(TOKEN_EOF, "missing an EOF token");
    return chunk_;
}

void Compiler::emitByte(uint8_t byte, int lineNo) {
    chunk_.addCode(byte, lineNo);
}

void Compiler::emitBytes(uint8_t b1, uint8_t b2, int lineNo) {
    chunk_.addCode(b1, lineNo);
    chunk_.addCode(b2, lineNo);
}

void Compiler::emitConstant(Value value, int lineNo) {
    int idx = chunk_.addConstant(value);
    emitBytes(OP_CONSTANT, idx, lineNo);
}

auto Compiler::match(TokenType t) -> bool {
    return currToken_->type_ == t;
}

void Compiler::advance() {
    prevToken_ = std::move(currToken_);
    currToken_ = scanner_.scanToken();
}

auto Compiler::advanceIfMatch(TokenType t) -> bool {
    if (match(t)) {
        advance();
        return true;
    }
    return false;
}

void Compiler::consume(TokenType t, std::string &&message) {
    if (!advanceIfMatch(t)) {
        throw CompilerException{std::move(message)};
    }
}

auto Compiler::hasNext() -> bool {
    return scanner_.hasNext();
}

void Compiler::parsePrecedence(Precedence p) {
    std::ostringstream oss;

    advance();
    if (parserRules_.find(prevToken_->type_) == parserRules_.end()) {
        oss << "token " << prevToken_->type_ << " has no parser rule";
        throw CompilerException{oss.str()};
    }
    auto rule = parserRules_[prevToken_->type_];

    // First consider a token as a prefix operator and compiles a prefix expression.
    // Each token is a prefix operator of itself
    (this->*(rule.prefix))();

    // Then check if this prefix expresison is an operand of an infix expression.
    while (hasNext()) {
        if (parserRules_.find(currToken_->type_) == parserRules_.end()) {
            oss << "token type: " << currToken_->type_ << " has no parser rule";
            throw CompilerException{oss.str()};
        }
        auto rule = parserRules_[currToken_->type_];
        if (p > rule.precedence) {
            break;
        }
        // Only advance to next token after ensuring this infix token can be consumed
        advance();
        (this->*(rule.infix))();
    }
}

void Compiler::declaration() {
    // if (advanceIfMatch(TOKEN_VAR)) {
    //     varDecl();
    // } else {
    //     statement();
    // }
    statement();
}

void Compiler::varDecl() {
    // if (advanceIfMatch(TOKEN_EQUAL)) {
    //     // initializer
    // }
    // consume(TOKEN_SEMICOLON, "variable declaration missing a ';'");
}

void Compiler::statement() {
    if (advanceIfMatch(TOKEN_PRINT)) {
        printStmt();
    } else {
        expressionStmt();
    }
}

void Compiler::printStmt() {
    auto lineNo = prevToken_->lineNo_;
    expression();
    consume(TOKEN_SEMICOLON, "statement missing a ';'");
    emitByte(OP_PRINT, lineNo);
}

void Compiler::expressionStmt() {
    auto lineNo = currToken_->lineNo_;
    expression();
    consume(TOKEN_SEMICOLON, "statement missing a ';'");
    emitByte(OP_POP, lineNo);
}

void Compiler::expression() {
    parsePrecedence(PREC_ASSIGNMENT);
}

void Compiler::binary() {
    // The left operand is compiled and binary operator is consumed
    auto opType = prevToken_->type_;
    auto opLineNo = prevToken_->lineNo_;
    auto rule = parserRules_[opType];

    // Compile the right operand. These binary operators are all left-associative,
    // i.e. 2 + 3 + 4 === ((2 + 3) + 4)
    parsePrecedence((Precedence)((int)rule.precedence + 1));

    switch (opType) {
        case TOKEN_PLUS:
            emitByte(OP_ADD, opLineNo);
            break;
        case TOKEN_MINUS:
            emitByte(OP_SUBTRACT, opLineNo);
            break;
        case TOKEN_STAR:
            emitByte(OP_MULTIPLY, opLineNo);
            break;
        case TOKEN_SLASH:
            emitByte(OP_DIVIDE, opLineNo);
            break;
        case TOKEN_BANG_EQUAL:
            // a != b is equiv to !(a == b)
            emitBytes(OP_EQUAL, OP_NOT, opLineNo);
            break;
        case TOKEN_EQUAL_EQUAL:
            emitByte(OP_EQUAL, opLineNo);
            break;
        case TOKEN_LESS_EQUAL:
            // a <= b is equiv to !(a > b)
            emitBytes(OP_GREATER, OP_NOT, opLineNo);
            break;
        case TOKEN_LESS:
            emitByte(OP_LESS, opLineNo);
            break;
        case TOKEN_GREATER_EQUAL:
            // a >= b is equiv to !(a < b)
            emitBytes(OP_LESS, OP_NOT, opLineNo);
            break;
        case TOKEN_GREATER:
            emitByte(OP_GREATER, opLineNo);
            break;
        default:
            break;
    }
}

void Compiler::unary() {
    auto opType = prevToken_->type_;
    auto opLineNo = prevToken_->lineNo_;
    parsePrecedence(PREC_UNARY);

    switch (opType) {
        case TOKEN_MINUS:
            emitByte(OP_NEGATE, opLineNo);
            break;
        case TOKEN_BANG:
            emitByte(OP_NOT, opLineNo);
            break;
        default:
            break;
    }
}

void Compiler::grouping() {
    expression();
    consume(TOKEN_RIGHT_PAREN, "grouping expr missing ')'");
}

void Compiler::number() {
    double value = std::stod(source_.substr(prevToken_->start_, prevToken_->length_));
    // Store the number constant in a separate constant_ array because
    // number cosntants can have billions of variants
    emitConstant(value, prevToken_->lineNo_);
}

void Compiler::string() {
    auto value = source_.substr(prevToken_->start_, prevToken_->length_);
    emitConstant(value, prevToken_->lineNo_);
}

void Compiler::literal() {
    // Technically, we can save execution time and space by not storing
    // these literals in a separate constant_ array. So a more optimal solution
    // is to emit a bytecode instruction.
    switch (prevToken_->type_) {
        case TOKEN_TRUE:
            emitConstant(true, prevToken_->lineNo_);
            break;
        case TOKEN_FALSE:
            emitConstant(false, prevToken_->lineNo_);
            break;
        case TOKEN_NIL:
            emitConstant(Value{}, prevToken_->lineNo_);
            break;
        default:
            throw CompilerException{"processing literal for invalid token"};
    }
}