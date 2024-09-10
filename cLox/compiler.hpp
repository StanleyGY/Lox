#ifndef COMPILE_H
#define COMPILE_H

#include "chunk.hpp"
#include "token.hpp"
#include "scanner.hpp"
#include "value.hpp"
#include <map>

enum Precedence {
    PREC_NONE,
    PREC_ASSIGNMENT,  // =
    PREC_OR,          // or
    PREC_AND,         // and
    PREC_EQUALITY,    // == !=
    PREC_COMPARISON,  // < > <= >=
    PREC_TERM,        // + -
    PREC_FACTOR,      // * /
    PREC_UNARY,       // ! -
    PREC_CALL,        // . ()
    PREC_PRIMARY
};

class CompilerException : std::exception {
   public:
    CompilerException(std::string &&message) : message_{message} {}

    virtual auto what() const throw() -> const char * {
        return message_.c_str();
    }

   private:
    std::string message_;
};

class Compiler {
   public:
    Compiler(const std::string &source);
    auto compile() -> Chunk;

   private:
    auto hasNext() -> bool;
    // auto previous() -> const Token *;
    // auto current() -> const Token *;
    auto match(TokenType t) -> bool;
    void advance();
    auto advanceIfMatch(TokenType t) -> bool;
    void consume(TokenType t, std::string &&message);

    void emitByte(uint8_t byte, int lineNo);
    void emitBytes(uint8_t b1, uint8_t b2, int lineNo);
    void emitReturn();
    void emitConstant(Value value, int lineNo);

    void parsePrecedence(Precedence p);
    void expression();
    void binary();
    void unary();
    void grouping();
    void number();
    void literal();

    using ParseFunc = void (Compiler::*)();

    struct Rule {
        ParseFunc prefix;
        ParseFunc infix;
        Precedence precedence;  // precedence for an infix operator
    };

    std::map<TokenType, Rule> parserRules_;

    const std::string &source_;
    Scanner scanner_;
    std::unique_ptr<Token> prevToken_;
    std::unique_ptr<Token> currToken_;

    Chunk chunk_;
    // int current_;
    // std::vector<std::unique_ptr<Token>> tokens_;
};

#endif