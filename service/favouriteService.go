package service

import (
  "github.com/RaymondCode/simple-demo/repository"
  "github.com/RaymondCode/simple-demo/utils"
  "sync"
)

type UserFavouriteResponse struct {
	repository.Response
}
var favoriteLock sync.Mutex 

func FavouriteOrNot(username string, VideoId int, ActionType int) (UserFavouriteResponse) {
  //点赞
	if ActionType == 1 { 
    exits, likedvideos := repository.NewFavouriteInstance().CheckLikedStatus(username, VideoId)
    
    //存在记录
		if exits == true { 

      //liked=1，再点赞favoritecount不增加
      if likedvideos[0].Liked == 1{
        UserFavouriteResponse := UserFavouriteResponse{
				Response: repository.Response{
					StatusCode: 0,
					StatusMsg:  "你已经点赞过了"},
        }
        return UserFavouriteResponse
			}else{ //liked=0，再点赞favoritecount增加
        err := repository.NewFavouriteInstance().ChangeUnfavToFav(username, VideoId)
  			if err != nil {
  				UserFavouriteResponse := UserFavouriteResponse{
  					Response: repository.Response{
  						StatusCode: 1,
  						StatusMsg:  " like failed" + err.Error(),
  					},
  				}
  				return UserFavouriteResponse
  			}
        //更新User表的TotalFavorited(视频作者获赞数量)
        authornames, err := repository.NewVideoDaoInstance().QueryAuthorNameByVideoId(VideoId)
        if err != nil {
          UserFavouriteResponse := UserFavouriteResponse{
            Response: repository.Response{
              StatusCode: 1,
              StatusMsg:  " favorite action query author name failed",
            },
          }
          return UserFavouriteResponse
        }

        favoriteLock.Lock()
        err = repository.NewUserDaoInstance().UpdateUserTotalFavorited(authornames[0], 1)
        favoriteLock.Unlock()
        
        if err != nil {
          UserFavouriteResponse := UserFavouriteResponse{
            Response: repository.Response{
              StatusCode: 1,
              StatusMsg:  " favorite action update totfavorited failed",
            },
          }
          return UserFavouriteResponse
        }       

        //更新Video表的favoritecount（红心下数字）
        favoriteLock.Lock()
        err = repository.NewFavouriteInstance().UpdateFavouriteCountPlus(VideoId)
        favoriteLock.Unlock()
        
        if err != nil {
          UserFavouriteResponse := UserFavouriteResponse{
            Response: repository.Response{
              StatusCode: 1,
              StatusMsg:  " favorite count add failed",
            },
          }
          return UserFavouriteResponse
        }
        UserFavouriteResponse := UserFavouriteResponse{
          Response: repository.Response{
            StatusCode: 0,
            StatusMsg:  "点赞成功"},
        }
        return UserFavouriteResponse
      }
    
		} else { //不存在记录
      favoriteLock.Lock()
			err := repository.NewFavouriteInstance().CreateLiked(username, VideoId)
      favoriteLock.Unlock()
			if err != nil {
				UserFavouriteResponse := UserFavouriteResponse{
					Response: repository.Response{
						StatusCode: 1,
						StatusMsg:  " like failed",
					},
				}
				return UserFavouriteResponse
			}

      //更新User表的TotalFavorited(视频作者获赞数量)
      authornames, err := repository.NewVideoDaoInstance().QueryAuthorNameByVideoId(VideoId)
      if err != nil {
        UserFavouriteResponse := UserFavouriteResponse{
          Response: repository.Response{
            StatusCode: 1,
            StatusMsg:  " favorite action query author name failed",
          },
        }
        return UserFavouriteResponse
      }

      favoriteLock.Lock()
      err = repository.NewUserDaoInstance().UpdateUserTotalFavorited(authornames[0], 1)
      favoriteLock.Unlock()
      
      if err != nil {
        UserFavouriteResponse := UserFavouriteResponse{
          Response: repository.Response{
            StatusCode: 1,
            StatusMsg:  " favorite action update totfavorited failed",
          },
        }
        return UserFavouriteResponse
      }

      //更新Video表的favoritecount（红心下数字）
      favoriteLock.Lock()
      err = repository.NewFavouriteInstance().UpdateFavouriteCountPlus(VideoId)
      favoriteLock.Unlock()
      
      if err != nil {
        UserFavouriteResponse := UserFavouriteResponse{
          Response: repository.Response{
            StatusCode: 1,
            StatusMsg:  " favorite count add failed",
          },
        }
        return UserFavouriteResponse
      } 
      
			UserFavouriteResponse := UserFavouriteResponse{
				Response: repository.Response{
					StatusCode: 0,
					StatusMsg:  "like success" + username},
			}
			return UserFavouriteResponse
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
			return UserFavouriteResponse
		}

    //更新User表的TotalFavorited(视频作者获赞数量)
    authornames, err := repository.NewVideoDaoInstance().QueryAuthorNameByVideoId(VideoId)
    if err != nil {
      UserFavouriteResponse := UserFavouriteResponse{
        Response: repository.Response{
          StatusCode: 1,
          StatusMsg:  " favorite action query author name failed",
        },
      }
      return UserFavouriteResponse
    }

    favoriteLock.Lock()
    err = repository.NewUserDaoInstance().UpdateUserTotalFavorited(authornames[0], 0)
    favoriteLock.Unlock()
    
    if err != nil {
      UserFavouriteResponse := UserFavouriteResponse{
        Response: repository.Response{
          StatusCode: 1,
          StatusMsg:  " favorite action update totfavorited failed",
        },
      }
      return UserFavouriteResponse
    }
    
    
    //更新Video表的favoritecount（红心下数字）
    favoriteLock.Lock()
    err = repository.NewFavouriteInstance().UpdateFavouriteCountMinus(VideoId)
    favoriteLock.Unlock()
    
    if err != nil {
      UserFavouriteResponse := UserFavouriteResponse{
        Response: repository.Response{
          StatusCode: 1,
          StatusMsg:  " favorite count minus failed",
        },
      }
      return UserFavouriteResponse
    } 
    
		UserFavouriteResponse := UserFavouriteResponse{
			Response: repository.Response{
				StatusCode: 0,
				StatusMsg:  "取消点赞成功"},
		}
		return UserFavouriteResponse
	}
}


var favoritelistLock sync.Mutex 
func FavoriteList(username string) (VideoListResponse) {
	videoIDlist, err := repository.NewFavouriteInstance().QueryFavoriteVideoIdbyUsername(username)
	if err != nil {
    favoritelist_res := VideoListResponse{
      Response: repository.Response{
				StatusCode: 1,
				StatusMsg:  "get favorite video id error：" + err.Error()},
    }
		return favoritelist_res
	}

  videolist, err := repository.NewVideoDaoInstance().QueryFavoriteVideoListbyVideoIds(videoIDlist)
  if err != nil {
    favoritelist_res := VideoListResponse{
      Response: repository.Response{
				StatusCode: 1,
				StatusMsg:  "get favorite video error：" + err.Error()},
    }
		return favoritelist_res
	}
  
  for index, _ := range videolist{
      videolist[index].IsFavorite = true
  }

  //更新用户喜欢的作品总数
  favoritelistLock.Lock()
  err = repository.NewUserDaoInstance().UpdateUserFavoriteCount(username, len(videolist))
  favoritelistLock.Unlock()
  
  if err != nil{
    videolist_res := VideoListResponse{
  		Response: repository.Response{
  			StatusCode: 1,
        StatusMsg:  "update user favoritecount falied：" + err.Error()},
	  }
    return videolist_res
  }
  videolist_final, err := utils.ConvertVideoDBToJSON(&videolist)
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
