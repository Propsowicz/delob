package parser

type regex_pattern string

const (
	add_player         regex_pattern = `(?i)^add player '[^)]+'`
	add_players        regex_pattern = `(?i)^add players \([^)]+\)`
	add_decisive_match regex_pattern = `(?i)^set (win for|lose for) ('[^']*'|\([^)]+\)) and (win for|lose for) ('[^']*'|\([^)]+\))$`
	add_draw_match     regex_pattern = `(?i)^set draw between ('[^']*'|\([^)]+\)) and ('[^']*'|\([^)]+\))$`
	select_statement   regex_pattern = `(?i)^select `
	order_by           regex_pattern = `(?i) order by (elo|key) (asc|desc)`

	// how many of select statemenst do we need? players with stats, players with match, matches. I guess no more. Bcuz of that I will use hardcoded regex to parse select statement
	select_players_with_stats regex_pattern = `(?i) players with stats`
)

// supported functions
// ADD PLAYER 'Tom';
// ADD PLAYERS ('Tom', 'Joe');
// SET WIN FOR 'Tom' AND LOSE FOR 'Joe';
// SET WIN FOR ('Tom', 'Bob') AND LOSE FOR ('Joe', 'Jim');
// SET DRAW BETWEEN 'Joe' AND 'Tom';
// SET DRAW BETWEEN ('Tom', 'Bob') AND ('Joe', 'Jim');
// SELECT Players;
// SELECT Players ORDER BY Elo DESC;
// SELECT Players ORDER BY Key DESC;
// SELECT Players WITH Stats;

const (
	valueInParanthesis regex_pattern = `'([^']+)'`
	valueInBrackets    regex_pattern = `\(\s*'([^']+)'(?:\s*,\s*'([^']+)')*\s*\)`
)
