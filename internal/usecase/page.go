package usecase

import (
	"fmt"
	"sort"
	"time"

	"github.com/daiki-trnsk/YoiYoi-API/internal/models"
	"github.com/daiki-trnsk/YoiYoi-API/internal/repositories"
	"github.com/daiki-trnsk/YoiYoi-API/pkg/dto"
	"github.com/google/uuid"
)

func GetHomeInfo(userID uuid.UUID) (*dto.HomeResponse, error) {
	user, err := repositories.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	oneWeekAgo := time.Now().AddDate(0, 0, -6).Truncate(24 * time.Hour)
	logs, err := repositories.GetDrinkLogs(userID, &oneWeekAgo, 0)
	if err != nil {
		return nil, err
	}

	// 直近30日間のログ数を取得
	thirtyDaysAgo := time.Now().AddDate(0, 0, -29).Truncate(24 * time.Hour)
	logs30, err := repositories.GetDrinkLogs(userID, &thirtyDaysAgo, 0)
	if err != nil {
		return nil, err
	}
	logsCount30 := len(logs30)

	weeklyStats := dto.WeeklyStats{
		TotalAlcoholMl:   0,
		AlcoholByWeekday: map[string]int{"Mon": 0, "Tue": 0, "Wed": 0, "Thu": 0, "Fri": 0, "Sat": 0, "Sun": 0},
	}
	for _, log := range logs {
		drinks, err := repositories.GetDrinkDetails(log.ID)
		if err != nil {
			return nil, err
		}
		weekday := log.DrinkDate.Weekday().String()[:3]
		for _, drink := range drinks {
			weeklyStats.TotalAlcoholMl += int(drink.AmountMl)
			weeklyStats.AlcoholByWeekday[weekday] += int(drink.AmountMl)
		}
	}

	recentLogs, err := repositories.GetDrinkLogs(userID, nil, 3)
	if err != nil {
		return nil, err
	}
	var recentLogDetails []dto.DrinkLogWithDetails
	for _, log := range recentLogs {
		drinks, err := repositories.GetDrinkDetails(log.ID)
		if err != nil {
			return nil, err
		}
		recentLogDetails = append(recentLogDetails, dto.DrinkLogWithDetails{
			DrinksLogs: log,
			Drinks:     drinks,
		})
	}

	resp := &dto.HomeResponse{
		UserInfo:    dto.ToUserResponse(*user),
		WeeklyStats: weeklyStats,
		RecentLogs:  recentLogDetails,
		LogsCount30: logsCount30,
	}
	return resp, nil
}

func GetTimelineInfo(userID uuid.UUID) (*dto.TimelineResponse, error) {
	friends, err := repositories.GetFriends(userID)
	if err != nil {
		return nil, err
	}
	friendIDs := make(map[uuid.UUID]string) // statusも保持
	for _, f := range friends {
		var friendID uuid.UUID
		if f.FollowerID == userID {
			friendID = f.FolloweeID
		} else {
			friendID = f.FollowerID
		}
		friendIDs[friendID] = f.Status // statusを記録
	}
	ids := make([]uuid.UUID, 0, len(friendIDs))
	for id := range friendIDs {
		ids = append(ids, id)
	}
	users, err := repositories.GetUsersByIDs(ids)
	if err != nil {
		return nil, err
	}

	// friendList: pending→acceptedの順で並べる
	var pendingList, acceptedList []dto.FriendWithStatus
	for _, f := range friends {
		var friendID uuid.UUID
		if f.FollowerID == userID {
			friendID = f.FolloweeID
		} else {
			friendID = f.FollowerID
		}
		// ユーザー情報取得
		u, err := findUserByID(users, friendID)
		if err != nil {
			continue // ユーザー情報がなければスキップ
		}
		friend := dto.FriendWithStatus{
			UserResponse: dto.ToUserResponse(*u),
			Status:       f.Status,
			FriendID:     f.ID,
		}
		if f.Status == "pending" {
			pendingList = append(pendingList, friend)
		} else {
			acceptedList = append(acceptedList, friend)
		}
	}
	friendList := append(pendingList, acceptedList...)

	// timeline: acceptedのみ
	var timeline []dto.Timeline
	for _, u := range users {
		status := friendIDs[u.ID]
		if status != "accepted" {
			continue
		}
		logs, err := repositories.GetDrinkLogs(u.ID, nil, 10)
		if err != nil {
			return nil, err
		}
		for _, log := range logs {
			drinks, err := repositories.GetDrinkDetails(log.ID)
			if err != nil {
				return nil, err
			}
			timeline = append(timeline, dto.Timeline{
				User:     dto.ToUserResponse(u),
				DrinkLog: dto.DrinkLogWithDetails{DrinksLogs: log, Drinks: drinks},
			})
		}
	}
	sort.Slice(timeline, func(i, j int) bool {
		return timeline[i].DrinkLog.DrinksLogs.DrinkDate.After(timeline[j].DrinkLog.DrinksLogs.DrinkDate)
	})

	return &dto.TimelineResponse{
		FriendList: friendList,
		Timeline:   timeline,
	}, nil
}

func GetPeriodStatsInfo(userID uuid.UUID, days int) (*dto.PeriodStatsResponse, error) {
	startDate := time.Now().AddDate(0, 0, -days+1).Truncate(24 * time.Hour)
	logs, err := repositories.GetDrinkLogs(userID, &startDate, 0)
	if err != nil {
		return nil, err
	}
	weekdayAlcohol := map[string]int{"Mon": 0, "Tue": 0, "Wed": 0, "Thu": 0, "Fri": 0, "Sat": 0, "Sun": 0}
	alcoholByDrinkType := make(map[string]int)
	totalAlcoholGram := 0
	dateSet := make(map[string]struct{})
	for _, log := range logs {
		dateStr := log.DrinkDate.Format("2006-01-02")
		dateSet[dateStr] = struct{}{}
		weekday := log.DrinkDate.Weekday().String()[:3]
		drinks, err := repositories.GetDrinkDetails(log.ID)
		if err != nil {
			return nil, err
		}
		for _, drink := range drinks {
			alcoholMl := float64(drink.AmountMl) * drink.Abv / 100.0
			alcoholGram := int(alcoholMl * 0.8)
			totalAlcoholGram += alcoholGram
			weekdayAlcohol[weekday] += alcoholGram
			alcoholByDrinkType[drink.Name] += alcoholGram
		}
	}
	periodDays := days
	actualDays := len(dateSet)
	var averageAlcoholGram float64
	if periodDays > 0 {
		averageAlcoholGram = float64(totalAlcoholGram) / float64(periodDays)
	}
	resp := &dto.PeriodStatsResponse{
		TotalAlcoholGram:   totalAlcoholGram,
		AverageAlcoholGram: averageAlcoholGram,
		PeriodDays:         periodDays,
		ActualDrinkDays:    actualDays,
		AlcoholByWeekday:   weekdayAlcohol,
		AlcoholByDrinkType: alcoholByDrinkType,
	}
	return resp, nil
}

// ユーティリティ関数
func findUserByID(users []models.Users, id uuid.UUID) (*models.Users, error) {
	for _, u := range users {
		if u.ID == id {
			return &u, nil
		}
	}
	return nil, fmt.Errorf("user not found")
}
