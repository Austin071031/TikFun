# simple-demo

一句话概括项目核心信息-本项目由Go语言完成开发，在用户流量较少的前提下的一个简易版抖音服务端,具有用户、视频等基本功能，点赞、评论等互动功能。

项目服务地址-https://1024code.com/codecubes/vo1jikt

项目文档地址-https://vqi5rytxpuq.feishu.cn/docx/SQRzdrvoLolFyIxOnIucrNXMnGe?from=from_copylink

### 接口功能说明

#### 用户模块实现
##### 用户注册操作（/douyin/user/register/）
    1. 通过username在数据库User表查询注册用户是否存在。
    2. 若注册用户不存在，则在用户表插入新的用户信息来注册用户。此外，用户的Token也将通过JWT编码用户的username来生成。因为用户登陆后的一系列操作（如视频点赞，评论）是需要依赖username来更新用户表的信息的，所以我们选择编码username来生成Token。
    3. 生成的Token和用户的密码会一起放到Loginuserdata表里面进行存储。用户的密码在放入数据库之前还会通过Golang的bcrypt加密库来生成用户密码的密码哈希来存放在数据库中，以防在数据库中的敏感信息的明文存储。最后返回用户注册成功的响应，包括注册用户的ID和Token，页面跳转至用户主页。
    4. 若注册用户存在，则返回注册用户已存在的响应。
    
##### 用户登录操作（/douyin/user/login/）
    1. 通过username在数据库Loginuserdata表查询登陆用户是否存在
    2. 若登陆用户信息在数据库Loginuserdata表不存在，则返回登陆用户不存在的响应。
    3. 若登陆用户信息在数据库Loginuserdata表存在，定义一个PasswordIsTrue的布尔变量，将用户登陆输入的密码与存储在数据库中的用户密码做比对。若密码一致，把PasswordIsTrue设为true，否则设为false。
    4. 若密码正确，则登陆成功，每次登陆成功会通过JWT编码用户的username来更新用户的Token，也将新Token更新到数据库的Loginuserdata表。最后通过username在数据库User表查询到用户ID，返回登陆成功响应，包括用户Token和用户ID。
    5. 若密码不正确，则返回密码不正确登陆失败响应。
    
##### 获取用户主页（/douyin/user）
    1. 通过JWT中间件解码用户Token来验证用户登陆的合法性。通过解码用户Token得到的username来查询数据库User表是否存在该用户信息。若数据库中无该用户信息则返回用户不存在响应。若存在该用户信息则在请求的上下文中将username传进来，以至于可以将username传递给后续的处理函数。
    2. 通过请求的上下文获取到用户的username，并以此在数据库User表中查询用户信息，并返回包含用户信息的响应。

#### 视频模块实现
##### 视频上传操作（/douyin/public/action/）
    1. 从POST请求中获取用户信息（username）和视频信息（title，data）
    2. 若当前用户不存在，返回statuscode为1，结束视频上传；若用户存在，为了防止文件名冲突，采用如下的命名方式：查询数据库获得该视频发布者的用户ID及该视频在数据库中即将对应的ID，使用ffmpeg获得视频第一帧作为视频方面并保存在本地，将视频标题一起作为视频文件和视频封面文件在本地的文件名保存在/public目录下。
    3. 获取当前服务器域名并拼接视频和封面静态数据URL，将视频相关信息保存数据库字段对应的结构体中,并以此添加到video数据库中；
    4. 投稿完成后返回响应。
    - Controlle层——publish.go
      - Publish函数：视频发布的处理
    - Service层——video.go
      - VideoListResponse结构体：表示视频切片及相关信息的响应结构体
         type VideoListResponse struct {
                     repository.Response
                     VideoList          []Video  json:"video_list"
                     Token             string   json:"token"
            }
      - GetServerDomain函数：获取服务器的域名
      - Publish函数：发布并保存上传文件到公共目录
    - repository层——video.go
      - TableName函数：Video 结构体的一个方法，它返回字符串值 "videos"。这个方法用于指定 Video 结构体在数据库中对应的表名
      - VideoDao结构体：封装视频数据的访问和操作逻辑（注：无具体的属性）
              type VideoDao struct {
              }
      - NewVideoDaoInstance函数：获取 VideoDao 的单例实例，确保只有一个实例存在
      - QueryVideoByAuthor函数：按照作者的用户名查询视频列表
      - QueryVideoFeed函数：按照 Feed 的要求返回视频列表
      - QueryVideoLatest函数：查询下一个主键 ID
      - CreateVideo函数：在数据库中创建视频记录
      - UpdateVideoUrl函数：更新视频的 URL
      - UpdateUserFavcount函数：更新用户的点赞总数和作品总数
##### 获取视频发布列表（/douyin/public/list/）
    1. 由于每次重启代码空间，域名会改变，所以需要根据当前域名修改数据库中的URL；
    2. 获取当前页面的用户，并在video数据库中查询该用户发布的所有视频信息，将其信息保存在videos列表中；
    3. videos中保存的数据是与数据库字段对应的信息，将其转换为前端需要的json列表；
    4. 返回响应和json数据
    - Controlle层——publish.go
      - PublishList函数：获取发布视频列表
    - Service层——videoService.go
      - Video结构体：表示视频对象的数据结构体
                type Video struct {
                       Id            int64           json:"id"
                       Author        repository.User  json:"author"
                       PlayUrl       string          json:"play_url"
                       CoverUrl      string          json:"cover_url"
                       FavoriteCount  int64           json:"favorite_count"
                       CommentCount int64           json:"comment_count"
                       IsFavorite      bool           json:"is_favorite"
                       Title          string          json:"title"
                }
      - ConvertVideoDBToJSON函数：将数据库中的视频对象转换为 JSON 格式返回 
    - repository层——video.go
      - PublishList函数：获取发布视频的列表，并对用户的令牌进行验证和处理
##### 用户首页视频推送（/douyin/feed）
    1. 获取当前登录用户的token，若用户信息不存在，则返回错误信息；
    2. 若用户存在或用户未登录，调用服务层的函数完成视频流查询；
    3. 服务层的Feed函数先根据当前服务器域名更新数据库中的URL，然后按照时间顺序倒序查询video数据库中的视频，并返回一个结构体切片，该切片表中包含视频的基本信息；
    4. 然后查询当前用户的点赞信息，并返回当前用户点赞的视频ID数组，若之前返回的包含视频信息的切片中与当前用户点赞视频存在交集，则将视频信息的is_favourte改为true，并返回视频切片和响应。
    - Controlle层——feed.go
      - Feed函数：查询当前登录用户，返回视频流信息
    - Service层——feed.go
      - Feed函数：从数据库中查询视频流信息和当前用户点赞信息，将其合并返回给APP；
    - Repository层——video.go
      - QueryFavoriteVideoIdbyUsername：获取当前用户的点赞视频ID

#### 点赞模块实现
##### 点赞操作（/douyin/favorite/action/）
    1. 通过JWT中间件在用户请求时获取用户姓名
    2. 获取请求中的video_id、用户的ActionType
    3. 用户点赞(ActionType=1)，查询数据库是否有这个用户对这条视频的操作记录。如有，则将用户的点赞状态设置为点赞状态。如果没有则新建记录。并且在Videos表里对该视频ID对应的视频的“点赞红心”下的数字加一；将User表里的该视频作者的获赞数加一。
    4. 用户取消点赞(ActionType=2)，将数据库中用户对这条视频的操作记录设置为取消点赞状态。并且在Videos表里对该视频ID对应的视频的“点赞红心”下的数字减一；将User表里的该视频作者的获赞数减一。
    5. 返回响应
##### 获取点赞列表（/douyin/favorite/list/）
    1. 通过JWT中间件在用户请求时获取用户姓名
    2. 根据用户姓名在Like表中查询该用户点赞过的视频ID
    3. 根据上一步查询到的视频ID在Videos表中查询所对应的视频并将最后查询到的视频列表返回。
    4. 更新用户喜欢的作品总数

#### 评论模块实现
##### 评论操作（/douyin/comment/action/）
    1. 通过JWT中间件在用户请求时获取用户姓名
    2. 获取请求中用户的action_type
    3. 新建评论(ActionType=1)，获取video_id和评论的内容comment_text，并且利用当前登录的用户信息、video_id、comment_text、评论创建时间来新建记录，然后在comments表里进行插入。再在videos表里对video_id对应的视频的commentCount进行加一操作。
    4. 删除评论(ActionType=2)，获取video_id和评论的comment_id，验证当前登入的用户身份，如果当前用户是评论的发布者，则允许该用户对comment_id对应的评论进行删除，然后再comments表里进行删除。再在videos表里对video_id对应的视频的commentCount进行减一操作。
    5. 返回响应
##### 获取评论列表（/douyin/comment/list/）
    1. 获取请求中视频的video_id
    2. 根据video_id在comments表中查询该视频所拥有的所有评论，并将查询到的评论列表返回。

### 测试

进入项目服务地址，运行项目
`go build`
`go run main.go`

然后复制服务端地址到抖声app后台即可测试相关功能。
