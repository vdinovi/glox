package lox

import (
	"fmt"
	"strconv"
)

type Parser struct {
	scan tokenScanner
}

func NewParser(tokens []Token) Parser {
	return Parser{
		scan: tokenScanner{tokens: tokens},
	}
}

func (p *Parser) Parse() (Expression, error) {
	return p.expression()
}

func (p *Parser) expression() (Expression, error) {
	return p.equality()
}

func (p *Parser) equality() (Expression, error) {
	expr, err := p.comparison()
	if err != nil {
		return nil, err

	}
	for p.scan.match(TokenBang, TokenBangEqual) {
		op, err := p.scan.previous().Operator()
		if err != nil {
			return nil, ParseError{Err: err}
		}
		right, err := p.comparison()
		if err != nil {
			return nil, err
		}
		expr = BinaryExpression{op: *op, left: expr, right: right}
	}
	return expr, nil
}

func (p *Parser) comparison() (Expression, error) {
	expr, err := p.term()
	if err != nil {
		return nil, err
	}
	for p.scan.match(TokenGreater, TokenGreaterEqual, TokenLess, TokenLessEqual) {
		op, err := p.scan.previous().Operator()
		if err != nil {
			return nil, ParseError{Err: err}
		}
		right, err := p.term()
		if err != nil {
			return nil, err
		}
		expr = BinaryExpression{op: *op, left: expr, right: right}
	}
	return expr, nil
}

func (p *Parser) term() (Expression, error) {
	expr, err := p.factor()
	if err != nil {
		return nil, err
	}
	for p.scan.match(TokenPlus, TokenMinus) {
		op, err := p.scan.previous().Operator()
		if err != nil {
			return nil, ParseError{Err: err}
		}
		right, err := p.factor()
		if err != nil {
			return nil, err
		}
		expr = BinaryExpression{op: *op, left: expr, right: right}
	}
	return expr, nil
}

func (p *Parser) factor() (Expression, error) {
	expr, err := p.unary()
	if err != nil {
		return nil, err
	}
	for p.scan.match(TokenSlash, TokenStar) {
		op, err := p.scan.previous().Operator()
		if err != nil {
			return nil, ParseError{Err: err}
		}
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		expr = BinaryExpression{op: *op, left: expr, right: right}
	}
	return expr, nil
}

func (p *Parser) unary() (Expression, error) {
	if p.scan.match(TokenBang, TokenMinus) {
		op, err := p.scan.previous().Operator()
		if err != nil {
			return nil, ParseError{Err: err}
		}
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		return UnaryExpression{op: *op, right: right}, nil
	}
	return p.primary()
}

func (p *Parser) primary() (Expression, error) {
	if expr, err := p.literal(); err != nil {
		return nil, err
	} else if expr != nil {
		return expr, nil
	}

	if expr, err := p.grouping(); err != nil {
		return nil, err
	} else if expr != nil {
		return expr, nil
	}

	return nil, &ParseError{
		Err: MissingTerminalError{},
	}
}

func (p *Parser) literal() (Expression, error) {
	if p.scan.match(TokenNumber) {
		numberToken := p.scan.previous()
		n, err := strconv.ParseFloat(p.scan.previous().Lexem, 64)
		if err != nil {
			return nil, ParseError{
				Err: NumberConversionError{
					Token: numberToken,
					Err:   err,
				},
			}
		}
		return NumberExpression(n), nil
	} else if p.scan.match(TokenString) {
		return StringExpression(p.scan.previous().Lexem), nil
	} else if p.scan.match(TokenTrue) {
		return BooleanExpression(true), nil
	} else if p.scan.match(TokenFalse) {
		return BooleanExpression(false), nil
	} else if p.scan.match(TokenNil) {
		return NilExpression{}, nil
	}
	return nil, nil
}

func (p *Parser) grouping() (Expression, error) {
	if p.scan.match(TokenLeftParen) {
		leftParenToken := p.scan.previous()
		expr, err := p.expression()
		if err != nil {
			return nil, err
		}
		if !p.scan.match(TokenRightParen) {
			return nil, ParseError{
				Err: UnmatchedTokenError{leftParenToken},
			}
		}
		return GroupingExpression{expr: expr}, nil
	}
	return nil, nil
}

type tokenScanner struct {
	tokens []Token
	offset int
}

func (s *tokenScanner) done() bool {
	return s.offset == len(s.tokens)
}

func (s *tokenScanner) peek() Token {
	if s.done() {
		return EofToken
	}
	return s.tokens[s.offset]
}

func (s *tokenScanner) previous() Token {
	if s.offset == 0 {
		return NoneToken
	}
	return s.tokens[s.offset-1]
}

func (s *tokenScanner) advance() Token {
	if s.done() {
		return EofToken
	}
	s.offset += 1
	return s.previous()
}

func (s *tokenScanner) match(ts ...TokenType) bool {
	for _, t := range ts {
		if s.check(t) {
			s.advance()
			return true
		}
	}
	return false
}

func (s *tokenScanner) check(t TokenType) bool {
	if s.done() {
		return false
	}
	token := s.peek()
	return token.Type == t
}

type ParseError struct {
	Err error
}

func (e ParseError) Error() string {
	return fmt.Sprintf("ParseError: %s", e.Err)
}

func (e ParseError) Unwrap() error {
	return e.Err
}

type UnmatchedTokenError struct {
	Token
}

func (e UnmatchedTokenError) Error() string {
	return fmt.Sprintf("unmatched %s at (%d, %d)", e.Token, e.Token.Line, e.Token.Column)
}

type MissingTerminalError struct{}

func (e MissingTerminalError) Error() string {
	return "missing terminal"
}

type NumberConversionError struct {
	Token
	Err error
}

func (e NumberConversionError) Error() string {
	return fmt.Sprintf("cannot convert number %q at (%d, %d): %s", e.Token.Lexem, e.Token.Line, e.Token.Column, e.Err)
}

func (e NumberConversionError) Unwrap() error {
	return e.Err
}
