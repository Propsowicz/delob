package parser

import "fmt"

func errorCannotGenerateParsedExpression(traceId string) error {
	return fmt.Errorf("cannot generate parsed expression from given tokens")
}

func errorWrongExpressionFormat(traceId, expression string) error {
	return fmt.Errorf("expression should ends with semicolon (;), got - %s", expression)
}

func errorCannotParseExpression(traceId, expression string) error {
	return fmt.Errorf("cannot parse expression - %s;", expression)
}

func errorInvalidNumberOfArguments(traceId, expression string, exprType ExpressionType, numberOfArgs int) error {
	return fmt.Errorf("invalid number of arguments in expression - %s;, function '%s' support exactly %d arguments, i.e.: %s",
		expression, exprType, numberOfArgs, exampleExpression(exprType))
}

func errorCannotUseUserKeyMoreThanOnce(key string) error {
	return fmt.Errorf("duplicate key '%s' detected - keys must be unique", key)
}

func exampleExpression(exprType ExpressionType) string {
	switch exprType {
	case AddPlayersCommandType:
		return "<ADD PLAYER 'Tom';> or <ADD PLAYERS ('Tom', 'Joe');>"
	case AddMatchCommandType:
		return "<SET WIN FOR 'Tom' AND LOSE FOR 'Joe';> or <SET WIN FOR ('Tom', 'Bob') AND LOSE FOR ('Joe', 'Jim');>"
	default:
		return ""
	}
}
