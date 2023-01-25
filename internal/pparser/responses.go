package pparser

// Response from /cabinet/get/
type PolemicaGameHistoryResponseRow struct {
	Id        string  `json:"id"`
	DateStart string  `json:"date_start"`
	DateEnds  string  `json:"date_ends"`
	Duration  string  `json:"duration"`
	GameMode  string  `json:"game_mode"`
	Points    float64 `json:"points"`
	Sp        struct {
		Gained int `json:"gained"`
		Total  int `json:"total"`
	} `json:"sp"`
	Role struct {
		Type  string `json:"type"`
		Title string `json:"title"`
	} `json:"role"`
	Result struct {
		Title string `json:"title"`
		Code  string `json:"code"`
	} `json:"result"`
}

type PolemicaGameHistoryResponse struct {
	Rows       []PolemicaGameHistoryResponseRow `json:"rows"`
	TotalCount string                           `json:"totalCount"`
}

// Response from /game-statistics/{gameID}
type GameStatisticsAchievementsSum struct {
	Points       float32 `json:"points"`
	Achievements struct {
		Voting struct {
			Title  string  `json:"title"`
			Points float64 `json:"points"`
		} `json:"voting,omitempty"`
		InconsistentShooting struct {
			Title  string  `json:"title"`
			Points float64 `json:"points"`
		} `json:"inconsistent_shooting,omitempty"`
		Victory struct {
			Title  string `json:"title"`
			Points int    `json:"points"`
		} `json:"victory,omitempty"`
		BestMove struct {
			Title  string  `json:"title"`
			Points float64 `json:"points"`
		} `json:"best_move,omitempty"`
		KillingSheriff struct {
			Title  string  `json:"title"`
			Points float64 `json:"points"`
		} `json:"killing_sheriff,omitempty"`
	} `json:"achievements"`
}

type GameStatisticsPlayerResponse struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	Image    string `json:"image"`
	Role     struct {
		Type  string `json:"type"`
		Title string `json:"title"`
	} `json:"role"`
	TablePosition int    `json:"tablePosition"`
	WL            string `json:"w_l"`
	Points        int    `json:"points"`
	Coins         int    `json:"coins"`
	Achievements  []struct {
		Sum   float64 `json:"sum"`
		Array []struct {
			Title  string  `json:"title"`
			Points float64 `json:"points"`
		} `json:"array"`
	} `json:"achievements"`
	AchievementsSum GameStatisticsAchievementsSum `json:"achievementsSum"`
}

type GameStatisticsResponse struct {
	Id         string  `json:"id"`
	Type       string  `json:"type"`
	DaysNumber float64 `json:"daysNumber"`
	WinnerCode int     `json:"winnerCode"`
	Judge      struct {
		Id       string `json:"id"`
		Username string `json:"username"`
		Coins    int    `json:"coins"`
	} `json:"judge"`
	Players     []GameStatisticsPlayerResponse `json:"players"`
	FirstKilled string                         `json:"firstKilled"`
}
