package main

import (
    "fmt"
    "strings"
)

type ASTNode interface{}

type SelectQuery struct {
    ColumnName  []string
    TableName string
    Condition string
}

func main() {
    input := "SELECT id, name FROM users WHERE id = 1;"
    tokens := tokenize(input)
	// fmt.Printf("%#v\n",tokens)
	for i := 0; i < len(tokens); i++{
		fmt.Printf("%s\n",tokens[i])
	}
    ast := parse(tokens)
    fmt.Printf("%#v\n", ast)
}

func tokenize(input string) []string {
    var tokens []string
    for _, word := range strings.FieldsFunc(input, func(r rune) bool {
        return r == ' ' || r == ',' || r == ';'
    }) {
        tokens = append(tokens, word)
    }
    return tokens
}

func parse(tokens []string) ASTNode {
    var query SelectQuery
    var fields []string
    for i := 0; i < len(tokens); i++ {
        switch strings.ToUpper(tokens[i]) {
        case "SELECT":
            i++
            for ; i < len(tokens) && tokens[i] != "FROM"; i++ {
                if tokens[i] != "," {
                    fields = append(fields, tokens[i])
                }
            }
			i--
            query.ColumnName = fields
        case "FROM":
            i++
            query.TableName = tokens[i]
        case "WHERE":
            i++
            query.Condition = tokens[i]
        }
    }
    return query
}
