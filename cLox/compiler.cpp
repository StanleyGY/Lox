#include "compiler.hpp"

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

auto Scanner::scanTokens() -> std::vector<std::unique_ptr<Token>> {
    std::vector<std::unique_ptr<Token>> tokens;

    while (true) {
        std::unique_ptr<Token> token{nullptr};

        start_ = current_;

        if (!hasNext()) {
            tokens.emplace_back(emitToken(TOKEN_EOF));
            break;
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
                if (isdigit(peek())) {
                    token = scanNumber();
                } else if (isalpha(peek())) {
                    token = scanIdentifier();
                } else {
                    token = emitErrorToken("unknown token");
                }
                break;
        }

        if (token) {
            tokens.emplace_back(std::move(token));
        }
    }
    return tokens;
}

auto Scanner::emitToken(TokenType type) -> std::unique_ptr<Token> {
    return std::make_unique<Token>(type, start_, current_ - start_, line_);
}

auto Scanner::emitErrorToken(const std::string &message) -> std::unique_ptr<Token> {
    return std::make_unique<ErrorToken>(TOKEN_ERROR, start_, current_ - start_, line_, message);
}

auto Scanner::scanString() -> std::unique_ptr<Token> {
    while (hasNext() && !match('"')) {
        advance();
    }
    if (!hasNext()) {
        return emitErrorToken("unterminated string");
    } else {
        advance();
        return emitToken(TOKEN_STRING);
    }
}

auto Scanner::scanNumber() -> std::unique_ptr<Token> {
    while (isdigit(peek())) {
        advance();
    }
    if (match('.') && isdigit(peekNext())) {
        advance();
        while (isdigit(peek())) {
            advance();
        }
    }
    return emitToken(TOKEN_NUMBER);
}

auto Scanner::scanIdentifier() -> std::unique_ptr<Token> {
    while (isalpha(peek()) || isdigit(peek())) {
        advance();
    }

    auto word = source_.substr(start_, current_);
    if (reservedWords.find(word) != reservedWords.end()) {
        return emitToken(reservedWords[word]);
    } else {
        return emitToken(TOKEN_IDENTIFIER);
    }
}

auto Scanner::advance() -> char {
    current_++;
    return source_[current_ - 1];
}

auto Scanner::advanceIfMatch(char r) -> bool {
    if (match(r)) {
        advance();
        return true;
    }
    return false;
}

auto Scanner::peek() -> char {
    if (!hasNext()) return '\0';
    return source_[current_];
}

auto Scanner::peekNext() -> char {
    if (current_ + 1 >= source_.size()) return '\0';
    return source_[current_ + 1];
}

auto Scanner::match(char r) -> bool {
    return hasNext() && source_[current_] == r;
}

auto Scanner::hasNext() -> bool {
    return current_ + 1 < source_.length();
}