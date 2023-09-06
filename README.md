# simple-demo

一句话概括项目核心信息-本项目由Go语言完成开发，在用户流量较少的前提下的一个简易版抖音服务端,具有用户、视频等基本功能，点赞、评论等互动功能。

项目服务地址-https://1024code.com/codecubes/vo1jikt

项目文档地址-https://vqi5rytxpuq.feishu.cn/docx/SQRzdrvoLolFyIxOnIucrNXMnGe?from=from_copylink

### 接口功能说明

#### 用户注册操作（/douyin/user/register/）
    1. 通过username在数据库User表查询注册用户是否存在。
    2. 若注册用户不存在，则在用户表插入新的用户信息来注册用户。此外，用户的Token也将通过JWT编码用户的username来生成。因为用户登陆后的一系列操作（如视频点赞，评论）是需要依赖username来更新用户表的信息的，所以我们选择编码username来生成Token。
    3. 生成的Token和用户的密码会一起放到Loginuserdata表里面进行存储。用户的密码在放入数据库之前还会通过Golang的bcrypt加密库来生成用户密码的密码哈希来存放在数据库中，以防在数据库中的敏感信息的明文存储。最后返回用户注册成功的响应，包括注册用户的ID和Token，页面跳转至用户主页。
    4. 若注册用户存在，则返回注册用户已存在的响应。
    
#### 用户登录操作（/douyin/user/login/）
    1. 通过username在数据库Loginuserdata表查询登陆用户是否存在
    2. 若登陆用户信息在数据库Loginuserdata表不存在，则返回登陆用户不存在的响应。
    3. 若登陆用户信息在数据库Loginuserdata表存在，定义一个PasswordIsTrue的布尔变量，将用户登陆输入的密码与存储在数据库中的用户密码做比对。若密码一致，把PasswordIsTrue设为true，否则设为false。
    4. 若密码正确，则登陆成功，每次登陆成功会通过JWT编码用户的username来更新用户的Token，也将新Token更新到数据库的Loginuserdata表。最后通过username在数据库User表查询到用户ID，返回登陆成功响应，包括用户Token和用户ID。
    5. 若密码不正确，则返回密码不正确登陆失败响应。
    
#### 获取用户主页（/douyin/user）
    1. 通过JWT中间件解码用户Token来验证用户登陆的合法性。通过解码用户Token得到的username来查询数据库User表是否存在该用户信息。若数据库中无该用户信息则返回用户不存在响应。若存在该用户信息则在请求的上下文中将username传进来，以至于可以将username传递给后续的处理函数。
    2. 通过请求的上下文获取到用户的username，并以此在数据库User表中查询用户信息，并返回包含用户信息的响应。

### 测试

进入项目服务地址，运行项目
`go build`
`go run main.go`

然后复制服务端地址到抖声app后台即可测试相关功能。
