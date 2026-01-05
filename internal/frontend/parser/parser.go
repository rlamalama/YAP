package parser

import (
	"fmt"
	"os"
	"strconv"

	yaperror "github.com/rlamalama/YAP/internal/error"
	"github.com/rlamalama/YAP/internal/frontend/lexer"
)

type Parser struct {
	filename string
	tokens   []*lexer.Token
	pos      int
}

func NewParser(file string) *Parser {
	return &Parser{
		filename: file,
		tokens:   []*lexer.Token{},
		pos:      0,
	}
}

func (p *Parser) peek() *lexer.Token {
	if p.pos >= len(p.tokens) {
		return &lexer.Token{Kind: lexer.TokenEOF}

	}
	return p.tokens[p.pos]
}

func (p *Parser) next() *lexer.Token {
	tok := p.peek()
	p.pos++
	return tok
}

func (p *Parser) expect(kind lexer.TokenKind) (*lexer.Token, error) {
	tok := p.next()
	if tok == nil {
		return nil, yaperror.NewExpectedTokenError(p.filename, 0, 0, kind.String())
	}

	if tok.Kind != kind {
		return tok, yaperror.NewUnexpectedTokenError(
			p.filename, tok.Line, tok.Col,
			tok.Kind.String(), kind.String(),
		)
	}
	return tok, nil
}

func (p *Parser) Parse() (*Program, error) {
	file, err := os.Open(p.filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	lexer := lexer.NewLexer(file, p.filename)

	p.tokens, err = lexer.Lex()
	if err != nil {
		return nil, err
	}

	return p.parseProgram()
}

func (p *Parser) parseProgram() (*Program, error) {
	prog := &Program{}

	if len(p.tokens) == 0 {
		return prog, nil
	}

	stmts := []Stmt{}
	for p.peek().Kind != lexer.TokenEOF {
		if p.peek().Kind == lexer.TokenComment {
			// pop comment token and the following new line
			_, err := p.expect(lexer.TokenComment)
			if err != nil {
				return nil, err
			}

			_, err = p.expect(lexer.TokenNewline)
			if err != nil {
				return nil, err
			}
			continue
		}

		stmt, err := p.parseStmt()
		if err != nil {
			return nil, err
		}
		stmts = append(stmts, stmt)
	}

	prog.Statements = stmts
	return prog, nil
}

func (p *Parser) parseStmt() (Stmt, error) {
	if _, err := p.expect(lexer.TokenDash); err != nil {
		return nil, err
	}

	key, err := p.expect(lexer.TokenKeyword)
	if err != nil {
		return nil, err
	}

	if _, err := p.expect(lexer.TokenColon); err != nil {
		return nil, err
	}

	switch key.Value {
	case lexer.KeywordPrint:
		return p.parsePrint()
	case lexer.KeywordSet:
		return p.parseSet()
	default:
		return nil, yaperror.NewUnknownStatementError(
			p.filename, key.Line, key.Col, key.Value,
		)
	}
}

func (p *Parser) parseValue() (Value, error) {
	// Skip any comment tokens (but not newlines - those are structural)
	for p.peek().Kind == lexer.TokenComment {
		p.next() // consume comment only
	}

	switch p.peek().Kind {
	case lexer.TokenIdentifier:
		tok := p.next()
		return &Identifier{Name: tok.Value}, nil

	case lexer.TokenString:
		tok := p.next()
		return &StringLiteral{Value: tok.Value}, nil

	case lexer.TokenNumerical:
		tok := p.next()
		num, err := strconv.Atoi(tok.Value)
		if err != nil {
			return nil, err
		}
		return &NumericLiteral{Value: num}, nil

	case lexer.TokenKeyword:
		tok := p.peek()
		// Handle boolean literals
		if tok.Value == lexer.KeywordTrue {
			p.next()
			return &BooleanLiteral{Value: true}, nil
		}
		if tok.Value == lexer.KeywordFalse {
			p.next()
			return &BooleanLiteral{Value: false}, nil
		}
		return nil, fmt.Errorf(
			"unexpected keyword %q at line %d, expected value",
			tok.Value,
			tok.Line,
		)

	default:
		tok := p.peek()
		return nil, fmt.Errorf(
			"expected value, got %v at line %d",
			tok.Kind,
			tok.Line,
		)
	}
}

func (p *Parser) parsePrint() (Stmt, error) {
	expr, err := p.parseExpr()
	if err != nil {
		return nil, err
	}

	// Skip any trailing comment before newline
	for p.peek().Kind == lexer.TokenComment {
		p.next() // consume comment only
	}

	if _, err := p.expect(lexer.TokenNewline); err != nil {
		return nil, err
	}

	return PrintStmt{
		Expr: expr,
	}, nil
}

func (p *Parser) parseExpr() (Value, error) {
	left, err := p.parseValue()
	if err != nil {
		return nil, err
	}

	// Check for binary operators
	for p.peek().Kind == lexer.TokenOperator {
		opTok := p.next()

		// Skip any comment tokens after operator (but not newlines - those are structural)
		for p.peek().Kind == lexer.TokenComment {
			p.next() // consume comment only
		}

		right, err := p.parseValue()
		if err != nil {
			return nil, err
		}
		left = &BinaryExpr{
			Left:     left,
			Operator: opTok.Value,
			Right:    right,
		}
	}

	return left, nil
}

func (p *Parser) parseSet() (Stmt, error) {
	val := p.next()
	switch val.Kind {
	case lexer.TokenNewline:

		_, err := p.expect(lexer.TokenIndent)
		if err != nil {
			return nil, err
		}

		assignments := []*Assignment{}
		for p.peek().Kind != lexer.TokenDedent {
			// Skip any comment lines within the set block
			for p.peek().Kind == lexer.TokenComment {
				p.next() // consume comment
				if p.peek().Kind == lexer.TokenNewline {
					p.next() // consume newline after comment
				}
			}

			// Check again after skipping comments
			if p.peek().Kind == lexer.TokenDedent {
				break
			}

			// Start of sets STMT
			_, err = p.expect(lexer.TokenDash)
			if err != nil {
				return nil, err
			}

			key, err := p.expect(lexer.TokenIdentifier)
			if err != nil {
				return nil, err
			}

			_, err = p.expect(lexer.TokenColon)
			if err != nil {
				return nil, err
			}

			expr, err := p.parseExpr()
			if err != nil {
				return nil, err
			}

			_, err = p.expect(lexer.TokenNewline)
			if err != nil {
				return nil, err
			}
			assignments = append(assignments, &Assignment{
				Name: key.Value,
				Expr: expr,
			})
		}

		// END OF STMT
		_, err = p.expect(lexer.TokenDedent)
		if err != nil {
			return nil, err
		}

		return SetStmt{
			Assignment: assignments,
		}, nil

	default:
		return nil, yaperror.NewUnexpectedTokenError(
			p.filename, val.Line, val.Col,
			val.Kind.String(), fmt.Sprintf("%s or", lexer.TokenNewline.String()),
		)
	}
}
