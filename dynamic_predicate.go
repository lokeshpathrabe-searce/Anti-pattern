package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/xwb1989/sqlparser"
)
 
func main() {
	var queryFile string
	flag.StringVar(&queryFile, "file", "", "Specify the file containing the query")
 
	flag.Parse()
 
	if queryFile == "" {
		fmt.Println("Please provide a file containing the query using the -file flag")
		return
	}
 
	file, err := os.Open(queryFile)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()
 
	// Read the query from the file
	var queryLines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		queryLines = append(queryLines, line)
	}
 
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
		return
	}
 
	// Join the query lines into a single string
	query := strings.Join(queryLines, "\n")
 
	stmt, err := sqlparser.Parse(query)
	if err != nil {
		fmt.Println("Error parsing SQL query:", err)
		return
	}


	// Walk through the AST to find the subquery in the WHERE clause
	sqlparser.Walk(func(node sqlparser.SQLNode) (kontinue bool, err error) {
		switch n := node.(type) {
		case *sqlparser.Subquery:
			// If a subquery is found, transform it into a JOIN statement
			transformSubqueryToJoin(n,stmt)
			return false, nil
		}
		return true, nil
	}, stmt)

	// Print the modified query
	fmt.Println("Modified Query:")
	fmt.Println(sqlparser.String(stmt))

}
 
func transformSubqueryToJoin(subquery *sqlparser.Subquery,stmt sqlparser.Statement) {
	mainSelect:= stmt.(*sqlparser.Select)
	subSelect := subquery.Select.(*sqlparser.Select)

	subTableName := subSelect.From[0].(*sqlparser.AliasedTableExpr).Expr.(sqlparser.TableName)
	mainTableName := mainSelect.From[0].(*sqlparser.AliasedTableExpr).Expr.(sqlparser.TableName)

	subColName := subSelect.SelectExprs[0].(*sqlparser.AliasedExpr).Expr.(*sqlparser.ColName).Name.String()
	// fmt.Printf("%#v\n",subSelect.SelectExprs[0].(*sqlparser.AliasedExpr).Expr.(*sqlparser.ColName).Name.String())

	// Create JOIN statement
	joinExpr := &sqlparser.JoinTableExpr{
		LeftExpr: mainSelect.From[0].(*sqlparser.AliasedTableExpr),
		Join: "join",
		RightExpr: subSelect.From[0].(*sqlparser.AliasedTableExpr),
	}

	leftCondColName := fmt.Sprintf("%s.%s",mainTableName.Name.String(),subColName)
	rightCondColName := fmt.Sprintf("%s.%s",subTableName.Name.String(),subColName)

	// Build the ON condition
	onCondition := &sqlparser.ComparisonExpr{
		Operator: sqlparser.EqualStr,
		Left:     &sqlparser.ColName{Name: sqlparser.NewColIdent(leftCondColName)},
		Right:    &sqlparser.ColName{Name: sqlparser.NewColIdent(rightCondColName)},
	}
	joinExpr.Condition = sqlparser.JoinCondition{On: onCondition}
	// fmt.Printf("%#v\n",mainSelect.From[0].(*sqlparser.JoinTableExpr))

	mainSelect.From = sqlparser.TableExprs{joinExpr}
	// fmt.Printf("%#v\n",mainSelect.Where.Expr.(*sqlparser.AndExpr).Left.(*sqlparser.AndExpr))

	// Check if there is only one subquery with IN clause in the WHERE clause
	if containsOnlySubquery(mainSelect.Where.Expr) {
		mainSelect.Where = nil
	}else{
		// Remove the condition AND column IN (Subquery)
		newConditions := make([]sqlparser.Expr, 0)
		for _, expr := range splitAndExpressions(mainSelect.Where.Expr) {
			if !containsSubquery(expr) {
				newConditions = append(newConditions, expr)
			}
		}

		// Update the WHERE clause with the new conditions
		mainSelect.Where = &sqlparser.Where{
			Type:"where",
			Expr: joinAndExpressions(newConditions),
		}
	}

	// fmt.Printf("%#v\n",mainSelect)
	// fmt.Printf("%#v\n",stmt.(*sqlparser.Select).Where)
}

// containsSubquery checks if the expression contains a subquery
func containsSubquery(expr sqlparser.Expr) bool {
	switch node := expr.(type) {
	case *sqlparser.ComparisonExpr:
		if node.Operator == "in" || node.Operator == "not in" {
			_, ok := node.Right.(*sqlparser.Subquery)
			return ok
		}
	case *sqlparser.ParenExpr:
		return containsSubquery(node.Expr)
	}
	return false
}

// splitAndExpressions splits the AND expressions in an expression tree
func splitAndExpressions(expr sqlparser.Expr) []sqlparser.Expr {
	switch node := expr.(type) {
	case *sqlparser.AndExpr:
		return append(splitAndExpressions(node.Left), splitAndExpressions(node.Right)...)
	default:
		return []sqlparser.Expr{expr}
	}
}

// joinAndExpressions joins multiple expressions with AND operators
func joinAndExpressions(exprs []sqlparser.Expr) sqlparser.Expr {
	if len(exprs) == 1 {
		return exprs[0]
	}
	andExpr := &sqlparser.AndExpr{}
	andExpr.Left = joinAndExpressions(exprs[:len(exprs)-1])
	andExpr.Right = exprs[len(exprs)-1]
	return andExpr
}

// containsOnlySubquery checks if the expression contains only one subquery with IN clause
func containsOnlySubquery(expr sqlparser.Expr) bool {
	switch node := expr.(type) {
	case *sqlparser.ComparisonExpr:
		if node.Operator == "in" || node.Operator == "not in" {
			_, ok := node.Right.(*sqlparser.Subquery)
			return ok
		}
	case *sqlparser.ParenExpr:
		return containsOnlySubquery(node.Expr)
	}
	return false
}