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

func (p *Parser) Parse() ([]Statement, error) {
	log.Debug().Msgf("(parser) scanning %d tokens", len(p.scan.tokens))
	stmts := []Statement{}
	for !p.done() {
		stmt, err := p.statement()
		if err != nil {
			log.Error().Msgf("(parser) error: %s", err)
			return nil, err
		}
		log.Debug().Msgf("(parser) produced statement: %s", stmt)
		stmts = append(stmts, stmt)
	}
	return stmts, nil
}

func (p *Parser) done() bool {
	_, eof := p.scan.match(TokenEOF)
	return eof
}

func (p *Parser) statement() (Statement, error) {
	p.skipComments()
	if token, ok := p.scan.match(TokenPrint); ok {
		log.Debug().Msgf("(parser) %s ... ;", *token)
		return p.printStatement()
	}
	return p.expressionStatement()
}

func (p *Parser) printStatement() (Statement, error) {
	expr, err := p.expression()
	if err != nil {
		return nil, err
	}
	if token, ok := p.scan.match(TokenSemicolon); !ok {
		return nil, ParseError{
			Err: UnexpectedTokenError{Expected: TokenSemicolon, Actual: *token},
		}
	}
	p.skipComments()
	return PrintStatement{expr: expr}, nil
}

func (p *Parser) expressionStatement() (Statement, error) {
	log.Debug().Msg("(parser) ... ;")
	expr, err := p.expression()
	if err != nil {
		return nil, err
	}
	if token, ok := p.scan.match(TokenSemicolon); !ok {
		return nil, ParseError{
			Err: UnexpectedTokenError{Expected: TokenSemicolon, Actual: *token},
		}
	}
	p.skipComments()
	return ExpressionStatement{expr: expr}, nil
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
		log.Debug().Msgf("(parser) terminal: %s", token)
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
		log.Debug().Msgf("(parser) terminal: %s", token)
		return StringExpression(token.Lexem), nil
	} else if token, ok := p.scan.match(TokenTrue); ok {
		log.Debug().Msgf("(parser) terminal: %s", token)
		return BooleanExpression(true), nil
	} else if token, ok := p.scan.match(TokenFalse); ok {
		log.Debug().Msgf("(parser) terminal: %s", token)
		return BooleanExpression(false), nil
	} else if token, ok := p.scan.match(TokenNil); ok {
		log.Debug().Msgf("(parser) terminal: %s", token)
		return NilExpression{}, nil
	}
	return nil, nil
}

func (p *Parser) skipComments() {
	for {
		if token, ok := p.scan.match(TokenComment); ok {
			log.Debug().Msgf("(parser) skip: %s", token)
			continue
		}
		break
	}
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
	token := s.peek()
	for _, t := range ts {
		if token.Type == t {
			token = s.advance()
			return &token, true
		}
	}
	return &token, false
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

type UnexpectedTokenError struct {
	Expected TokenType
	Actual   Token
}

func (e UnexpectedTokenError) Error() string {
	return fmt.Sprintf("expected %s token but got %s", e.Expected, e.Actual)
}
