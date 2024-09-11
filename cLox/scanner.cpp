#include "scanner.hpp"
#include "chunk.hpp"
#include <ctype.h>
#include <map>

auto reservedWords = std::map<std::string, TokenType>{
    {"and", TOKEN_AND},
    {"class", TOKEN_CLASS},
    {"else", TOKEN_ELSE},
    {"false", TOKEN_FALSE},
    {"fun", TOKEN_FUN},
    {"for", TOKEN_FOR},
    {"if", TOKEN_IF},
    {"nil", TOKEN_NIL},
    {"or", TOKEN_OR},
    {"print", TOKEN_PRINT},
    {"return", TOKEN_RETURN},
    {"super", TOKEN_SUPER},
    {"this", TOKEN_THIS},
    {"true", TOKEN_TRUE},
    {"var", TOKEN_VAR},
    {"while", TOKEN_WHILE},
};

auto Scanner::scanToken() -> std::unique_ptr<Token> {
    std::unique_ptr<Token> token{nullptr};

    while (true) {
        start_ = current_;
        if (!hasNext()) {
            return emitToken(TOKEN_EOF);
        }

        char c = advance();
        switch (c) {
            case '(':
                token = emitToken(TOKEN_LEFT_PAREN);
                break;
            case ')':
                token = emitToken(TOKEN_RIGHT_PAREN);
                break;
            case '{':
                token = emitToken(TOKEN_LEFT_BRACE);
                break;
            case '}':
                token = emitToken(TOKEN_RIGHT_BRACE);
                break;
            case ';':
                token = emitToken(TOKEN_SEMICOLON);
                break;
            case ',':
                token = emitToken(TOKEN_COMMA);
                break;
            case '.':
                token = emitToken(TOKEN_DOT);
                break;
            case '-':
                token = emitToken(TOKEN_MINUS);
                break;
            case '+':
                token = emitToken(TOKEN_PLUS);
                break;
            case '/':
                if (match('/')) {
                    // Skip comment
                    while (hasNext() && advance() != '\n') {
                        ;
                    }
                } else {
                    token = emitToken(TOKEN_SLASH);
                }
                break;
            case '*':
                token = emitToken(TOKEN_STAR);
                break;
            case '!':
                token = emitToken(advanceIfMatch('=') ? TOKEN_BANG_EQUAL : TOKEN_BANG);
                break;
            case '=':
                token = emitToken(advanceIfMatch('=') ? TOKEN_EQUAL_EQUAL : TOKEN_EQUAL);
                break;
            case '<':
                token = emitToken(advanceIfMatch('=') ? TOKEN_LESS_EQUAL : TOKEN_LESS);
                break;
            case '>':
                token = emitToken(advanceIfMatch('=') ? TOKEN_GREATER_EQUAL : TOKEN_GREATER);
                break;
            case ' ':
            case '\t':
            case '\r':
                break;
            case '\n':
                line_++;
                break;
            case '"':
                token = scanString();
                break;
            default:
                if (isdigit(previous())) {
                    token = scanNumber();
                } else if (isalpha(previous())) {
                    token = scanIdentifier();
                } else {
                    token = emitErrorToken("unknown token");
                }
                break;
        }

        if (token) {
            return token;
        }
    }
}

auto Scanner::emitToken(TokenType type) -> std::unique_ptr<Token> {
    return std::make_unique<Token>(type, start_, current_ - start_, line_);
}

auto Scanner::emitToken(TokenType type, int s, int e) -> std::unique_ptr<Token> {
    return std::make_unique<Token>(type, s, e - s, line_);
}

auto Scanner::emitErrorToken(const std::string &message) -> std::unique_ptr<Token> {
    return std::make_unique<ErrorToken>(TOKEN_ERROR, start_, current_ - start_, line_, message);
}

auto Scanner::scanString() -> std::unique_ptr<Token> {
    while (hasNext() && !match('"')) {
        advance();
    }
    if (hasNext()) {
        advance();
        return emitToken(TOKEN_STRING, start_ + 1, current_ - 1);
    }
    return emitErrorToken("unterminated string");
}

auto Scanner::scanNumber() -> std::unique_ptr<Token> {
    while (isdigit(current())) {
        advance();
    }
    if (match('.') && isdigit(next())) {
        advance();
        while (isdigit(current())) {
            advance();
        }
    }
    return emitToken(TOKEN_NUMBER);
}

auto Scanner::scanIdentifier() -> std::unique_ptr<Token> {
    while (isalpha(current()) || isdigit(current())) {
        advance();
    }

    auto word = source_.substr(start_, current_ - start_);
    if (reservedWords.find(word) != reservedWords.end()) {
        return emitToken(reservedWords[word]);
    } else {
        return emitToken(TOKEN_IDENTIFIER);
    }
}

auto Scanner::advance() -> char {
    return source_[current_++];
}

auto Scanner::advanceIfMatch(char r) -> bool {
    if (match(r)) {
        advance();
        return true;
    }
    return false;
}

auto Scanner::previous() -> char {
    return source_[current_ - 1];
}

auto Scanner::current() -> char {
    if (!hasNext()) return '\0';
    return source_[current_];
}

auto Scanner::next() -> char {
    if (current_ + 1 >= source_.size()) return '\0';
    return source_[current_ + 1];
}

auto Scanner::match(char r) -> bool {
    return hasNext() && source_[current_] == r;
}

auto Scanner::hasNext() -> bool {
    return current_ < source_.length();
}