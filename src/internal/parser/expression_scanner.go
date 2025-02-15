package parser

import (
	"fmt"
	"strings"
)

type ExpressionScanner struct {
	ExpressionType ExpressionType
	traceId        string
	tokens         []Token
}

type Token struct {
	Token TokenType
	Value []string
}

type ExpressionType string

const (
	AddPlayersCommandType ExpressionType = "AddPlayersCommandType"
	AddMatchCommandType   ExpressionType = "AddMatchCommandType"
	SelectQueryType       ExpressionType = "SelectQueryType"
)

type TokenType string

const (
	AddPlayer     TokenType = "AddPlayer"
	AddPlayers    TokenType = "AddPlayers"
	SetWin        TokenType = "SetWin"
	SetLose       TokenType = "SetLose"
	SetDraw       TokenType = "SetDraw"
	SelectPlayers TokenType = "SelectPlayers"
	OrderByAsc    TokenType = "OrderByAsc"
	OrderByDesc   TokenType = "OrderByDesc"
)

func newExpressionScanner(traceId string) ExpressionScanner {
	return ExpressionScanner{}
}

func (sc *ExpressionScanner) getLogicalTokens() (ExpressionType, []Token, error) {
	if sc.ExpressionType == "" {
		return "", nil, fmt.Errorf("expression type is not determined")
	}
	return sc.ExpressionType, sc.tokens, nil
}

func (sc *ExpressionScanner) scanRawExpression(expression string) error {
	sanitazedExpression, err := sanitazeExpression(sc.traceId, expression)
	if err != nil {
		return err
	}

	switch {
	case isMatch(add_player, sanitazedExpression):
		return sc.tryTokenizeAddPlayer(sanitazedExpression)
	case isMatch(add_players, sanitazedExpression):
		return sc.tryTokenizeAddPlayers(sanitazedExpression)
	case isMatch(add_decisive_match, sanitazedExpression):
		return sc.tryTokenizeSetWinAndLose(sanitazedExpression)
	case isMatch(add_draw_match, sanitazedExpression):
		return sc.tryTokenizeSetDraw(sanitazedExpression)
	case isMatch(select_players, sanitazedExpression):
		return sc.tryTokenizeSelectPlayers(sanitazedExpression)
	default:
		return errorCannotParseExpression(sc.traceId, expression)
	}
}

func (sc *ExpressionScanner) tryTokenizeAddPlayer(sanitazedExpression string) error {
	numOfMatches, tokenValue := extractKeyFromParanthesis(sanitazedExpression)
	if numOfMatches != 1 {
		return errorInvalidNumberOfArguments(sc.traceId, sanitazedExpression, AddPlayersCommandType, 1)
	}
	sc.ExpressionType = AddPlayersCommandType
	sc.tokens = append(sc.tokens, Token{
		AddPlayer,
		tokenValue,
	})
	return nil
}

func (sc *ExpressionScanner) tryTokenizeAddPlayers(sanitazedExpression string) error {
	numOfMatches, tokenValue := extractKeysFromBrackets(sanitazedExpression)
	if numOfMatches != 1 {
		return errorInvalidNumberOfArguments(sc.traceId, sanitazedExpression, AddPlayersCommandType, 1)
	}
	sc.ExpressionType = AddPlayersCommandType
	sc.tokens = append(sc.tokens, Token{
		AddPlayers,
		tokenValue[0],
	})
	return nil
}

func (sc *ExpressionScanner) tryTokenizeSetWinAndLose(sanitazedExpression string) error {
	setWinIdx := strings.Index(strings.ToLower(sanitazedExpression), "win")
	setLoseIdx := strings.Index(strings.ToLower(sanitazedExpression), "lose")
	setWinFirst := setWinIdx < setLoseIdx

	numOfBracketMatches, tokenBracketValue := extractKeysFromBrackets(sanitazedExpression)
	if numOfBracketMatches == 2 {
		sc.setWinAndLoseTokens(setWinFirst, tokenBracketValue[0], tokenBracketValue[1])
		return nil
	}

	numOfParanthesisMatches, tokenParanthesisValue := extractKeyFromParanthesis(sanitazedExpression)
	if numOfParanthesisMatches == 2 {
		sc.setWinAndLoseTokens(setWinFirst, tokenParanthesisValue[:1], tokenParanthesisValue[1:])
		return nil
	}
	return errorInvalidNumberOfArguments(sc.traceId, sanitazedExpression, AddMatchCommandType, 2)
}

func (sc *ExpressionScanner) setWinAndLoseTokens(setWinFirst bool, firstTokenValues, secondTokenValues []string) {
	sc.ExpressionType = AddMatchCommandType
	if setWinFirst {
		sc.tokens = append(sc.tokens, Token{
			SetWin,
			firstTokenValues,
		})
		sc.tokens = append(sc.tokens, Token{
			SetLose,
			secondTokenValues,
		})
	} else {
		sc.tokens = append(sc.tokens, Token{
			SetWin,
			secondTokenValues,
		})
		sc.tokens = append(sc.tokens, Token{
			SetLose,
			firstTokenValues,
		})
	}
}

func (sc *ExpressionScanner) tryTokenizeSetDraw(sanitazedExpression string) error {
	numOfBracketMatches, tokenBracketValue := extractKeysFromBrackets(sanitazedExpression)
	if numOfBracketMatches == 2 {
		sc.setDrawTokens(tokenBracketValue[0], tokenBracketValue[1])
		return nil
	}

	numOfParanthesisMatches, tokenParanthesisValue := extractKeyFromParanthesis(sanitazedExpression)
	if numOfParanthesisMatches == 2 {
		sc.setDrawTokens(tokenParanthesisValue[:1], tokenParanthesisValue[1:])
		return nil
	}
	return errorInvalidNumberOfArguments(sc.traceId, sanitazedExpression, AddMatchCommandType, 2)
}

func (sc *ExpressionScanner) setDrawTokens(firstTokenValues, secondTokenValues []string) {
	sc.ExpressionType = ExpressionType(AddMatchCommandType)
	sc.tokens = append(sc.tokens, Token{
		SetDraw,
		firstTokenValues,
	})
	sc.tokens = append(sc.tokens, Token{
		SetDraw,
		secondTokenValues,
	})
}

func (sc *ExpressionScanner) tryTokenizeSelectPlayers(sanitazedExpression string) error {
	sc.ExpressionType = SelectQueryType
	sc.tokens = append(sc.tokens, Token{
		SelectPlayers,
		[]string{"*"},
	})

	sc.tryTokenizeOrderBySubExpression(sanitazedExpression)

	return nil
}

func (sc *ExpressionScanner) tryTokenizeOrderBySubExpression(sanitazedExpression string) {
	isMatch, orderSubEpression := findMatch(order_by, sanitazedExpression)
	if !isMatch {
		return
	}

	splittedSubExpression := strings.Split(orderSubEpression, " ")
	orderKey := splittedSubExpression[3:4]

	if strings.ToLower(splittedSubExpression[4]) == "asc" {
		sc.tokens = append(sc.tokens, Token{
			OrderByAsc,
			orderKey,
		})
	} else {
		sc.tokens = append(sc.tokens, Token{
			OrderByDesc,
			orderKey,
		})
	}
}
