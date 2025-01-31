package processor

const k int16 = 25

type Calculator struct {
	playerWinElo  int16
	playerLoseElo int16
}

func NewCalculator(playerWin Player, playerLose Player) Calculator {
	playerWinElo := playerWin.Elo + k
	playerLoseElo := playerLose.Elo - k

	return Calculator{
		playerWinElo:  playerWinElo,
		playerLoseElo: playerLoseElo,
	}
}

func (c *Calculator) GetWinElo() int16 {
	return c.playerWinElo
}

func (c *Calculator) GetLoseElo() int16 {
	return c.playerLoseElo
}
