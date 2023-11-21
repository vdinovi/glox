package lox

import (
	"strconv"

	"github.com/rs/zerolog/log"
)

func Parse(tokens []Token) ([]Statement, error) {
	parser := NewParser(tokens)
	return parser.Parse()
}

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
	for {
		if p.skipComments(); p.done() {
			log.Debug().Msg("(parser) done")
			break
		}
		stmt, err := p.declaration()

		if err != nil {
			log.Error().Msgf("(parser) error: %s", err)
			return nil, err
		}
		log.Debug().Msgf("(parser) statement: %s", stmt)
		stmts = append(stmts, stmt)
	}
	return stmts, nil
}

func (p *Parser) done() bool {
	_, eof := p.scan.match(TokenEOF)
	return eof
}

func (p *Parser) declaration() (Statement, error) {
	if fn, ok := p.scan.match(TokenFun); ok {
		return p.funcDeclaration(fn.Position)
	}
	if vr, ok := p.scan.match(TokenVar); ok {
		return p.varDeclaration(vr.Position)
	}
	return p.statement()
}

func (p *Parser) funcDeclaration(pos Position) (*FunctionStatement, error) {
	stmt := FunctionStatement{pos: pos}
	id, ok := p.scan.match(TokenIdentifier)
	if !ok {
		return nil, NewSyntaxError(
			NewUnexpectedTokenError(TokenIdentifier.String(), id), id.Position,
		)
	}
	stmt.name = id.Lexem
	if lparen, ok := p.scan.match(TokenLeftParen); !ok {
		return nil, NewSyntaxError(
			NewUnexpectedTokenError(TokenLeftParen.String(), lparen), lparen.Position,
		)
	}
	for {
		if _, ok := p.scan.match(TokenRightParen); ok {
			break
		}
		id, ok := p.scan.match(TokenIdentifier)
		if !ok {
			return nil, NewSyntaxError(
				NewUnexpectedTokenError(TokenIdentifier.String(), id), id.Position,
			)
		}
		stmt.params = append(stmt.params, id.Lexem)
		comma, ok := p.scan.match(TokenComma)
		if ok {
			continue
		}
		rparen, ok := p.scan.match(TokenRightParen)
		if ok {
			break
		}
		return nil, NewSyntaxError(
			NewUnexpectedTokenError(TokenRightParen.String(), rparen), comma.Position,
		)
	}
	lbrace, ok := p.scan.match(TokenLeftBrace)
	if !ok {
		return nil, NewSyntaxError(
			NewUnexpectedTokenError(TokenLeftBrace.String(), lbrace), lbrace.Position,
		)
	}
	for {
		if p.done() {
			return nil, NewSyntaxError(
				NewUnexpectedTokenError(TokenRightBrace.String(), eofToken), lbrace.Position,
			)
		}
		if _, ok := p.scan.match(TokenRightBrace); ok {
			return &stmt, nil
		}
		bodyStmt, err := p.declaration()
		if err != nil {
			return nil, err
		}
		stmt.body = append(stmt.body, bodyStmt)
	}
}

func (p *Parser) varDeclaration(pos Position) (*DeclarationStatement, error) {
	stmt := DeclarationStatement{pos: pos}
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
		stmt.expr = &NilExpression{pos: stmt.pos}
	}

	if token, ok := p.scan.match(TokenSemicolon); !ok {
		return nil, NewSyntaxError(
			NewUnexpectedTokenError(TokenSemicolon.String(), token), token.Position,
		)
	}

	return &stmt, nil
}

func (p *Parser) statement() (Statement, error) {
	if if_, ok := p.scan.match(TokenIf); ok {
		return p.condStatement(if_.Position)
	}
	if print, ok := p.scan.match(TokenPrint); ok {
		return p.printStatement(print.Position)
	}
	if while, ok := p.scan.match(TokenWhile); ok {
		return p.whileStatement(while.Position)
	}
	if for_, ok := p.scan.match(TokenFor); ok {
		return p.forStatement(for_.Position)
	}
	if lbrace, ok := p.scan.match(TokenLeftBrace); ok {
		return p.blockStatement(lbrace.Position)
	}
	return p.expressionStatement()
}

func (p *Parser) condStatement(pos Position) (*ConditionalStatement, error) {
	var err error
	stmt := ConditionalStatement{pos: pos}
	stmt.expr, err = p.condition()
	if err != nil {
		return nil, err
	}
	stmt.thenBranch, err = p.statement()
	if err != nil {
		return nil, err
	}
	if _, ok := p.scan.match(TokenElse); !ok {
		stmt.elseBranch = nil
		return &stmt, nil
	}
	stmt.elseBranch, err = p.statement()
	if err != nil {
		return nil, err
	}
	return &stmt, nil
}

func (p *Parser) whileStatement(pos Position) (*WhileStatement, error) {
	var err error
	stmt := WhileStatement{pos: pos}
	stmt.expr, err = p.condition()
	if err != nil {
		return nil, err
	}
	stmt.body, err = p.statement()
	if err != nil {
		return nil, err
	}
	return &stmt, nil
}

func (p *Parser) forStatement(pos Position) (*ForStatement, error) {
	var err error
	stmt := ForStatement{pos: pos}
	if lparen, ok := p.scan.match(TokenLeftParen); !ok {
		return nil, NewSyntaxError(
			NewUnexpectedTokenError(TokenLeftParen.String(), lparen), lparen.Position,
		)
	}
	if _, ok := p.scan.match(TokenSemicolon); ok {
		stmt.init = nil
	} else if var_, ok := p.scan.match(TokenVar); ok {
		stmt.init, err = p.varDeclaration(var_.Position)
		if err != nil {
			return nil, err
		}
	} else {
		stmt.init, err = p.expressionStatement()
		if err != nil {
			return nil, err
		}

	}
	if _, ok := p.scan.match(TokenSemicolon); ok {
		stmt.cond = nil
	} else {
		stmt.cond, err = p.expression()
		if err != nil {
			return nil, err
		}
		if semicolon, ok := p.scan.match(TokenSemicolon); !ok {
			return nil, NewSyntaxError(
				NewUnexpectedTokenError(TokenSemicolon.String(), semicolon), semicolon.Position,
			)
		}
	}
	if _, ok := p.scan.match(TokenRightParen); ok {
		stmt.incr = nil
	} else {
		stmt.incr, err = p.expression()
		if err != nil {
			return nil, err
		}
		if rparen, ok := p.scan.match(TokenRightParen); !ok {
			return nil, NewSyntaxError(
				NewUnexpectedTokenError(TokenRightParen.String(), rparen), rparen.Position,
			)
		}
	}
	stmt.body, err = p.statement()
	if err != nil {
		return nil, err
	}
	return &stmt, nil
}

func (p *Parser) printStatement(pos Position) (*PrintStatement, error) {
	expr, err := p.expression()
	if err != nil {
		return nil, err
	}
	stmt := PrintStatement{expr: expr, pos: pos}

	if token, ok := p.scan.match(TokenSemicolon); !ok {
		return nil, NewSyntaxError(
			NewUnexpectedTokenError(TokenSemicolon.String(), token), token.Position,
		)
	}

	return &stmt, nil
}

func (p *Parser) blockStatement(pos Position) (*BlockStatement, error) {
	block := BlockStatement{stmts: make([]Statement, 0), pos: pos}
	for {
		if _, ok := p.scan.match(TokenRightBrace); ok || p.scan.done() {
			return &block, nil
		}
		stmt, err := p.declaration()
		if err != nil {
			return nil, err
		}
		block.stmts = append(block.stmts, stmt)
	}
}

func (p *Parser) expressionStatement() (*ExpressionStatement, error) {
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
	return &stmt, nil
}

func (p *Parser) condition() (expr Expression, err error) {
	if lparen, ok := p.scan.match(TokenLeftParen); !ok {
		return nil, NewSyntaxError(
			NewUnexpectedTokenError(TokenLeftParen.String(), lparen), lparen.Position,
		)
	}
	expr, err = p.expression()
	if err != nil {
		return nil, err
	}
	if rparen, ok := p.scan.match(TokenRightParen); !ok {
		return nil, NewSyntaxError(
			NewUnexpectedTokenError(TokenLeftParen.String(), rparen), rparen.Position,
		)
	}
	return expr, nil
}

func (p *Parser) expression() (Expression, error) {
	expr, err := p.assignment()
	return expr, err
}

func (p *Parser) assignment() (Expression, error) {
	expr, err := p.or()
	if err != nil {
		return nil, err
	}
	if _, ok := p.scan.match(TokenEqual); ok {
		right, err := p.assignment()
		if err != nil {
			return nil, err
		}
		if left, ok := expr.(*VariableExpression); ok {
			return &AssignmentExpression{name: left.name, right: right, pos: left.Position()}, nil
		}
		return nil, NewSyntaxError(NewInvalidAssignmentTargetError(expr.String()), expr.Position())
	}
	return expr, nil
}

func (p *Parser) or() (Expression, error) {
	expr, err := p.and()
	if err != nil {
		return nil, err
	}
	for {
		token, ok := p.scan.match(TokenOr)
		if !ok {
			break
		}
		right, err := p.and()
		if err != nil {
			return nil, err
		}
		op, err := token.Operator()
		if err != nil {
			return nil, err
		}
		expr = &BinaryExpression{op: op, left: expr, right: right, pos: expr.Position()}
	}
	return expr, nil
}

func (p *Parser) and() (Expression, error) {
	expr, err := p.equality()
	if err != nil {
		return nil, err
	}
	for {
		token, ok := p.scan.match(TokenAnd)
		if !ok {
			break
		}
		right, err := p.and()
		if err != nil {
			return nil, err
		}
		op, err := token.Operator()
		if err != nil {
			return nil, err
		}
		expr = &BinaryExpression{op: op, left: expr, right: right, pos: expr.Position()}
	}
	return expr, nil
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
			return nil, err
		}
		right, err := p.comparison()
		if err != nil {
			return nil, err
		}
		expr = &BinaryExpression{op: op, left: expr, right: right, pos: expr.Position()}
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
			return nil, err
		}
		right, err := p.term()
		if err != nil {
			return nil, err
		}
		expr = &BinaryExpression{op: op, left: expr, right: right, pos: expr.Position()}
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
			return nil, err
		}
		right, err := p.factor()
		if err != nil {
			return nil, err
		}
		expr = &BinaryExpression{op: op, left: expr, right: right, pos: expr.Position()}
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
			return nil, err
		}
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		expr = &BinaryExpression{op: op, left: expr, right: right, pos: expr.Position()}
	}
	return expr, nil
}

func (p *Parser) unary() (Expression, error) {
	if token, ok := p.scan.match(TokenBang, TokenMinus, TokenPlus); ok {
		op, err := token.Operator()
		if err != nil {
			return nil, err
		}
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		return &UnaryExpression{op: op, right: right, pos: token.Position}, nil
	}
	return p.call()
}

func (p *Parser) call() (expr Expression, err error) {
	expr, err = p.primary()
	if err != nil {
		return nil, err
	}
	for {
		if lparen, ok := p.scan.match(TokenLeftParen); !ok {
			break
		} else {
			expr, err = p.finishCall(lparen.Position, expr)
			if err != nil {
				return nil, err
			}
		}
	}
	return expr, err
}

const maxArguments = 255

func (p *Parser) finishCall(pos Position, callee Expression) (Expression, error) {
	expr := CallExpression{callee: callee, pos: pos}
	if _, ok := p.scan.match(TokenRightParen); !ok {
		for {
			arg, err := p.expression()
			if err != nil {
				return nil, err
			}
			if nargs := len(expr.args); nargs >= maxArguments {
				return nil, NewSyntaxError(
					NewMaximumArgumentCountExceededError(nargs+1, maxArguments),
					arg.Position(),
				)
			}

			expr.args = append(expr.args, arg)
			if _, ok := p.scan.match(TokenComma); !ok {
				if rparen, ok := p.scan.match(TokenRightParen); !ok {
					return nil, NewSyntaxError(
						NewUnexpectedTokenError(TokenRightParen.String(), rparen),
						rparen.Position,
					)
				}
				break
			}
		}
	}
	return &expr, nil
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
		n, err := strconv.ParseFloat(token.Lexem, 64)
		if err != nil {
			return nil, NewSyntaxError(NewNumberConversionError(err, token), token.Position)
		}
		return &NumericExpression{value: n, pos: token.Position}, nil
	} else if token, ok := p.scan.match(TokenString); ok {
		return &StringExpression{value: token.Lexem, pos: token.Position}, nil
	} else if token, ok := p.scan.match(TokenTrue); ok {
		return &BooleanExpression{value: true, pos: token.Position}, nil
	} else if token, ok := p.scan.match(TokenFalse); ok {
		return &BooleanExpression{value: false, pos: token.Position}, nil
	} else if token, ok := p.scan.match(TokenNil); ok {
		return &NilExpression{pos: token.Position}, nil
	} else if token, ok := p.scan.match(TokenIdentifier); ok {
		return &VariableExpression{name: token.Lexem, pos: token.Position}, nil
	}

	return nil, nil
}

func (p *Parser) skipComments() {
	for {
		if _, ok := p.scan.match(TokenComment); ok {
			continue
		}
		break
	}
}

func (p *Parser) grouping() (Expression, error) {
	if token, ok := p.scan.match(TokenLeftParen); ok {
		expr, err := p.expression()
		if err != nil {
			return nil, err
		}
		if _, ok := p.scan.match(TokenRightParen); !ok {
			return nil, NewSyntaxError(
				NewUnexpectedTokenError(TokenRightParen.String(), token), token.Position,
			)
		}
		return &GroupingExpression{expr: expr, pos: token.Position}, nil
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
