package controller

import (
	"DEMO01/handler"
	"DEMO01/tools"
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func HandleLogin(c *gin.Context) {
	//登录校验
	db := handler.DB
	var user User_info
	err := c.ShouldBindJSON(&user)
	tools.CheckErr(err)
	//密码进行hash
	code_hash := tools.CodeHash(user.Code)
	var password, user_id string
	//数据库查询用户名和密码
	err = db.QueryRow("select password,user_id from user_info where user_name=?", user.Mobile).Scan(&password, &user_id)
	if err != nil {
		c.JSON(400, gin.H{
			"data":    "",
			"message": "用户名错误",
		})
		return
	}
	//校验正确返回token
	if password == code_hash {

		token := handler.GenerateToken(&handler.UserClaims{
			ID:             user_id,
			Name:           user.Mobile,
			Password:       password,
			StandardClaims: jwt.StandardClaims{},
		})

		data := make(map[string]string)
		data["token"] = token
		c.JSON(200, gin.H{
			"data":    data,
			"message": "OK",
		})
	} else {
		c.JSON(400, gin.H{
			"data":    "",
			"message": "密码错误",
		})
	}
}

// 解决跨域问题
func Cors(context *gin.Context) {
	method := context.Request.Method
	// 必须，接受指定域的请求，可以使用*不加以限制，但不安全
	context.Header("Access-Control-Allow-Origin", "*")
	// context.Header("Access-Control-Allow-Origin", context.GetHeader("Origin"))
	fmt.Println(context.GetHeader("Origin"))
	// 必须，设置服务器支持的所有跨域请求的方法
	context.Header("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE, OPTIONS")
	// 服务器支持的所有头信息字段，不限于浏览器在"预检"中请求的字段
	context.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Token")
	// 可选，设置XMLHttpRequest的响应对象能拿到的额外字段
	context.Header("Access-Control-Expose-Headers", "Access-Control-Allow-Headers, Token")
	// 可选，是否允许后续请求携带认证信息Cookir，该值只能是true，不需要则不设置
	context.Header("Access-Control-Allow-Credentials", "true")
	// 放行所有OPTIONS方法
	if method == "OPTIONS" {
		context.AbortWithStatus(http.StatusNoContent)
		return
	}
	context.Next()
}

// 新用户注册
func HandleRegister(c *gin.Context) {
	//绑定传入json
	var resdata Register_info
	err := c.ShouldBindJSON(&resdata)
	tools.CheckErr(err)
	//连接redis对邮箱验证码进行校验
	rdb := handler.Rdb
	ctx := context.Background()
	num_key := resdata.Mobile + resdata.Email
	num, err := rdb.Get(ctx, num_key).Result()
	if num != resdata.Captcha || err != nil {
		c.JSON(400, gin.H{
			"message": "验证码错误",
		})
		return
	} else {
		c.JSON(200, gin.H{
			"message": "注册成功",
		})
	}
	// 将用户注册信息存入数据库
	db := handler.DB
	password_hash := tools.CodeHash(resdata.Password)
	smtp, err := db.Prepare("insert into user_info (user_name, password, user_email, submission_date)  VALUES (?,?, ?, ?)")
	tools.CheckErr(err)
	_, err = smtp.Exec(resdata.Mobile, password_hash, resdata.Email, time.Now().Format("2006-01-02"))
	tools.CheckErr(err)

}

//发送邮箱验证码

func HandleRegisterCaptcha(c *gin.Context) {
	var registercaptcha Register_captcha
	err := c.ShouldBindJSON(&registercaptcha)
	tools.CheckErr(err)
	//验证用户名是否已经存在
	db := handler.DB
	var name string
	_ = db.QueryRow("select user_name from user_info where user_name=?", registercaptcha.Mobile).Scan(&name)
	if name != "" {
		c.JSON(400, gin.H{
			"message": "用户已存在",
		})
		return
	}
	//验证redis中是否有注册邮箱，没有发送邮箱验证码
	rdb := handler.Rdb
	ctx := context.Background()
	num_key := registercaptcha.Mobile + registercaptcha.Email
	val, _ := rdb.Get(ctx, num_key).Result()
	if val != "" {
		return
	}
	num, err := tools.SendEmailValidate([]string{registercaptcha.Email})
	tools.CheckErr(err)
	rdb.SetNX(ctx, num_key, num, 5*time.Minute)
}

// 返回用户个人信息
func HandleUserinfo(c *gin.Context) {
	user_info, _ := c.Get("user")
	claims := user_info.(*handler.UserClaims)
	c.JSON(200, gin.H{
		"name": claims.Name,
	})
}

// 返回文章分类
func HandleChannels(c *gin.Context) {
	res := `{"data":{"channels":[{"id":0,"name":"推荐"},{"id":1,"name":"html"},{"id":2,"name":"开发者资讯"},{"id":4,"name":"c++"},{"id":6,"name":"css"},{"id":7,"name":"数据库"},{"id":8,"name":"区块链"},{"id":9,"name":"go"},{"id":10,"name":"产品"},{"id":11,"name":"后端"},{"id":12,"name":"linux"},{"id":13,"name":"人工智能"},{"id":14,"name":"php"},{"id":15,"name":"javascript"},{"id":16,"name":"架构"},{"id":17,"name":"前端"},{"id":18,"name":"python"},{"id":19,"name":"java"},{"id":20,"name":"算法"},{"id":21,"name":"面试"},{"id":22,"name":"科技动态"},{"id":23,"name":"js"},{"id":24,"name":"设计"},{"id":25,"name":"数码产品"},{"id":26,"name":"软件测试"}]},"message":"OK"}`
	c.String(200, res)
}

//返回文章内容

func HandleArticles(c *gin.Context) {
	var res Article_res
	page := c.Query("page")
	per_page := c.Query("per_page")
	res.Data.Page, _ = strconv.Atoi(page)
	res.Data.Per_page, _ = strconv.Atoi(per_page)
	res.Data.Total_count = 15
	results := []Article_content{}
	result := Article_content{
			ID:            "8218",
			Title:         "加载",
			Comment_count: 0,
			Pubdate:       "2019-03-11 09:00:00",
			Read_count:    2,
			Status:        2,
			Article_cover: Article_cover{
				Type:   "3",
				Iamges: []string{"http://geek.itheima.net/resources/images/15.jpg"},
			},
		}
	for i:=0;i<15;i++{
		result2 := result
		result2.Title = result2.Title+strconv.Itoa(i)
		result2.ID = strconv.Itoa(i)
		results = append(results, result2)
	}
	res.Data.Results = results
	res.Message = "OK"
	c.JSON(200, res)

}


//上传文章内容

func HandleUpload(c *gin.Context){

}

//更新文章内容

func HandleUpdate(c *gin.Context) {
	
}