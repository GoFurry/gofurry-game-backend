package controller

type reviewApi struct{}

var ReviewApi *reviewApi

func init() {
	ReviewApi = &reviewApi{}
}
