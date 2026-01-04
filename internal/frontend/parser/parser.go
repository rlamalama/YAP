package parser

import (
	"fmt"
	"os"
)

type Parser struct {
	filename string
	tokens   []*Token
	pos      int
}

func NewParser(file string) *Parser {
	return &Parser{
		filename: file,
		tokens:   []*Token{},
		pos:      0,
	}
}

func (p *Parser) peek() *Token {
	if p.pos >= len(p.tokens) {
		return &Token{Kind: TokenEOF}

	}
	return p.tokens[p.pos]
}

func (p *Parser) next() *Token {
	tok := p.peek()
	p.pos++
	return tok
}

func (p *Parser) expect(kind TokenKind) (*Token, error) {
	tok := p.next()
	if tok == nil {
		return nil, fmt.Errorf("expected %v, no more tokens", kind)
	}

	if tok.Kind != kind {
		return tok, fmt.Errorf(
			"expected %v, got %v at line %d",
			kind, tok.Kind, tok.Line,
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

	lexer := NewLexer(file)

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
	for p.peek().Kind != TokenEOF {

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
	if _, err := p.expect(TokenDash); err != nil {
		return nil, err
	}

	key, err := p.expect(TokenIdentifier)
	if err != nil {
		return nil, err
	}

	if _, err := p.expect(TokenColon); err != nil {
		return nil, err
	}

	switch key.Value {
	case "print":
		return p.parsePrint()
	default:
		return nil, fmt.Errorf(
			"unknown statement %q at line %d",
			key.Value,
			key.Line,
		)
	}
}

func (p *Parser) parsePrint() (Stmt, error) {
	val, err := p.expect(TokenScalar)
	if err != nil {
		return nil, err
	}

	// Optional: consume newline if present
	if p.peek().Kind == TokenEOL {
		p.next()
	}

	return PrintStmt{
		Value: val.Value,
	}, nil
}
