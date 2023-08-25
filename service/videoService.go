package service

import "github.com/RaymondCode/simple-demo/repository"

type Video struct {
	Id            int64           `json:"id"`
	Author        repository.User `json:"author"`
	PlayUrl       string          `json:"play_url"`
	CoverUrl      string          `json:"cover_url"`
	FavoriteCount int64           `json:"favorite_count"`
	CommentCount  int64           `json:"comment_count"`
	IsFavorite    bool            `json:"is_favorite"`
	Title         string          `json:"title"`
}

func ConvertVideoDBToJSON(dbVideos *[]repository.Video) ([]Video, error) {
	var jsonVideos []Video
	for _, dbVideo := range *dbVideos {
		user, err := repository.NewUserDaoInstance().QueryUserByName(dbVideo.Name)
		if err != nil {
			return jsonVideos, err
		}
		jsonVideo := Video{
			Id:            dbVideo.Id,
			Author:        user[0],
			PlayUrl:       dbVideo.PlayUrl,
			CoverUrl:      dbVideo.CoverUrl,
			FavoriteCount: dbVideo.FavoriteCount,
			CommentCount:  dbVideo.CommentCount,
			IsFavorite:    dbVideo.IsFavorite,
			Title:         dbVideo.Title,
		}
		jsonVideos = append(jsonVideos, jsonVideo)
	}

	return jsonVideos, nil
}
