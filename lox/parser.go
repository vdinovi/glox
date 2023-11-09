package lox

import "strconv"

type Parser struct {
	scan tokenScanner
}

func NewParser(tokens []Token) Parser {
	return Parser{
		scan: tokenScanner{tokens: tokens},
	}
}

func (p *Parser) Parse() (Expression, error) {
	return p.expression(), nil
}

func (p *Parser) expression() Expression {
	return p.equality()
}

func (p *Parser) equality() Expression {
	expr := p.comparison()
	for p.scan.match(TokenBang, TokenBangEqual) {
		op := p.scan.previous().Operator()
		right := p.comparison()
		expr = BinaryExpression{op: op, left: expr, right: right}
	}
	return expr
}

func (p *Parser) comparison() Expression {
	expr := p.term()
	for p.scan.match(TokenGreater, TokenGreaterEqual, TokenLess, TokenLessEqual) {
		op := p.scan.previous().Operator()
		right := p.term()
		expr = BinaryExpression{op: op, left: expr, right: right}
	}
	return expr
}

func (p *Parser) term() Expression {
	expr := p.factor()
	for p.scan.match(TokenPlus, TokenMinus) {
		op := p.scan.previous().Operator()
		right := p.factor()
		expr = BinaryExpression{op: op, left: expr, right: right}
	}
	return expr
}

func (p *Parser) factor() Expression {
	expr := p.unary()
	for p.scan.match(TokenSlash, TokenStar) {
		op := p.scan.previous().Operator()
		right := p.unary()
		expr = BinaryExpression{op: op, left: expr, right: right}
	}
	return expr
}

func (p *Parser) unary() Expression {
	if p.scan.match(TokenBang, TokenMinus) {
		op := p.scan.previous().Operator()
		right := p.unary()
		return UnaryExpression{op: op, right: right}
	}
	return p.primary()
}

func (p *Parser) primary() Expression {
	if p.scan.match(TokenNumber) {
		n, err := strconv.ParseFloat(p.scan.previous().Lexem, 64)
		if err != nil {
			panic(err)
		}
		return NumberExpression(n)
	} else if p.scan.match(TokenString) {
		return StringExpression(p.scan.previous().Lexem)
	} else if p.scan.match(TokenTrue) {
		return BoolExpression(true)
	} else if p.scan.match(TokenFalse) {
		return BoolExpression(false)
	} else if p.scan.match(TokenNil) {
		return NilExpression{}
	} else if p.scan.match(TokenLeftParen) {
		expr := p.expression()
		if !p.scan.match(TokenRightParen) {
			panic("unmatched '('")
		}
		return GroupingExpression{expr: expr}
	}
	panic("reached end of Parse without terminal")
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
