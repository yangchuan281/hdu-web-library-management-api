

# 图书管理系统 - 从零开始看懂项目

> 同学你好！这份文档是专门为你写的，假设你完全不懂 Web 开发。
> 我会用最通俗的语言解释每个文件是干什么的、它们之间怎么配合的。
> 读完这份文档，你就能跟老师讲清楚这个项目了。

---

## 一、先搞懂三个基本概念

### 1.1 什么是Web应用？

你用浏览器打开的网站（比如淘宝、B站）就是 Web 应用。
一个 Web 应用 = 前端 + 后端 + 数据库

形象地说：
你（用户）-> 打开浏览器（前端页面）-> 发请求到服务器（后端程序）-> 查数据库 -> 返回数据 -> 页面展示

### 1.2 什么是前后端分离？

我们这个项目就是前后端分离的：

- 后端（Go语言写的）：只负责处理数据、做业务逻辑，像个后台管理员
- 前端（Vue3写的）：只负责展示界面、接收用户点击，像个前台接待员

前后端通过 API（接口）交流，就像两个人通过微信发消息。

### 1.3 什么是 API？

API 就是前后端约定好的暗号。
每个 API 包含：

- 请求方法：POST（创建）、GET（查询）、PUT（更新）、DELETE（删除）
- 路径：比如 /api/books 表示图书相关
- 数据：请求时带什么参数，返回什么结果

---

## 二、项目总体结构（先看全景图）

library-management-api/
- main.go                后端的开关，启动程序从这里开始
- go.mod / go.sum        Go语言的依赖清单

- api/                   API说明书：定义前端和后端怎么通信
  - book/                关于图书的通信规则
  - borrow/              关于借阅的通信规则
  - user/                关于用户的通信规则

- internal/              后端核心代码（最重要的文件夹）
  - cmd/cmd.go           路由注册（哪个请求交给哪个函数处理）
  - controller/          控制器（接待员：收到请求后转给业务层）
  - service/             业务层（真正的干活的人）
  - dao/                 数据层（跟数据库打交道）
  - model/               数据模型（定义数据长什么样）

- manifest/              配置文件
  - config/config.yaml   配置（数据库密码、端口等）
  - sql/create.sql       数据库建表SQL

- library-frontend/      前端项目（你看到的页面）
  - src/api/             前端发消息的工具
  - src/views/           页面长什么样
  - src/stores/          前端记忆力（记住登录状态）
  - src/router/          页面跳转规则
  - src/layouts/         页面整体布局
  - index.html           前端入口
  - vite.config.js       前端配置


---

## 三、每个文件夹/文件干什么——最详细版

### 3.1 后端入口：main.go

这个文件就是按下开机键。程序一运行，就跳到 internal/cmd/cmd.go 去执行。

代码做了三件事：
1. 导入MySQL驱动（让Go能连接MySQL数据库）
2. 导入cmd包（真正的启动逻辑在里面）
3. 调用cmd.Main.Run(gctx.New()) 启动HTTP服务器

### 3.2 路由注册：internal/cmd/cmd.go（最重要的文件）

这个文件是整个后端的大脑。它做了三件事：

第一件事：启动 HTTP 服务器
- g.Server() 创建一个HTTP服务器
- s.Run() 启动它，默认监听8000端口

第二件事：注册中间件（安检关卡）
- middlewareSvc.Ctx：每个请求都解析用户身份
- ghttp.MiddlewareCORS：允许跨域（前后端端口不同）

第三件事：注册路由（最关键！）

完整的API路由表如下（按分组）：

/api/auth（公开，无需登录）
- POST /verification-codes  发送验证码
- POST /register            邮箱注册
- POST /login               登录
- DELETE /logout            退出登录
- PUT /password             重置密码

/api/users（公开）
- POST /                    手机号注册
- GET /me                   获取个人信息

/api/books（公开查询）
- GET /                     图书列表（分页+搜索）
- GET /{id}                 图书详情

以上需要登录的区域（Auth中间件检查）：
  需要管理员权限（Admin中间件检查）：
  - POST /books             创建图书
  - PUT /books/{id}         更新图书
  - DELETE /books/{id}      删除图书

  登录就能访问：
  - POST /borrows           借书
  - GET /borrows            借阅列表
  - PUT /borrows/{id}/return 还书


### 3.3 三层架构详解（老师必问）

GoFrame框架推荐的分层方式，逻辑非常清晰：

用户请求 -> Controller（控制层，餐厅服务员）
               -> Service（业务层，后厨厨师）
                    -> DAO/Model（数据层，仓库管理员）
                         -> MySQL数据库

#### API 层（api/ 文件夹）

作用：定义请求长什么样、响应长什么样。

以 api/book/v1/book.go 为例：

创建图书时，前端要发这些数据：
- Title（书名，必填）
- Author（作者，必填）
- Isbn（ISBN号，必填）
- PublishDate（出版日期，可选）

创建成功后，后端返回：
- Id（新书的ID）
- Title（书名）
- Isbn（ISBN号）

这一层的作用是规范前后端通信格式。

#### Controller 层（internal/controller/ 文件夹）

作用：接收请求 -> 调用 Service -> 返回结果。
代码很薄，就是个传话筒。

每个Controller文件的结构：
1. ControllerV1 结构体（保存了Service的引用）
2. NewV1() 函数（创建Controller实例）
3. 各种处理方法（Create、List、Get、Update、Delete等）

以图书的创建为例，流程是：
1. 从请求中取出书名、作者、ISBN
2. 调用 bookSvc.Create() 去创建
3. 如果出错，返回错误
4. 如果成功，把结果包装成API定义的格式返回

#### Service 层（internal/service/ 文件夹）——核心逻辑都在这里

图书服务（internal/service/book/book.go）的功能：

1. Create（创建图书）：
   - 先检查ISBN号有没有被用过
   - 如果没有被用过，就插入数据库

2. List（图书列表）：
   - 支持按书名模糊搜索（比如搜三能找到三国演义）
   - 支持按状态筛选（可借阅/已借出）
   - 支持分页（一页显示10本）

3. Get（图书详情）：根据ID查一本书

4. Update（更新图书）：
   - 只更新用户填了的字段（没填的不修改）
   - 如果改了ISBN，要检查新的ISBN有没有被其他书占用

5. Delete（删除图书）：
   - 先检查这本书有没有未归还的借阅记录
   - 如果有，不能删除

用户服务（internal/service/user/user.go）的功能：

1. Create（注册）：检查手机号/邮箱是否已被注册，用MD5加密密码
2. Login（登录）：比对密码，生成JWT令牌
3. SendRegisterCode（发送注册验证码）：生成4位数字验证码，存数据库，发邮件
4. EmailSignUp（邮箱注册）：验证验证码，创建用户
5. ResetPassword（重置密码）：验证重置码，更新密码
6. LoginByEmail（邮箱登录）：支持邮箱+密码登录，返回JWT Token

借阅服务（internal/service/borrow/borrow.go）的功能：

1. Borrow（借书）：
   - 检查图书存在且可借阅
   - 检查用户是否已经借了同一本书还没还
   - 创建借阅记录
   - 把图书状态改为已借出
   - 以上四步在一个事务中完成

2. Return（还书）：
   - 检查借阅记录存在
   - 检查是否已经还过
   - 更新归还时间
   - 把图书状态恢复为可借阅

3. List（借阅列表）：
   - 支持按状态筛选
   - 支持分页

#### DAO 层（internal/dao/ 文件夹）

作用：跟数据库打交道。GoFrame框架自动生成的。

你只需要在Service里这样调用：
dao.Books.Ctx(ctx).Where(id, 1).Scan(&book)    查ID为1的图书
dao.Books.Ctx(ctx).Data(data).InsertAndGetId()  插入一条记录

#### Model 层（internal/model/ 文件夹）

作用：定义数据库表在Go语言里的长相。

Books结构体：
- Id（uint64）：图书ID
- Title（string）：书名
- Author（string）：作者
- Isbn（string）：ISBN号
- PublishDate（time）：出版日期
- Status（int8）：状态（1可借阅/0已借出/2下架）
- CreatedAt（time）：创建时间
- UpdatedAt（time）：更新时间


### 3.4 中间件（internal/service/middleware/middleware.go）

中间件就像是安检关卡，每个请求都要经过它。

三个中间件：

1. Ctx 中间件：每个请求都执行，解析用户信息
   - 先看 Session 里有没有用户信息
   - 再看请求头有没有 JWT Token
   - 解析出用户信息后放入上下文

2. Auth 中间件：检查用户是否已登录
   - 没登录 -> 返回403错误请先登录
   - 已登录 -> 放行

3. Admin 中间件：检查用户是否是管理员
   - 不是管理员 -> 返回403错误权限不足
   - 是管理员 -> 放行

请求经过的顺序：
Ctx（解析身份）-> Auth（检查登录）-> Admin（检查管理员）-> 真正的处理函数

### 3.5 JWT 认证（internal/service/jwt/jwt.go）

JWT（JSON Web Token）是一个加密的身份令牌。

登录过程：
1. 用户输入邮箱+密码
2. 后端校验密码正确
3. 后端生成JWT Token，里面包含用户ID、角色等信息
4. 用HS256加密算法签名（相当于盖防伪章）
5. 返回Token给前端
6. 前端存到浏览器的localStorage里
7. 后续每次请求，前端都把Token放在请求头里
8. 后端每次收到请求，先验证Token是否有效、是否过期

Token 24小时有效。

### 3.6 邮件服务（internal/service/email/email.go）

用网易163的SMTP服务器发邮件。

配置：
- 服务器：smtp.163.com
- 端口：465（SSL加密）
- 账号：18295659278@163.com

可以发送两种邮件：
1. 注册验证码（4位数字，5分钟有效）
2. 重置密码验证码

### 3.7 配置文件（manifest/config/config.yaml）

server:
  address: :8000   后端启动在8000端口

database:
  default:
    link: mysql:root:123456@tcp(127.0.0.1:3306)/library_management
    意思是：用root用户，密码123456，连接本机MySQL，用library_management数据库

### 3.8 数据库4张表

users表（用户）：
- id（主键，自动增长）
- name（姓名）
- email（邮箱，唯一）
- phone（手机号，唯一）
- password（MD5加密后的密码）
- role（角色：user普通用户/admin管理员）
- status（状态：1正常/0禁用）

books表（图书）：
- id（主键）
- title（书名）
- author（作者）
- isbn（ISBN号，唯一）
- status（1可借阅/0已借出/2下架）

borrows表（借阅记录）：
- id（主键）
- user_id（用户ID，外键关联users表）
- book_id（图书ID，外键关联books表）
- borrow_at（借书时间）
- return_at（还书时间，NULL表示未还）

verification_codes表（验证码）：
- email（邮箱）
- code（4位数字验证码）
- type（类型：register注册/reset重置）
- used（是否已使用：0未用/1已用）
- expires_at（过期时间）


---

## 四、前端部分详解（library-frontend/）

### 4.1 前端是什么？

就是你打开浏览器看到的那个漂亮页面。用 Vue3 框架 + Element Plus 组件库写的。

### 4.2 前端文件结构

library-frontend/
- index.html           首页HTML（浏览器最先加载这个）
- vite.config.js       配置（端口3000，API代理到8000）
- package.json         依赖清单

src/
- main.js              前端入口（加载所有组件）
- App.vue              根组件

src/api/               跟后端通信的辅助函数
- request.js           Axios配置（发HTTP请求的工具）
- auth.js              登录/注册/重置密码的API
- book.js              图书增删改查的API
- borrow.js            借书/还书的API

src/views/             页面
- Login.vue            登录/注册/重置密码页面
- Books.vue            图书浏览（所有人可见）
- BookManage.vue       图书管理（仅管理员）
- Borrows.vue          我的借阅
- Profile.vue          个人中心

src/stores/auth.js     全局状态管理（记住登录状态）
src/router/index.js    路由配置（URL对应哪个页面）
src/layouts/AppLayout.vue  页面整体框架（侧边栏+顶部栏）

### 4.3 核心文件解释

request.js（所有API调用的基础工具）：
- 创建Axios实例，所有请求自动加上/api前缀
- 请求拦截器：每次发请求前，从localStorage取出Token，放到请求头里
- 响应拦截器：收到响应后，如果出错就弹窗提示
- 如果收到401状态码（Token过期），自动跳回登录页

auth.js（登录注册相关API）：
- loginApi(data)：调用POST /api/auth/login
- registerApi(data)：调用POST /api/auth/register
- sendCodeApi(email, type)：调用POST /api/auth/verification-codes
- resetPasswordApi(data)：调用PUT /api/auth/password

stores/auth.js（记住用户登录状态）：
- 用户登录后，把Token和用户信息存起来
- 刷新页面时，从localStorage恢复Token
- isLoggedIn 判断是否已登录

router/index.js（页面跳转规则）：
- /login -> 登录页
- /books -> 图书浏览
- /books/manage -> 图书管理（需要登录+管理员）
- /borrows -> 我的借阅（需要登录）
- /profile -> 个人中心（需要登录）
- 路由守卫：每次跳转前检查登录状态和权限


---

## 五、一次完整的登录流程

假设你打开浏览器访问 http://localhost:3000：

第1步：浏览器加载前端页面
- index.html -> main.js -> App.vue -> 路由判断 -> Login.vue
- 屏幕上显示登录页面

第2步：你在登录页面输入邮箱和密码，点击登录按钮

第3步：Login.vue 调用 auth.login({ email: xxx, password: 123456 })

第4步：auth.js 里的 login 函数调用 loginApi(data)
       loginApi 调用 request.post(/auth/login, data)

第5步：request.js 里的 Axios 发 HTTP 请求到：
       http://localhost:3000/api/auth/login
       （Vite配置了代理，自动转发到后端的 localhost:8000）

第6步：后端的 cmd.go 匹配到路由 POST /api/auth/login
       交给 Controller -> userCtrl.Login

第7步：Controller 调用 Service 的 Login 方法

第8步：Service 去数据库查 users 表：email=xxx

第9步：比对密码（MD5加密后比对）

第10步：查到了 -> 生成 JWT Token

第11步：返回结果：{code: 0, data: {id: 1, name: 张三, token: eyJ...}}

第12步：前端收到响应
       - 把Token存到localStorage
       - 把用户信息存到auth store
       - 跳转到/books页面
       - 显示图书列表

---

## 六、一次完整的借书流程

第1步：用户在借阅页面输入图书ID，点击借书

第2步：Borrows.vue 调用 borrowBookApi(bookId)

第3步：Axios发POST请求到 /api/borrows
       请求头自动带上了 Bearer Token

第4步：后端cmd.go匹配路由 POST /api/borrows
       先经过Ctx中间件（解析用户身份）
       再经过Auth中间件（检查是否登录）
       然后交给 borrowCtrl.Borrow

第5步：Controller 调用 Service 的 Borrow 方法

第6步：Service 开启数据库事务

第7步：在事务中：
       - 查books表：图书存在且status=1（可借阅）
       - 查borrows表：用户没有借同一本书未还
       - 插入borrows表：创建借阅记录
       - 更新books表：status改为0（已借出）

第8步：事务提交成功，返回结果

第9步：前端显示借书成功，刷新借阅列表


---

## 七、老师验收时可能问的问题（附参考答案）

### Q1：这个项目用了什么技术？
答：后端用 Go 语言和 GoFrame 框架，数据库用 MySQL，前端用 Vue3 + Element Plus 组件库，用 Vite 做构建工具。

### Q2：为什么选择前后端分离？
答：前后端分离可以让前端和后端独立开发和部署。以后如果要加手机App，前端重写一套，后端API不用改直接用。而且分工明确，前端只负责页面展示，后端只负责数据处理。

### Q3：说说你是怎么处理用户登录的？
答：用户输入邮箱和密码后，后端先查数据库校验密码是否正确。校验通过后，用JWT生成一个加密的令牌返回给前端。前端把令牌存在浏览器的localStorage里，之后的每个请求都会在请求头里带上这个令牌。后端通过中间件每次解析令牌，就知道是谁在访问。

### Q4：借书的时候怎么保证数据一致性？
答：借书的四个操作（检查图书、检查用户是否已借、创建借阅记录、更新图书状态）包装在一个数据库事务里。事务的意思是：要么全部成功，要么全部失败。如果中间某一步出错了，前面的操作都会撤销，不会出现书借出去了但状态没更新的情况。

### Q5：怎么区分管理员和普通用户？
答：数据库的users表有个role字段，普通用户是user，管理员是admin。后端有两个中间件：Auth中间件检查是否登录，Admin中间件检查角色是否为admin。图书的增删改接口挂了两个中间件，普通用户访问会被拦截。

### Q6：数据表之间有什么关系？
答：一共4张表。users和borrows通过user_id关联（一个用户可以有多条借阅记录），books和borrows通过book_id关联（一本书可以被多次借阅）。borrows表里用外键约束保证引用的user_id和book_id必须在users表和books表里存在。

### Q7：密码是怎么存储的？
答：密码不是明文存储的，而是用MD5加密后存储。用户登录时，把用户输入的密码也做MD5加密，然后比对加密后的结果是否一致。这样即使数据库泄露，攻击者也拿不到原始密码。

### Q8：验证码是怎么工作的？
答：用户输入邮箱，点击获取验证码，后端生成4位随机数字，存到verification_codes表里（含5分钟过期时间），然后通过163邮箱的SMTP服务发送邮件。用户注册或重置密码时提交验证码，后端检查验证码是否正确、是否过期、是否已被使用。

### Q9：什么是中间件？你用了哪几个？
答：中间件就是在请求到达真正的处理函数之前，先经过的一道道安检关卡。我用了三个：Ctx中间件解析用户身份，Auth中间件检查登录状态，Admin中间件检查管理员权限。


---

## 八、如何运行项目（给老师演示的步骤）

### 第1步：启动 MySQL
确保 MySQL 服务已启动。
你可以在电脑上找到 MySQL 服务并打开它。

### 第2步：建数据库和表
打开 MySQL 的命令行或者用 Navicat/DBeaver 等工具，
执行 manifest/sql/create.sql 文件里的 SQL 语句。
这会创建 library_management 数据库和4张表。

### 第3步：启动后端
打开终端（命令提示符），输入：
cd C:\Users\橘生淮南\Desktop\library-management-api
go run main.go

看到 server started at :8000 就说明后端启动成功。

### 第4步：启动前端
再打开一个新的终端，输入：
cd C:\Users\橘生淮南\Desktop\library-management-api\library-frontend
npm run dev

看到 Local: http://localhost:3000/ 就说明前端启动成功。

### 第5步：打开浏览器
访问 http://localhost:3000
就可以开始使用了。

### 第6步：体验功能
1. 点击登录页面的邮箱注册，注册一个新账号
2. 登录后进入图书浏览页面
3. 在我的借阅页面输入图书ID借书
4. 如果想体验管理员功能，在数据库把某个用户的role改为admin
5. 刷新页面，就能看到图书管理菜单了

---

## 九、如何跟老师展示这是我写的

老师可能会问你哪些是自己写的，你可以这样说：

我自己写的部分：
1. internal/cmd/cmd.go - 所有API路由是我设计的
2. internal/controller/ - 控制器层是我写的
3. internal/service/ - 业务逻辑是我写的（图书CRUD、借书还书、用户注册登录）
4. internal/service/middleware/middleware.go - 鉴权中间件
5. internal/service/jwt/jwt.go - JWT令牌生成和验证
6. internal/service/email/email.go - 邮件发送功能
7. 前端所有.vue文件、api/、stores/、router/ - 前端页面和交互
8. manifest/sql/create.sql - 数据库表设计

GoFrame框架自动生成的（不用说是自己写的）：
- api/ 文件夹
- internal/dao/ 文件夹
- internal/model/entity/ 文件夹

---

## 十、常见问题解决

### 问题1：前端页面显示空白
原因：后端没有启动
解决：检查后端终端，确保 go run main.go 正在运行

### 问题2：登录提示网络错误
原因：MySQL没有启动
解决：启动MySQL服务

### 问题3：注册时验证码收到了但没看到邮件
原因：邮件服务配置了163邮箱，但可能需要授权码
解决：查看后端终端输出，里面有验证码（代码里打了日志可以直接看到）

### 问题4：图书管理菜单看不到
原因：你的用户不是admin角色
解决：在数据库里把users表的role字段改为admin

### 问题5：npm run dev 报错
原因：没有安装依赖
解决：先执行 npm install，再执行 npm run dev
