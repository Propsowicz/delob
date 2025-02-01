package tokenizer

func selectAllTokenizer(expression string) ([]TokenizedExpression, error) {
	return []TokenizedExpression{
		TokenizedExpression{
			ProcessMethod: SelectAll,
			Arguments:     []string{},
		},
	}, nil
}
