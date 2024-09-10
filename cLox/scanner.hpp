#ifndef SCANNER_H
#define SCANNER_H

#include <string>
#include <vector>
#include <memory>
#include "token.hpp"

class Scanner {
   public:
    Scanner(const std::string &source) : source_(source) {
        start_ = 0;
        current_ = 0;
        line_ = 0;
    }
    auto scanToken() -> std::unique_ptr<Token>;
    auto hasNext() -> bool;

   private:
    auto emitToken(TokenType type) -> std::unique_ptr<Token>;
    auto emitErrorToken(const std::string &message) -> std::unique_ptr<Token>;

    auto scanString() -> std::unique_ptr<Token>;
    auto scanNumber() -> std::unique_ptr<Token>;
    auto scanIdentifier() -> std::unique_ptr<Token>;

    auto previous() -> char;
    auto current() -> char;
    auto next() -> char;
    auto match(char r) -> bool;
    auto advance() -> char;
    auto advanceIfMatch(char r) -> bool;

    const std::string &source_;
    int start_;
    int current_;
    int line_;
};

#endif