package lox

import (
	"fmt"
	"strconv"

	"github.com/rs/zerolog/log"
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
	log.Debug().Msgf("(parser) scanning %d tokens", len(p.scan.tokens))
	expr, err := p.expression()
	if err != nil {
		log.Error().Msgf("(parser) error: %s", err)
		return nil, err
	}
	log.Debug().Msgf("(parser) produced expression %s", expr)
	return expr, err
}

func (p *Parser) expression() (Expression, error) {
	return p.equality()
}

func (p *Parser) equality() (Expression, error) {
	expr, err := p.comparison()
	if err != nil {
		return nil, err
	}
	for {
		token, ok := p.scan.match(TokenBangEqual, TokenEqualEqual)
		if !ok {
			break
		}
		op, err := token.Operator()
		if err != nil {
			return nil, ParseError{Err: err}
		}
		log.Debug().Msgf("(parser) %s %s ...", expr, op)
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
	for {
		token, ok := p.scan.match(TokenGreater, TokenGreaterEqual, TokenLess, TokenLessEqual)
		if !ok {
			break
		}
		op, err := token.Operator()
		if err != nil {
			return nil, ParseError{Err: err}
		}
		log.Debug().Msgf("(parser) %s %s ...", expr, op)
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
	for {
		token, ok := p.scan.match(TokenPlus, TokenMinus)
		if !ok {
			break
		}
		op, err := token.Operator()
		if err != nil {
			return nil, ParseError{Err: err}
		}
		log.Debug().Msgf("(parser) %s %s ...", expr, op)
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
	for {
		token, ok := p.scan.match(TokenSlash, TokenStar)
		if !ok {
			break
		}
		op, err := token.Operator()
		if err != nil {
			return nil, ParseError{Err: err}
		}
		log.Debug().Msgf("(parser) %s %s ...", expr, op)
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		expr = BinaryExpression{op: *op, left: expr, right: right}
	}
	return expr, nil
}

func (p *Parser) unary() (Expression, error) {
	if token, ok := p.scan.match(TokenBang, TokenMinus); ok {
		op, err := token.Operator()
		if err != nil {
			return nil, ParseError{Err: err}
		}
		log.Debug().Msgf("(parser) %s ...", op)
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
	if token, ok := p.scan.match(TokenNumber); ok {
		log.Debug().Msgf("(parser) number(%s)", token.Lexem)
		n, err := strconv.ParseFloat(token.Lexem, 64)
		if err != nil {
			return nil, ParseError{
				Err: NumberConversionError{
					Token: *token,
					Err:   err,
				},
			}
		}
		return NumericExpression(n), nil
	} else if token, ok := p.scan.match(TokenString); ok {
		log.Debug().Msgf("(parser) \"%s\"", token.Lexem)
		return StringExpression(token.Lexem), nil
	} else if _, ok := p.scan.match(TokenTrue); ok {
		log.Debug().Msg("(parser) true")
		return BooleanExpression(true), nil
	} else if _, ok := p.scan.match(TokenFalse); ok {
		log.Debug().Msg("(parser) false")
		return BooleanExpression(false), nil
	} else if _, ok := p.scan.match(TokenNil); ok {
		log.Debug().Msg("(parser) nil")
		return NilExpression{}, nil
	}
	return nil, nil
}

func (p *Parser) grouping() (Expression, error) {
	if token, ok := p.scan.match(TokenLeftParen); ok {
		log.Debug().Msg("(parser) ( ... )")
		expr, err := p.expression()
		if err != nil {
			return nil, err
		}
		if _, ok := p.scan.match(TokenRightParen); !ok {
			return nil, ParseError{
				Err: UnmatchedTokenError{*token},
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

func (s *tokenScanner) advance() Token {
	if s.done() {
		return EofToken
	}
	s.offset += 1
	return s.tokens[s.offset-1]
}

func (s *tokenScanner) match(ts ...TokenType) (*Token, bool) {
	for _, t := range ts {
		if s.check(t) {
			token := s.advance()
			return &token, true
		}
	}
	return nil, false
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
