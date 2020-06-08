package parser

import (
	"github.com/antlr/antlr4/runtime/Go/antlr"
	"github.com/liturgiko/doxa/pkg/template"
	lml "gitlab.com/ocmc/liturgiko/lml-go/parser"
)

/**
LML provides access to a lexer and parser for the input stream.
Use NewLMLParser to get an instance.
 */
type LML struct {
	Lexer *lml.LMLLexer
	Parser *lml.LMLParser
	Listener *LMLListener
	ErrorListener *LMLErrorListener
	TemplateID string
	Input string
}
// NewLMLParser provides an LML instance with a lexer and parser initialized for the input stream.
// Upon completion of the parsing, a slice of parse errors will be available, and a template.ATEM will be populated.
// If the ATEM will be used for generating HTML or PDF or something else, then genLibs and resolver may not be nil.
// If genLibs and resolver are nil, the purpose of calling the parser is presumed to be simply to validate the syntax of the input.
// If genLibs and resolver are not nil, they will be used to populate the Values and Version acronyms in the ATEM, so it can be used for generation.
// Returns an error if the database path is invalid.
// templateID is the id of the template.
// genLibs provides the primary libraries and fallback libraries to use for generation. If
func NewLMLParser(templateID, input string, genLibs []template.GenLib, resolver Resolver) (*LML, error) {
	l := new(LML)
	var err error
	l.TemplateID = templateID
	l.Input = input
	l.Lexer = lml.NewLMLLexer(antlr.NewInputStream(input))
	stream := antlr.NewCommonTokenStream(l.Lexer, antlr.TokenDefaultChannel)
	l.Parser = lml.NewLMLParser(stream)
	l.Listener, err = NewLMLListener(genLibs, resolver)
	if err != nil {
		return nil, err
	}
	l.ErrorListener = NewLMLErrorListener(templateID)
	l.Parser.RemoveErrorListeners()
	l.Parser.AddErrorListener(l.ErrorListener)
	return l, nil
}
// Performs a walk on the template parse tree starting at the root and going down recursively with depth-first search. On each node, EnterRule is called before recursively walking down into child nodes, then ExitRule is called after the recursive call to wind up.
// Returns the Errors found in the input stream by the parser.
// Note that antlr only returns an error of a specify type once.
// So, there may be more errors than are reported by this function.
func (l *LML) WalkTemplate()  []ParseError {
	return l.walk(l.Parser.Template())
}
// Performs a walk on the given parse tree starting at the root and going down recursively with depth-first search. On each node, EnterRule is called before recursively walking down into child nodes, then ExitRule is called after the recursive call to wind up.
// This local function is provided for test purposes.
// Returns the Errors found in the input stream by the parser.
// Note that antlr only returns an error of a specify type once.
// So, there may be more errors than are reported by this function.
func (l *LML) walk(t antlr.Tree)  []ParseError {
	antlr.ParseTreeWalkerDefault.Walk(l.Listener, t)
	return l.ErrorListener.Errors
}
// Tokens provides back information about each identified token. It will not indicate any errors.
func (l *LML) Tokens() []antlr.Token {
	var tokens []antlr.Token
	for {
		t := l.Lexer.NextToken()
		if t.GetTokenType() == antlr.TokenEOF {
			break
		}
		tokens = append(tokens,t)
	}
	return tokens
}
