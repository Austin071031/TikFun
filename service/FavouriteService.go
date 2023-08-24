package service

import (
  "github.com/RaymondCode/simple-demo/repository"
)

type UserFavouriteResponse struct {
	repository.Response
}

func FavouriteOrNot(username string, VideoId int, ActionType int) (UserFavouriteResponse, error) {
  //点赞
	if ActionType == 1 { 
    exits, likedvideos := repository.NewFavouriteInstance().CheckLikedStatus(username, VideoId)
    
    //存在记录
		if exits == true { 

      //已经有点赞记录，再点赞favoritecount不增加
      if likedvideos[0].Liked == 1{
        UserFavouriteResponse := UserFavouriteResponse{
				Response: repository.Response{
					StatusCode: 0,
					StatusMsg:  "你已经点赞过了"},
        }
        return UserFavouriteResponse, nil
			}else{ //没有点赞记录，再点赞favoritecount增加
        err := repository.NewFavouriteInstance().ChangeUnfavToFav(username, VideoId)
  			if err != nil {
  				UserFavouriteResponse := UserFavouriteResponse{
  					Response: repository.Response{
  						StatusCode: 1,
  						StatusMsg:  " like failed" + err.Error(),
  					},
  				}
  				return UserFavouriteResponse, err
  			}
      
       err = repository.NewFavouriteInstance().UpdateFavouriteCountPlus(VideoId)
       if err != nil {
  				UserFavouriteResponse := UserFavouriteResponse{
  					Response: repository.Response{
  						StatusCode: 1,
  						StatusMsg:  " favorite count add failed",
  					},
  				}
  				return UserFavouriteResponse, err
  			}
  			UserFavouriteResponse := UserFavouriteResponse{
  				Response: repository.Response{
  					StatusCode: 0,
  					StatusMsg:  "点赞成功"},
  			}
  			return UserFavouriteResponse, nil
      }
    
		} else { //不存在记录
      
			err := repository.NewFavouriteInstance().CreateLiked(username, VideoId)
			if err != nil {
				UserFavouriteResponse := UserFavouriteResponse{
					Response: repository.Response{
						StatusCode: 1,
						StatusMsg:  " like failed",
					},
				}
				return UserFavouriteResponse, err
			}

     err = repository.NewFavouriteInstance().UpdateFavouriteCountPlus(VideoId)
     if err != nil {
        UserFavouriteResponse := UserFavouriteResponse{
          Response: repository.Response{
            StatusCode: 1,
            StatusMsg:  " favorite count add failed",
          },
        }
        return UserFavouriteResponse, err
      } 
      
			UserFavouriteResponse := UserFavouriteResponse{
				Response: repository.Response{
					StatusCode: 0,
					StatusMsg:  "like success" + username},
			}
			return UserFavouriteResponse, nil
		}
	} else { //取消点赞
		err := repository.NewFavouriteInstance().ChangeFavToUnfav(username, VideoId)
		if err != nil {
			UserFavouriteResponse := UserFavouriteResponse{
				Response: repository.Response{
					StatusCode: 1,
					StatusMsg:  " unlike failed",
				},
			}
			return UserFavouriteResponse, err
		}

    err = repository.NewFavouriteInstance().UpdateFavouriteCountMinus(VideoId)
   if err != nil {
      UserFavouriteResponse := UserFavouriteResponse{
        Response: repository.Response{
          StatusCode: 1,
          StatusMsg:  " favorite count minus failed",
        },
      }
      return UserFavouriteResponse, err
    } 
    
		UserFavouriteResponse := UserFavouriteResponse{
			Response: repository.Response{
				StatusCode: 0,
				StatusMsg:  "取消点赞成功"},
		}
		return UserFavouriteResponse, nil
	}
}

func FavoriteList(username string) (VideoListResponse) {
	videolist, err := repository.NewFavouriteInstance().GetFavoriteList(username)
	if err != nil {
    favoritelist_res := VideoListResponse{
      Response: repository.Response{
				StatusCode: 1,
				StatusMsg:  "获取点赞列表失败：" + err.Error()},
    }
		return favoritelist_res
	}

  for _, video := range videolist{
    video.IsFavorite = true
  }

  videolist_final, err := ConvertVideoDBToJSON(&videolist)
  if err != nil{
    videolist_res := VideoListResponse{
  		Response: repository.Response{
  			StatusCode: 1,
        StatusMsg:  "liked videos append user falied：" + err.Error()},
	  }
    return videolist_res
  }
  videolist_res := VideoListResponse{
		Response: repository.Response{
			StatusCode: 0,
		},
		VideoList: videolist_final,
	}
	return videolist_res
}
