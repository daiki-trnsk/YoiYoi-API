package handlers

import (
	"net/http"
	"sort"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/daiki-trnsk/YoiYoi-API/internal/models"
	"github.com/daiki-trnsk/YoiYoi-API/pkg/database"
	"github.com/daiki-trnsk/YoiYoi-API/pkg/dto"
)

func GetHome(c echo.Context) error {
	userIDRaw := c.Get("user_id")
	userID, ok := userIDRaw.(uuid.UUID)
	if !ok {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "user id not found"})
	}

	// ユーザー情報取得
	var user models.Users
	if err := database.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "user not found"})
	}

	// 直近7日分のログ取得
	oneWeekAgo := time.Now().AddDate(0, 0, -6).Truncate(24 * time.Hour)
	var logs []models.DrinksLogs
	if err := database.DB.
		Where("user_id = ? AND drink_date >= ?", userID, oneWeekAgo).
		Order("drink_date desc").
		Find(&logs).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to fetch logs"})
	}

	// 統計計算
	weeklyStats := dto.WeeklyStats{
		TotalAlcoholMl:      0,
		TotalNumberOfDrinks: 0,
		AlcoholByWeekday:    map[string]int{"Mon": 0, "Tue": 0, "Wed": 0, "Thu": 0, "Fri": 0, "Sat": 0, "Sun": 0},
	}
	for _, log := range logs {
		var drinks []models.DrinksDetails
		if err := database.DB.Where("drink_log_id = ?", log.ID).Find(&drinks).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to fetch drink details"})
		}
		weekday := log.DrinkDate.Weekday().String() // "Monday" など
		weekdayShort := weekday[:3]
		for _, drink := range drinks {
			alcoholMl := int(float64(drink.AmountMl) * drink.Abv / 100.0)
			weeklyStats.TotalAlcoholMl += alcoholMl
			weeklyStats.TotalNumberOfDrinks++
			weeklyStats.AlcoholByWeekday[weekdayShort] += alcoholMl
		}
	}

	// 最新3件のログ取得
	var recentLogs []models.DrinksLogs
	if err := database.DB.
		Where("user_id = ?", userID).
		Order("drink_date desc").
		Limit(3).
		Find(&recentLogs).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to fetch recent logs"})
	}

	var recentLogDetails []dto.DrinkLogWithDetails
	for _, log := range recentLogs {
		var drinks []models.DrinksDetails
		if err := database.DB.Where("drink_log_id = ?", log.ID).Find(&drinks).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to fetch drink details"})
		}
		recentLogDetails = append(recentLogDetails, dto.DrinkLogWithDetails{
			DrinksLogs: log,
			Drinks:     drinks,
		})
	}

	resp := dto.HomeResponse{
		UserInfo:    dto.ToUserResponse(user),
		WeeklyStats: weeklyStats,
		RecentLogs:  recentLogDetails,
	}
	return c.JSON(http.StatusOK, resp)
}

func GetTimeline(c echo.Context) error {
	userIDRaw := c.Get("user_id")
	userID, ok := userIDRaw.(uuid.UUID)
	if !ok {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "user id not found"})
	}

	// 1. acceptedなフレンドレコード取得
	var friends []models.Friends
	if err := database.DB.
		Where("(follower_id = ? OR followee_id = ?) AND status = ?", userID, userID, "accepted").
		Find(&friends).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to fetch friends"})
	}

	// 2. 友達のユーザーID一覧を作成
	friendIDs := make(map[uuid.UUID]struct{})
	for _, f := range friends {
		if f.FollowerID == userID {
			friendIDs[f.FolloweeID] = struct{}{}
		} else {
			friendIDs[f.FollowerID] = struct{}{}
		}
	}

	// 3. 友達ユーザー情報取得
	var users []models.Users
	userMap := make(map[uuid.UUID]models.Users)
	if len(friendIDs) > 0 {
		ids := make([]uuid.UUID, 0, len(friendIDs))
		for id := range friendIDs {
			ids = append(ids, id)
		}
		if err := database.DB.Where("id IN ?", ids).Find(&users).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to fetch users"})
		}
		for _, u := range users {
			userMap[u.ID] = u
		}
	}

	// 4. 各友達の全ログをフラットに集める
	var timeline []dto.Timeline
	for _, u := range users {
		var logs []models.DrinksLogs
		if err := database.DB.
			Where("user_id = ?", u.ID).
			Order("drink_date desc").
			Limit(10).
			Find(&logs).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to fetch logs"})
		}
		for _, log := range logs {
			var drinks []models.DrinksDetails
			if err := database.DB.Where("drink_log_id = ?", log.ID).Find(&drinks).Error; err != nil {
				return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to fetch drink details"})
			}
			timeline = append(timeline, dto.Timeline{
				User: dto.ToUserResponse(u),
				DrinkLog: dto.DrinkLogWithDetails{
					DrinksLogs: log,
					Drinks:     drinks,
				},
			})
		}
	}

	// 投稿日順で降順ソート
	sort.Slice(timeline, func(i, j int) bool {
		return timeline[i].DrinkLog.DrinksLogs.DrinkDate.After(timeline[j].DrinkLog.DrinksLogs.DrinkDate)
	})

	// 5. 友達リスト
	friendList := make([]dto.UserResponse, 0, len(users))
	for _, u := range users {
		friendList = append(friendList, dto.ToUserResponse(u))
	}

	resp := dto.TimelineResponse{
		FriendList: friendList,
		Timeline:   timeline,
	}
	return c.JSON(http.StatusOK, resp)
}

func GetWeeklyStats(c echo.Context) error {
	return getPeriodStats(c, 7)
}

func GetMonthlyStats(c echo.Context) error {
	return getPeriodStats(c, 30)
}

// 共通ロジック
func getPeriodStats(c echo.Context, days int) error {
	userIDRaw := c.Get("user_id")
	userID, ok := userIDRaw.(uuid.UUID)
	if !ok {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "user id not found"})
	}

	startDate := time.Now().AddDate(0, 0, -days+1).Truncate(24 * time.Hour)
	var logs []models.DrinksLogs
	if err := database.DB.
		Where("user_id = ? AND drink_date >= ?", userID, startDate).
		Order("drink_date desc").
		Find(&logs).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to fetch logs"})
	}

	// 曜日ごとの合計
	weekdayAlcohol := map[string]int{"Mon": 0, "Tue": 0, "Wed": 0, "Thu": 0, "Fri": 0, "Sat": 0, "Sun": 0}
	totalAlcoholGram := 0
	dateSet := make(map[string]struct{})

	for _, log := range logs {
		dateStr := log.DrinkDate.Format("2006-01-02")
		dateSet[dateStr] = struct{}{}
		weekday := log.DrinkDate.Weekday().String() // "Monday" など
		weekdayShort := weekday[:3]
		var drinks []models.DrinksDetails
		if err := database.DB.Where("drink_log_id = ?", log.ID).Find(&drinks).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to fetch drink details"})
		}
		for _, drink := range drinks {
			alcoholMl := float64(drink.AmountMl) * drink.Abv / 100.0
			alcoholGram := int(alcoholMl * 0.8)
			totalAlcoholGram += alcoholGram
			weekdayAlcohol[weekdayShort] += alcoholGram
		}
	}

	periodDays := days
	// 実際に飲酒記録があった日数
	actualDays := len(dateSet)
	// 平均（記録があった日数で割る場合はactualDays、常に期間日数で割る場合はperiodDays）
	var averageAlcoholGram float64
	if periodDays > 0 {
		averageAlcoholGram = float64(totalAlcoholGram) / float64(periodDays)
	}

	resp := dto.PeriodStatsResponse{
		TotalAlcoholGram:   totalAlcoholGram,
		AverageAlcoholGram: averageAlcoholGram,
		PeriodDays:         periodDays,
		ActualDrinkDays:    actualDays,
		AlcoholByWeekday:   weekdayAlcohol,
	}
	return c.JSON(http.StatusOK, resp)
}
