package lox

import (
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

func (p *Parser) Parse() (Program, error) {
	log.Debug().Msgf("(parser) scanning %d tokens", len(p.scan.tokens))
	stmts := []Statement{}
	for !p.done() {
		p.skipComments()
		stmt, err := p.declaration()
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

func (p *Parser) declaration() (Statement, error) {
	if token, ok := p.scan.match(TokenVar); ok {
		log.Debug().Msgf("(parser) %s ... ;", token.Lexem)
		return p.varDeclaration()
	}
	return p.statement()
}

func (p *Parser) varDeclaration() (Statement, error) {
	stmt := DeclarationStatement{}
	if token, ok := p.scan.match(TokenIdentifier); ok {
		stmt.name = token.Lexem
		stmt.pos = token.Position
	} else {
		return nil, NewSyntaxError(
			NewUnexpectedTokenError(TokenIdentifier.String(), token), token.Position,
		)
	}

	if _, ok := p.scan.match(TokenEqual); ok {
		var err error
		stmt.expr, err = p.expression()
		if err != nil {
			return nil, err
		}
	} else {
		stmt.expr = NilExpression{pos: stmt.pos}
	}

	if token, ok := p.scan.match(TokenSemicolon); !ok {
		return nil, NewSyntaxError(
			NewUnexpectedTokenError(TokenSemicolon.String(), token), token.Position,
		)
	}

	return stmt, nil
}

func (p *Parser) statement() (Statement, error) {
	if token, ok := p.scan.match(TokenPrint); ok {
		log.Debug().Msgf("(parser) %s ... ;", token.Lexem)
		return p.printStatement(token.Position)
	}
	return p.expressionStatement()
}

func (p *Parser) printStatement(pos Position) (Statement, error) {
	expr, err := p.expression()
	if err != nil {
		return nil, err
	}
	stmt := DeclarationStatement{expr: expr, pos: pos}

	if token, ok := p.scan.match(TokenSemicolon); !ok {
		return nil, NewSyntaxError(
			NewUnexpectedTokenError(TokenSemicolon.String(), token), token.Position,
		)
	}

	p.skipComments()
	return stmt, nil
}

func (p *Parser) expressionStatement() (Statement, error) {
	log.Debug().Msg("(parser) ... ;")
	expr, err := p.expression()
	if err != nil {
		return nil, err
	}
	stmt := ExpressionStatement{expr: expr, pos: expr.Position()}
	if token, ok := p.scan.match(TokenSemicolon); !ok {
		return nil, NewSyntaxError(
			NewUnexpectedTokenError(TokenSemicolon.String(), token), token.Position,
		)
	}
	p.skipComments()
	return stmt, nil
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
			// TODO: should this be a syntax error?
			return nil, err
		}
		log.Debug().Msgf("(parser) %s %s ...", expr, op)
		right, err := p.comparison()
		if err != nil {
			return nil, err
		}
		expr = BinaryExpression{op: op, left: expr, right: right, pos: expr.Position()}
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
			// TODO: should this be a syntax error?
			return nil, err
		}
		log.Debug().Msgf("(parser) %s %s ...", expr, op)
		right, err := p.term()
		if err != nil {
			return nil, err
		}
		expr = BinaryExpression{op: op, left: expr, right: right, pos: expr.Position()}
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
			// TODO: should this be a syntax error?
			return nil, err
		}
		log.Debug().Msgf("(parser) %s %s ...", expr, op)
		right, err := p.factor()
		if err != nil {
			return nil, err
		}
		expr = BinaryExpression{op: op, left: expr, right: right, pos: expr.Position()}
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
			// TODO: should this be a syntax error?
			return nil, err
		}
		log.Debug().Msgf("(parser) %s %s ...", expr, op)
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		expr = BinaryExpression{op: op, left: expr, right: right, pos: expr.Position()}
	}
	return expr, nil
}

func (p *Parser) unary() (Expression, error) {
	if token, ok := p.scan.match(TokenBang, TokenMinus); ok {
		op, err := token.Operator()
		if err != nil {
			// TODO: should this be a syntax error?
			return nil, err
		}
		log.Debug().Msgf("(parser) %s ...", op)
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		return UnaryExpression{op: op, right: right, pos: token.Position}, nil
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

	token := p.scan.peek()
	return nil, NewSyntaxError(NewMissingTerminalError(token), token.Position)
}

func (p *Parser) literal() (Expression, error) {
	if token, ok := p.scan.match(TokenNumber); ok {
		log.Debug().Msgf("(parser) terminal: %s", token)
		n, err := strconv.ParseFloat(token.Lexem, 64)
		if err != nil {
			return nil, NewSyntaxError(NewNumberConversionError(err, token), token.Position)
		}
		return NumericExpression{value: n, pos: token.Position}, nil
	} else if token, ok := p.scan.match(TokenString); ok {
		log.Debug().Msgf("(parser) terminal: %s", token)
		return StringExpression{value: token.Lexem, pos: token.Position}, nil
	} else if token, ok := p.scan.match(TokenTrue); ok {
		log.Debug().Msgf("(parser) terminal: %s", token)
		return BooleanExpression{value: true, pos: token.Position}, nil
	} else if token, ok := p.scan.match(TokenFalse); ok {
		log.Debug().Msgf("(parser) terminal: %s", token)
		return BooleanExpression{value: false, pos: token.Position}, nil
	} else if token, ok := p.scan.match(TokenNil); ok {
		log.Debug().Msgf("(parser) terminal: %s", token)
		return NilExpression{pos: token.Position}, nil
	} else if token, ok := p.scan.match(TokenIdentifier); ok {
		log.Debug().Msgf("(parser) terminal: %s", token)
		// TODO:
		panic("NYI")
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
			return nil, NewSyntaxError(
				NewUnexpectedTokenError(TokenRightParen.String(), token), token.Position,
			)
		}
		return GroupingExpression{expr: expr, pos: token.Position}, nil
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
		return eofToken
	}
	return s.tokens[s.offset]
}

func (s *tokenScanner) advance() Token {
	if s.done() {
		return eofToken
	}
	s.offset += 1
	return s.tokens[s.offset-1]
}

func (s *tokenScanner) match(ts ...TokenType) (Token, bool) {
	token := s.peek()
	for _, t := range ts {
		if token.Type == t {
			token = s.advance()
			return token, true
		}
	}
	return token, false
}
