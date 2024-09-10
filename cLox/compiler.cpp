#include "compiler.hpp"
#include "scanner.hpp"
#include <map>
#include <string>
#include <iostream>
#include <sstream>

Compiler::Compiler(const std::string &source) : source_(source) {
    current_ = 0;

    parserRules_ = std::map<TokenType, Rule>{
        {TOKEN_LEFT_PAREN, {&Compiler::grouping, nullptr, PREC_NONE}},
        {TOKEN_RIGHT_PAREN, {nullptr, nullptr, PREC_NONE}},
        {TOKEN_MINUS, {&Compiler::unary, &Compiler::binary, PREC_TERM}},
        {TOKEN_PLUS, {nullptr, &Compiler::binary, PREC_TERM}},
        {TOKEN_STAR, {nullptr, &Compiler::binary, PREC_FACTOR}},
        {TOKEN_SLASH, {nullptr, &Compiler::binary, PREC_FACTOR}},
        {TOKEN_NUMBER, {&Compiler::number, nullptr, PREC_NONE}},
        {TOKEN_EOF, {nullptr, nullptr, PREC_NONE}},
    };
}

auto Compiler::compile() -> Chunk {
    Scanner scanner{source_};
    tokens_ = scanner.scanTokens();

    // for (auto &token : tokens_) {
    //     printf("%d: %d %d %s\n", token->type_, token->start_, token->length_, source_.substr(token->start_, token->length_).c_str());
    // }
    expression();
    consume(TOKEN_EOF, "missing an EOF token");

    return chunk_;
}

void Compiler::emitByte(uint8_t byte, int lineNo) {
    // TODO: replace line number
    printf("emiting byte: %d\n", byte);
    chunk_.addCode(byte, lineNo);
}

void Compiler::emitBytes(uint8_t b1, uint8_t b2, int lineNo) {
    // TODO: replace line number
    chunk_.addCode(b1, lineNo);
    chunk_.addCode(b2, lineNo);
}

void Compiler::emitConstant(Value value, int lineNo) {
    printf("emiting constant: %lf\n", value);

    int idx = chunk_.addConstant(value);
    emitBytes(OP_CONSTANT, idx, lineNo);
}

auto Compiler::match(TokenType t) -> bool {
    if (current_ == tokens_.size()) {
        return false;
    }
    return tokens_[current_]->type_ == t;
}

auto Compiler::advance() -> const Token * {
    return tokens_[current_++].get();
}

auto Compiler::advanceIfMatch(TokenType t) -> bool {
    if (match(t)) {
        current_++;
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
    return current_ < tokens_.size();
}

auto Compiler::current() -> const Token * {
    return tokens_[current_].get();
}

auto Compiler::previous() -> const Token * {
    return tokens_[current_ - 1].get();
}

void Compiler::parsePrecedence(Precedence p) {
    std::ostringstream oss;

    auto prefixToken = advance();

    if (parserRules_.find(prefixToken->type_) == parserRules_.end()) {
        oss << "token " << prefixToken->type_ << " has no parser rule";
        throw CompilerException{oss.str()};
    }
    auto rule = parserRules_[prefixToken->type_];

    // First consider a token as a prefix operator and compiles a prefix expression.
    // Each token is a prefix operator of itself
    (this->*(rule.prefix))();

    // Then check if this prefix expresison is an operand of an infix expression.
    while (hasNext()) {
        auto infixToken = current();
        if (parserRules_.find(infixToken->type_) == parserRules_.end()) {
            oss << "token type: " << infixToken->type_ << " has no parser rule";
            throw CompilerException{oss.str()};
        }
        auto rule = parserRules_[infixToken->type_];
        if (p > rule.precedence) {
            break;
        }
        // Only advance to next token after making this infix token can be consumed
        advance();
        (this->*(rule.infix))();
    }
}

void Compiler::expression() {
    parsePrecedence(PREC_ASSIGNMENT);
}

void Compiler::binary() {
    // The left operand is compiled and binary operator is consumed
    auto token = previous();
    auto rule = parserRules_[token->type_];

    // Compile the right operand. These binary operators are all left-associative,
    // i.e. 2 + 3 + 4 === ((2 + 3) + 4)
    parsePrecedence((Precedence)((int)rule.precedence + 1));

    switch (token->type_) {
        case TOKEN_PLUS:
            emitByte(OP_ADD, token->lineNo_);
            break;
        case TOKEN_MINUS:
            emitByte(OP_SUBTRACT, token->lineNo_);
            break;
        case TOKEN_STAR:
            emitByte(OP_MULTIPLY, token->lineNo_);
            break;
        case TOKEN_SLASH:
            emitByte(OP_DIVIDE, token->lineNo_);
            break;
        default:
            break;
    }
}

void Compiler::unary() {
    auto token = previous();
    parsePrecedence(PREC_UNARY);

    switch (token->type_) {
        case TOKEN_MINUS:
            emitByte(OP_NEGATE, token->lineNo_);
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
    auto token = previous();
    Value value = std::stod(source_.substr(token->start_, token->length_));
    emitConstant(value, token->lineNo_);
}
