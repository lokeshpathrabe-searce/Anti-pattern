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


	// Check if the query contains an IN or NOT IN clause
	var inClause *sqlparser.ComparisonExpr
	sqlparser.Walk(func(node sqlparser.SQLNode) (kontinue bool, err error) {
		if comparison, ok := node.(*sqlparser.ComparisonExpr); ok {
			// fmt.Println("Operator:", comparison.Operator)
			if comparison.Operator == "in" || comparison.Operator == "not in" {
				// fmt.Println("Found IN or NOT IN clause")
				inClause = comparison
				return false, nil
			}
		}
		return true, nil
	}, stmt)
	
	// If IN or NOT IN clause found, modify the subquery/subqueries
	if inClause != nil {
		sqlparser.Walk(func(node sqlparser.SQLNode) (kontinue bool, err error) {
			switch n := node.(type) {
			case *sqlparser.Subquery:
				if needsModification(n) {
					modifySubquery(n)
				}
			}
			return true, nil
		}, stmt)
	}

	// Print the modified query
	fmt.Println("Modified Query:")
	fmt.Println(sqlparser.String(stmt))

}
 
func modifySubquery(subquery *sqlparser.Subquery) {
	if selectStmt, ok := subquery.Select.(*sqlparser.Select); ok {
		// Add Distinct to the subquery
		selectStmt.Distinct = sqlparser.DistinctStr

		// Extract column names from the subquery
		var groupByColumns []sqlparser.Expr
		for _, expr := range selectStmt.SelectExprs {
			groupByColumns = append(groupByColumns, expr.(*sqlparser.AliasedExpr).Expr)
		}

		// Add GROUP BY clause to the subquery
		selectStmt.GroupBy = groupByColumns
	}
}

func needsModification(subquery *sqlparser.Subquery) bool {
	if selectStmt, ok := subquery.Select.(*sqlparser.Select); ok {
		// Check if the subquery already has DISTINCT or GROUP BY
		return selectStmt.Distinct == "" && len(selectStmt.GroupBy) == 0
	}
	return false
}
