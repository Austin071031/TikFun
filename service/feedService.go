package service

import (
	"github.com/RaymondCode/simple-demo/repository"
	"github.com/RaymondCode/simple-demo/utils"
	"time"
)

type FeedResponse struct {
	repository.Response
	VideoList []utils.Video `json:"video_list,omitempty"`
	NextTime  int64   `json:"next_time,omitempty"`
}

func Feed(username string) (FeedResponse) {

	var VideoList []utils.Video
	err := repository.NewVideoDaoInstance().UpdateVideoUrl(serverDomain)
	if err != nil {
		feedResponse := FeedResponse{
			Response: repository.Response{
				StatusCode: 1,
				StatusMsg:  "fail to update url!\n"},
		}
		return feedResponse
	}
	videos, err := repository.NewVideoDaoInstance().QueryVideoFeed()
	if err != nil {
		feedResponse := FeedResponse{
			Response: repository.Response{
				StatusCode: 1,
				StatusMsg:  "Fail to query videos!\n"},
		}
		return feedResponse
	}

    //已登陆用户更新视频的点赞状态
  if(username != ""){
    likedVideosId, err := repository.NewFavouriteInstance().QueryFavoriteVideoIdbyUsername(username)
    if err != nil{
      feedResponse := FeedResponse{
			Response: repository.Response{
				StatusCode: 1,
				StatusMsg:  "Fail to query liked videos!\n" + err.Error()},
		  }
		return feedResponse
    }
    
    likedVideosId_map := utils.Tomap(likedVideosId)
    
    for index, video := range *videos{
      _, isliked := likedVideosId_map[int(video.Id)]
      if isliked == true{ 
        (*videos)[index].IsFavorite = true
      }
    }
  }

	VideoList, err = utils.ConvertVideoDBToJSON(videos)
	if err != nil {
		feedResponse := FeedResponse{
			Response: repository.Response{StatusCode: 1,
				StatusMsg: "Videos are found, but Convert is failed!\n"},
		}
		return feedResponse
	}
	feedResponse := FeedResponse{
		Response:  repository.Response{StatusCode: 0},
		VideoList: VideoList,
		NextTime:  time.Now().Unix(),
	}
	return feedResponse
}
