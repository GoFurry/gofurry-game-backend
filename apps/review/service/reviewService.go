package service

type reviewService struct{}

var reviewSingleton = new(reviewService)

func GetReviewService() *reviewService { return reviewSingleton }
