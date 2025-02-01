package processor

const k int16 = 25

type Calculator struct {
	playerWinEloChange  int16
	playerLoseEloChange int16
}

func NewCalculator(playerWin Player, playerLose Player) Calculator {
	return Calculator{
		playerWinEloChange:  k,
		playerLoseEloChange: k,
	}
}

func (c *Calculator) GetWinElo() int16 {
	return c.playerWinEloChange
}

func (c *Calculator) GetLoseElo() int16 {
	return c.playerLoseEloChange
}
