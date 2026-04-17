package routeTests

import (
	"fmt"
	"inside-athletics/internal/models"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestGetPopularPostsOrdersByEngagementRecencyAndRelevance(t *testing.T) {
	t.Parallel()
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)

	api := testDB.API

	currentUser := models.User{
		ID:                      JohnID,
		FirstName:               "Current",
		LastName:                "User",
		Email:                   "current@example.com",
		Username:                "current-user",
		Verified_Athlete_Status: models.VerifiedAthleteStatusPending,
	}
	if err := testDB.DB.Create(&currentUser).Error; err != nil {
		t.Fatalf("create current user: %v", err)
	}

	popularity := int32(100)
	soccer := models.Sport{ID: SoccerID, Name: "Soccer", Popularity: &popularity}
	if err := testDB.DB.Create(&soccer).Error; err != nil {
		t.Fatalf("create sport: %v", err)
	}

	college := models.College{
		ID:           NortheasternID,
		Name:         "Northeastern University",
		State:        "MA",
		City:         "Boston",
		DivisionRank: models.DivisionI,
		Website:      "https://example.edu",
	}
	if err := testDB.DB.Create(&college).Error; err != nil {
		t.Fatalf("create college: %v", err)
	}

	subscribedTag := models.Tag{Name: "recruiting"}
	otherTag := models.Tag{Name: "general"}
	if err := testDB.DB.Create(&subscribedTag).Error; err != nil {
		t.Fatalf("create subscribed tag: %v", err)
	}
	if err := testDB.DB.Create(&otherTag).Error; err != nil {
		t.Fatalf("create other tag: %v", err)
	}

	relevantPost := models.Post{
		ID:          uuid.New(),
		AuthorID:    currentUser.ID,
		SportID:     &SoccerID,
		CollegeID:   &NortheasternID,
		Title:       "Relevant and active",
		Content:     "This should rank first.",
		IsAnonymous: false,
		CreatedAt:   time.Now().Add(-2 * time.Hour),
		UpdatedAt:   time.Now().Add(-2 * time.Hour),
	}
	likedPost := models.Post{
		ID:          uuid.New(),
		AuthorID:    currentUser.ID,
		Title:       "Lots of likes, less discussion",
		Content:     "This should rank below the discussion-heavy post.",
		IsAnonymous: false,
		CreatedAt:   time.Now().Add(-2 * time.Hour),
		UpdatedAt:   time.Now().Add(-2 * time.Hour),
	}
	oldPost := models.Post{
		ID:          uuid.New(),
		AuthorID:    currentUser.ID,
		Title:       "Old post",
		Content:     "Old engagement should not dominate.",
		IsAnonymous: false,
		CreatedAt:   time.Now().Add(-10 * 24 * time.Hour),
		UpdatedAt:   time.Now().Add(-10 * 24 * time.Hour),
	}
	if err := testDB.DB.Create(&relevantPost).Error; err != nil {
		t.Fatalf("create relevant post: %v", err)
	}
	if err := testDB.DB.Create(&likedPost).Error; err != nil {
		t.Fatalf("create liked post: %v", err)
	}
	if err := testDB.DB.Create(&oldPost).Error; err != nil {
		t.Fatalf("create old post: %v", err)
	}

	if err := testDB.DB.Create(&models.TagPost{ID: uuid.New(), PostableID: relevantPost.ID, PostableType: "post", TagID: subscribedTag.ID}).Error; err != nil {
		t.Fatalf("tag relevant post: %v", err)
	}
	if err := testDB.DB.Create(&models.TagPost{ID: uuid.New(), PostableID: likedPost.ID, PostableType: "post", TagID: otherTag.ID}).Error; err != nil {
		t.Fatalf("tag liked post: %v", err)
	}
	if err := testDB.DB.Create(&models.TagFollow{ID: uuid.New(), UserID: currentUser.ID, TagID: subscribedTag.ID}).Error; err != nil {
		t.Fatalf("create tag follow: %v", err)
	}
	if err := testDB.DB.Create(&models.SportFollow{ID: uuid.New(), UserID: currentUser.ID, SportID: SoccerID}).Error; err != nil {
		t.Fatalf("create sport follow: %v", err)
	}
	if err := testDB.DB.Create(&models.CollegeFollow{ID: uuid.New(), UserID: currentUser.ID, CollegeID: NortheasternID}).Error; err != nil {
		t.Fatalf("create college follow: %v", err)
	}

	for i := 0; i < 2; i++ {
		user := models.User{
			ID:                      uuid.New(),
			FirstName:               fmt.Sprintf("Commenter%d", i),
			LastName:                "User",
			Email:                   fmt.Sprintf("commenter-%d@example.com", i),
			Username:                fmt.Sprintf("commenter-%d", i),
			Verified_Athlete_Status: models.VerifiedAthleteStatusPending,
		}
		if err := testDB.DB.Create(&user).Error; err != nil {
			t.Fatalf("create commenter %d: %v", i, err)
		}

		comment := models.Comment{
			ID:          uuid.New(),
			UserID:      user.ID,
			PostID:      relevantPost.ID,
			Description: "High-velocity comment",
			CreatedAt:   time.Now().Add(-90 * time.Minute),
			UpdatedAt:   time.Now().Add(-90 * time.Minute),
		}
		if err := testDB.DB.Create(&comment).Error; err != nil {
			t.Fatalf("create comment %d: %v", i, err)
		}
	}

	for i := 0; i < 4; i++ {
		user := models.User{
			ID:                      uuid.New(),
			FirstName:               fmt.Sprintf("Liker%d", i),
			LastName:                "User",
			Email:                   fmt.Sprintf("liker-%d@example.com", i),
			Username:                fmt.Sprintf("liker-%d", i),
			Verified_Athlete_Status: models.VerifiedAthleteStatusPending,
		}
		if err := testDB.DB.Create(&user).Error; err != nil {
			t.Fatalf("create liker %d: %v", i, err)
		}

		like := models.PostLike{
			ID:        uuid.New(),
			UserID:    user.ID,
			PostID:    likedPost.ID,
			CreatedAt: time.Now().Add(-90 * time.Minute),
		}
		if err := testDB.DB.Create(&like).Error; err != nil {
			t.Fatalf("create like %d: %v", i, err)
		}
	}

	oldLiker := models.User{
		ID:                      uuid.New(),
		FirstName:               "Old",
		LastName:                "Liker",
		Email:                   "old-liker@example.com",
		Username:                "old-liker",
		Verified_Athlete_Status: models.VerifiedAthleteStatusPending,
	}
	if err := testDB.DB.Create(&oldLiker).Error; err != nil {
		t.Fatalf("create old liker: %v", err)
	}
	if err := testDB.DB.Create(&models.PostLike{
		ID:        uuid.New(),
		UserID:    oldLiker.ID,
		PostID:    oldPost.ID,
		CreatedAt: time.Now().Add(-9 * 24 * time.Hour),
	}).Error; err != nil {
		t.Fatalf("create old like: %v", err)
	}

	authHeader := authHeaderWithPermissionsGivenUser(t, testDB.DB, nil, currentUser.ID)

	resp := api.Get("/api/v1/posts/popular?limit=3&window_hours=72", authHeader)
	if resp.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", resp.Code, resp.Body.String())
	}

	var result struct {
		Posts []struct {
			ID              uuid.UUID `json:"id"`
			Title           string    `json:"title"`
			LikeCount       int64     `json:"like_count"`
			CommentCount    int64     `json:"comment_count"`
			PopularityScore float64   `json:"popularity_score"`
		} `json:"posts"`
		Total int `json:"total"`
	}
	DecodeTo(&result, resp)

	if len(result.Posts) != 3 {
		t.Fatalf("expected 3 posts, got %d", len(result.Posts))
	}
	if result.Posts[0].ID != relevantPost.ID {
		t.Fatalf("expected relevant post first, got %s", result.Posts[0].ID)
	}
	if result.Posts[1].ID != likedPost.ID {
		t.Fatalf("expected liked post second, got %s", result.Posts[1].ID)
	}
	if result.Posts[2].ID != oldPost.ID {
		t.Fatalf("expected old post third, got %s", result.Posts[2].ID)
	}
	if result.Posts[0].CommentCount != 2 {
		t.Fatalf("expected 2 comments on top post, got %d", result.Posts[0].CommentCount)
	}
	if result.Posts[1].LikeCount != 4 {
		t.Fatalf("expected 4 likes on second post, got %d", result.Posts[1].LikeCount)
	}
	if result.Posts[0].PopularityScore <= result.Posts[1].PopularityScore {
		t.Fatalf("expected top post score %.2f to exceed second post score %.2f", result.Posts[0].PopularityScore, result.Posts[1].PopularityScore)
	}
}
