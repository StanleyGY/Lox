#ifndef COMPILE_H
#define COMPILE_H

#include <memory>
#include <string>
#include <vector>

enum TokenType {
    // Single-character tokens.
    TOKEN_LEFT_PAREN,
    TOKEN_RIGHT_PAREN,
    TOKEN_LEFT_BRACE,
    TOKEN_RIGHT_BRACE,
    TOKEN_COMMA,
    TOKEN_DOT,
    TOKEN_MINUS,
    TOKEN_PLUS,
    TOKEN_SEMICOLON,
    TOKEN_SLASH,
    TOKEN_STAR,
    // One or two character tokens.
    TOKEN_BANG,
    TOKEN_BANG_EQUAL,
    TOKEN_EQUAL,
    TOKEN_EQUAL_EQUAL,
    TOKEN_GREATER,
    TOKEN_GREATER_EQUAL,
    TOKEN_LESS,
    TOKEN_LESS_EQUAL,
    // Literals.
    TOKEN_IDENTIFIER,
    TOKEN_STRING,
    TOKEN_NUMBER,
    // Keywords.
    TOKEN_AND,
    TOKEN_CLASS,
    TOKEN_ELSE,
    TOKEN_FALSE,
    TOKEN_FOR,
    TOKEN_FUN,
    TOKEN_IF,
    TOKEN_NIL,
    TOKEN_OR,
    TOKEN_PRINT,
    TOKEN_RETURN,
    TOKEN_SUPER,
    TOKEN_THIS,
    TOKEN_TRUE,
    TOKEN_VAR,
    TOKEN_WHILE,

    TOKEN_ERROR,
    TOKEN_EOF
};

class Token {
   public:
    Token(TokenType type, int start, int length, int lineNo)
        : type_(type), start_(start), length_(length), lineNo_(lineNo) {}

    TokenType type_;

    // Store pointers to the string for lexeme
    int start_;
    int length_;
    int lineNo_;
};

class ErrorToken : public Token {
   public:
    ErrorToken(TokenType type, int start, int length, int lineNo, const std::string &message)
        : Token{type, start, length, lineNo}, message_{message} {}

    const std::string &message_;
};

class Scanner {
   public:
    Scanner(const std::string &source) : source_(source) {
        start_ = 0;
        current_ = 0;
        line_ = 0;
    }
    auto scanTokens() -> std::vector<std::unique_ptr<Token>>;

   private:
    auto emitToken(TokenType type) -> std::unique_ptr<Token>;
    auto emitErrorToken(const std::string &message) -> std::unique_ptr<Token>;

    auto scanString() -> std::unique_ptr<Token>;
    auto scanNumber() -> std::unique_ptr<Token>;
    auto scanIdentifier() -> std::unique_ptr<Token>;

    auto hasNext() -> bool;
    auto peek() -> char;
    auto peekNext() -> char;
    auto match(char r) -> bool;
    auto advance() -> char;
    auto advanceIfMatch(char r) -> bool;

    const std::string &source_;
    int start_;
    int current_;
    int line_;
};

class Compiler {
   public:
    void compile(std::string source);
};

#endif