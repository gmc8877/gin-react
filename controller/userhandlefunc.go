package controller

import (
	"DEMO01/handler"
	"DEMO01/tools"
	"context"
	"fmt"
	"io/ioutil"

	"os"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/satori/go.uuid"
)

var DB_NUM int = 5

//返回图片

func Getimages(c *gin.Context) {
	imageName := "./assets/" + c.Param("path")
	file, err := ioutil.ReadFile(imageName)
	tools.CheckErr(err)
	c.Writer.WriteString(string(file))
}

//用户登录
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

	// 必须，设置服务器支持的所有跨域请求的方法
	context.Header("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE, OPTIONS")
	// 服务器支持的所有头信息字段，不限于浏览器在"预检"中请求的字段
	context.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Token, X-Requested-With")
	// 可选，设置XMLHttpRequest的响应对象能拿到的额外字段
	// context.Header("Access-Control-Expose-Headers", "Access-Control-Allow-Headers, Token, X-Requested-With")
	// 放行所有OPTIONS方法
	if method == "OPTIONS" {
		context.AbortWithStatus(200)
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
	val, err := rdb.Get(ctx, num_key).Result()
	if val != "" {
		if err != nil {
			tools.CheckErr(err)
		}
		return
	}
	num, err := tools.SendEmailValidate([]string{registercaptcha.Email})
	tools.CheckErr(err)
	err = rdb.SetNX(ctx, num_key, num, 5*time.Minute).Err()
	tools.CheckErr(err)
}

// 返回用户个人信息
func HandleUserinfo(c *gin.Context) {
	user_info, _ := c.Get("user")
	claims := user_info.(*handler.UserClaims)
	c.JSON(200, gin.H{
		"name": claims.Name,
	})
}

//返回root用户个人信息

func HandleRootinfo(c *gin.Context) {
	user_info, _ := c.Get("user")
	claims := user_info.(*handler.UserClaims)
	c.JSON(200, gin.H{
		"name": claims.Name,
	})
}

// 返回文章分类
func HandleChannels(c *gin.Context) {
	res := `{"data":{"channels":[{"id":0,"name":"推荐"},{"id":1,"name":"html"},{"id":2,"name":"开发者资讯"},{"id":3,"name":"go"},{"id":4,"name":"c++"}]},"message":"OK"}`
	c.String(200, res)
}

//返回修改面板文章列表

func HandleArticlesList(c *gin.Context) {
	var res Response

	page := c.Query("page")
	per_page := c.Query("per_page")
	data := Article_data{}
	data.Page, _ = strconv.Atoi(page)
	data.Per_page, _ = strconv.Atoi(per_page)

	results := []Article_content{}
	//数据库查询
	channel_id, ok_c := c.GetQuery("channel_id")
	status, ok_s := c.GetQuery("status")
	begin_pubdate, ok_b := c.GetQuery("begin_pubdate")
	end_pubdate := c.Query("end_pubdate")
	var sql_l string
	db := handler.DB
	if ok_c {
		if ok_s {
			if ok_b {
				status_v, _ := strconv.Atoi(status)
				sql_l = fmt.Sprintf("select uuid, title, Comment_count, Read_count, Like_count, Status, Type, Images_url, pubdate from channel_%s where status=%v && pubdate>='%s' && pubdate<='%s' ", channel_id, status_v, begin_pubdate, end_pubdate)
			} else {
				status_v, _ := strconv.Atoi(status)
				sql_l = fmt.Sprintf("select uuid, title, Comment_count, Read_count, Like_count, Status, Type, Images_url, pubdate from channel_%s where status=%v ", channel_id, status_v)
			}
		} else {
			if ok_b {
				sql_l = fmt.Sprintf("select uuid, title, Comment_count, Read_count, Like_count, Status, Type, Images_url, pubdate from channel_%s where pubdate>='%s' && pubdate<='%s' ", channel_id, begin_pubdate, end_pubdate)
			} else {
				sql_l = fmt.Sprintf("select uuid, title, Comment_count, Read_count, Like_count, Status, Type, Images_url, pubdate from channel_%s", channel_id)
			}
		}
		rows, err := db.Query(sql_l)
		tools.CheckErr(err)
		for rows.Next() {
			//uuid, title, Comment_count, Read_count, Like_count, Status, Type, Iamges_url, pubdate
			var result Article_content
			var image_url, pubdate string
			err = rows.Scan(&result.ID, &result.Title, &result.Comment_count, &result.Read_count, &result.Like_count, &result.Status, &result.Type, &image_url, &pubdate)
			tools.CheckErr(err)
			result.Iamges = []string{image_url}
			result.Pubdate = pubdate[:10]
			results = append(results, result)
		}
	} else {
		for i := 0; i < DB_NUM; i++ {
			if ok_s {
				if ok_b {
					status_v, _ := strconv.Atoi(status)
					sql_l = fmt.Sprintf("select uuid, title, Comment_count, Read_count, Like_count, Status, Type, Images_url, pubdate from channel_%v where status=%v && pubdate>='%s' && pubdate<='%s' ", i, status_v, begin_pubdate, end_pubdate)
				} else {
					status_v, _ := strconv.Atoi(status)
					sql_l = fmt.Sprintf("select uuid, title, Comment_count, Read_count, Like_count, Status, Type, Images_url, pubdate from channel_%v where status=%v ", i, status_v)
				}
			} else {
				if ok_b {
					sql_l = fmt.Sprintf("select uuid, title, Comment_count, Read_count, Like_count, Status, Type, Images_url, pubdate from channel_%v where pubdate>='%s' && pubdate<='%s' ", i, begin_pubdate, end_pubdate)
				} else {
					sql_l = fmt.Sprintf("select uuid, title, Comment_count, Read_count, Like_count, Status, Type, Images_url, pubdate from channel_%v", i)
				}
			}
			rows, err := db.Query(sql_l)
			tools.CheckErr(err)
			for rows.Next() {
				//uuid, title, Comment_count, Read_count, Like_count, Status, Type, Iamges_url, pubdate
				var result Article_content
				var image_url, pubdate string

				err = rows.Scan(&result.ID, &result.Title, &result.Comment_count, &result.Read_count, &result.Like_count, &result.Status, &result.Type, &image_url, &pubdate)
				tools.CheckErr(err)
				result.Iamges = []string{image_url}
				result.Pubdate = pubdate[:10]
				results = append(results, result)
			}
		}
	}
	data.Total_count = len(results)
	data.Results = results
	res.Message = "OK"
	res.Data = data
	c.JSON(200, res)
}

//返回修改面板文章内容

func HandleArticles(c *gin.Context) {

	id := c.Param("id")
	db := handler.DB
	var article Article_update
	channel_id, err := strconv.Atoi(id[36:])
	tools.CheckErr(err)
	sql_l := fmt.Sprintf("select title, content, images_url from channel_%v where uuid=?", channel_id)
	var title, content, images_url string
	err = db.QueryRow(sql_l, id).Scan(&title, &content, &images_url)
	tools.CheckErr(err)
	article = Article_update{
		Id:         id,
		Title:      title,
		Channel_id: channel_id,
		Content:    content,
		Article_cover: Article_cover{
			Type:   0,
			Iamges: []string{images_url},
		},
	}
	var res Response
	res.Data = article
	res.Message = "OK"
	c.JSON(200, res)
}

//上传文章内容

func HandleUpload(c *gin.Context) {
	var article Article_upload
	c.Bind(&article)
	channel_id := article.Channel_id
	user_info, _ := c.Get("user")
	user_name := user_info.(*handler.UserClaims).Name
	db := handler.DB
	sql := fmt.Sprintf("insert into channel_%v (uuid, status, title, type, content, images_url,like_count,read_count,comment_count, pubdate, user_name)  VALUES (?,?, ?, ?, ?,?,?,?, ?, ?,?)", channel_id)
	smtp, err := db.Prepare(sql)
	tools.CheckErr(err)
	image_url := article.Article_cover.Iamges[0]
	id := article.Id + strconv.Itoa(channel_id)
	_, err = smtp.Exec(id, 2, article.Title, article.Type, article.Content, image_url, 0, 0, 0, time.Now().Format("2006-01-02"), user_name)
	tools.CheckErr(err)
	c.JSON(200, gin.H{
		"message": "OK",
	})
}


//更新文章内容

func HandleUpdate(c *gin.Context) {
	id := c.Param("id")
	var article Article_update_put
	c.Bind(&article)
	user_info, _ := c.Get("user")
	user_name := user_info.(*handler.UserClaims).Name
	db := handler.DB
	sql_l := fmt.Sprintf("update channel_%v set title=?,content=?,type=?,images_url=? , pubdate=? where uuid=? && user_name=?", article.Channel_id)
	smtp, err := db.Prepare(sql_l)
	tools.CheckErr(err)
	_, err = smtp.Exec(article.Title, article.Content, article.Type, article.Iamges[0], time.Now().Format("2006-01-02"), id, user_name)
	tools.CheckErr(err)

	data := make(map[string]string)
	data["id"] = id
	c.JSON(200, gin.H{
		"data":    data,
		"message": "OK",
	})
}



// 上传图片
func HandleImagesUpload(c *gin.Context) {
	id := c.PostForm("uuid")
	if len(id) == 0 {
		u2 := uuid.NewV4()
		id = u2.String()
	}
	image_id := id + ".jpg"
	file, _ := c.FormFile("image")
	dir := "./assets/" + image_id
	c.SaveUploadedFile(file, dir)
	data := make(map[string]string)
	data["url"] = "http://101.35.210.115/api/assets/" + image_id
	data["uuid"] = id
	c.JSON(200, gin.H{
		"data":    data,
		"message": "OK",
	})
}

//管理员删除文章

func HandleDelete(c *gin.Context) {
	id := c.Param("id")
	var article Article_update_put
	c.Bind(&article)
	db := handler.DB
	channel_id := id[36:]
	sql_l := fmt.Sprintf("delete from channel_%s where uuid=?", channel_id)
	smtp, err := db.Prepare(sql_l)
	tools.CheckErr(err)
	_, err = smtp.Exec(id)
	tools.CheckErr(err)
	os.Remove("./assets/" + id[:36] + ".jpg")
	data := make(map[string]string)
	data["id"] = id
	c.JSON(200, gin.H{
		"data":    data,
		"message": "OK",
	})
}
 //用户删除文章
func HandleUsrDelete(c *gin.Context) {
	id := c.Param("id")
	user_info, _ := c.Get("user")
	user_name := user_info.(*handler.UserClaims).Name
	var article Article_update_put
	c.Bind(&article)
	db := handler.DB
	channel_id := id[36:]
	sql_l := fmt.Sprintf("delete from channel_%s where uuid=? && user_name=?", channel_id)
	smtp, err := db.Prepare(sql_l)
	tools.CheckErr(err)
	_, err = smtp.Exec(id, user_name)
	tools.CheckErr(err)
	os.Remove("./assets/" + id[:36] + ".jpg")
	data := make(map[string]string)
	data["id"] = id
	c.JSON(200, gin.H{
		"data":    data,
		"message": "OK",
	})
}

// 返回展示页面文章列表
func HandleShows(c *gin.Context) {
	var res Response

	page := c.Query("page")
	per_page := c.Query("per_page")
	data := Article_data_show{}
	data.Page, _ = strconv.Atoi(page)
	data.Per_page, _ = strconv.Atoi(per_page)

	results := []Article_content_show{}
	//数据库查询
	channel_id, ok_c := c.GetQuery("channel_id")
	status, ok_s := c.GetQuery("status")
	begin_pubdate, ok_b := c.GetQuery("begin_pubdate")
	end_pubdate := c.Query("end_pubdate")
	var sql_l string
	db := handler.DB
	if ok_c {
		if ok_s {
			if ok_b {
				status_v, _ := strconv.Atoi(status)
				sql_l = fmt.Sprintf("select uuid, title, Comment_count, Read_count, Like_count, Status,  Images_url, pubdate , user_name from channel_%s where status=%v && pubdate>='%s' && pubdate<='%s' ", channel_id, status_v, begin_pubdate, end_pubdate)
			} else {
				status_v, _ := strconv.Atoi(status)
				sql_l = fmt.Sprintf("select uuid, title, Comment_count, Read_count, Like_count, Status, Images_url, pubdate  , user_name from channel_%s where status=%v ", channel_id, status_v)
			}
		} else {
			if ok_b {
				sql_l = fmt.Sprintf("select uuid, title, Comment_count, Read_count, Like_count, Status,Images_url, pubdate  , user_name from channel_%s where pubdate>='%s' && pubdate<='%s' ", channel_id, begin_pubdate, end_pubdate)
			} else {
				sql_l = fmt.Sprintf("select uuid, title, Comment_count, Read_count, Like_count, Status, Images_url, pubdate  , user_name from channel_%s", channel_id)
			}
		}
		rows, err := db.Query(sql_l)
		tools.CheckErr(err)
		for rows.Next() {
			//uuid, title, Comment_count, Read_count, Like_count, Status, Type, Iamges_url, pubdate
			var result Article_content_show
			var image_url, pubdate string
			err = rows.Scan(&result.ID, &result.Title, &result.Comment_count, &result.Read_count, &result.Like_count, &result.Status, &image_url, &pubdate, &result.Name)
			tools.CheckErr(err)
			result.Image = image_url
			result.Pubdate = pubdate
			results = append(results, result)
		}
	} else {
		for i := 0; i < DB_NUM; i++ {
			if ok_s {
				if ok_b {
					status_v, _ := strconv.Atoi(status)
					sql_l = fmt.Sprintf("select uuid, title, Comment_count, Read_count, Like_count, Status,  Images_url, pubdate , user_name  from channel_%v where status=%v && pubdate>='%s' && pubdate<='%s'", i, status_v, begin_pubdate, end_pubdate)
				} else {
					status_v, _ := strconv.Atoi(status)
					sql_l = fmt.Sprintf("select uuid, title, Comment_count, Read_count, Like_count, Status, Images_url, pubdate  , user_name from channel_%v where status=%v ", i, status_v)
				}
			} else {
				if ok_b {
					sql_l = fmt.Sprintf("select uuid, title, Comment_count, Read_count, Like_count, Status,  Images_url, pubdate  , user_name from channel_%v where pubdate>='%s' && pubdate<='%s' ", i, begin_pubdate, end_pubdate)
				} else {
					sql_l = fmt.Sprintf("select uuid, title, Comment_count, Read_count, Like_count, Status, Images_url, pubdate  , user_name from channel_%v", i)
				}
			}
			rows, err := db.Query(sql_l)
			tools.CheckErr(err)
			for rows.Next() {
				//uuid, title, Comment_count, Read_count, Like_count, Status, Type, Iamges_url, pubdate
				var result Article_content_show
				var image_url, pubdate string

				err = rows.Scan(&result.ID, &result.Title, &result.Comment_count, &result.Read_count, &result.Like_count, &result.Status, &image_url, &pubdate, &result.Name)
				tools.CheckErr(err)
				result.Image = image_url
				result.Pubdate = pubdate
				results = append(results, result)
			}
		}
	}
	data.Total_count = len(results)
	data.Results = results
	res.Message = "OK"
	res.Data = data
	c.JSON(200, res)
}

//处理管理员账户登录

func HandleRootLogin(c *gin.Context) {
	//登录校验
	db := handler.DB
	var user User_info
	err := c.ShouldBindJSON(&user)
	tools.CheckErr(err)
	//密码进行hash
	code_hash := tools.CodeHash(user.Code)
	var password, user_id string
	//数据库查询用户名和密码
	err = db.QueryRow("select password, user_id from root_info where user_name=?", user.Mobile).Scan(&password, &user_id)
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

//返回用户文章列表

func HandleUsrArticlesList(c *gin.Context) {
	var res Response
	page := c.Query("page")
	per_page := c.Query("per_page")
	data := Article_data{}
	data.Page, _ = strconv.Atoi(page)
	data.Per_page, _ = strconv.Atoi(per_page)
	results := []Article_content{}
	//用户名
	user_info, _ := c.Get("user")
	claims := user_info.(*handler.UserClaims)
	usr_name := claims.Name
	//数据库查询
	channel_id, ok_c := c.GetQuery("channel_id")
	status, ok_s := c.GetQuery("status")
	begin_pubdate, ok_b := c.GetQuery("begin_pubdate")
	end_pubdate := c.Query("end_pubdate")
	var sql_l string
	db := handler.DB
	if ok_c {
		if ok_s {
			if ok_b {
				status_v, _ := strconv.Atoi(status)
				sql_l = fmt.Sprintf("select uuid, title, Comment_count, Read_count, Like_count, Status, Type, Images_url, pubdate from channel_%s where status=%v && pubdate>='%s' && pubdate<='%s' && user_name=%s", channel_id, status_v, begin_pubdate, end_pubdate, usr_name)
			} else {
				status_v, _ := strconv.Atoi(status)
				sql_l = fmt.Sprintf("select uuid, title, Comment_count, Read_count, Like_count, Status, Type, Images_url, pubdate from channel_%s where status=%v && user_name=%s", channel_id, status_v, usr_name)
			}
		} else {
			if ok_b {
				sql_l = fmt.Sprintf("select uuid, title, Comment_count, Read_count, Like_count, Status, Type, Images_url, pubdate from channel_%s where pubdate>='%s' && pubdate<='%s' && user_name=%s ", channel_id, begin_pubdate, end_pubdate, usr_name)
			} else {
				sql_l = fmt.Sprintf("select uuid, title, Comment_count, Read_count, Like_count, Status, Type, Images_url, pubdate from channel_%s && user_name=%s", channel_id, usr_name)
			}
		}
		rows, err := db.Query(sql_l)
		tools.CheckErr(err)
		for rows.Next() {
			//uuid, title, Comment_count, Read_count, Like_count, Status, Type, Iamges_url, pubdate
			var result Article_content
			var image_url, pubdate string
			err = rows.Scan(&result.ID, &result.Title, &result.Comment_count, &result.Read_count, &result.Like_count, &result.Status, &result.Type, &image_url, &pubdate)
			tools.CheckErr(err)
			result.Iamges = []string{image_url}
			result.Pubdate = pubdate[:10]
			results = append(results, result)
		}
	} else {
		for i := 0; i < DB_NUM; i++ {
			if ok_s {
				if ok_b {
					status_v, _ := strconv.Atoi(status)
					sql_l = fmt.Sprintf("select uuid, title, Comment_count, Read_count, Like_count, Status, Type, Images_url, pubdate from channel_%v where status=%v && pubdate>='%s' && pubdate<='%s' && user_name=%s", i, status_v, begin_pubdate, end_pubdate, usr_name)
				} else {
					status_v, _ := strconv.Atoi(status)
					sql_l = fmt.Sprintf("select uuid, title, Comment_count, Read_count, Like_count, Status, Type, Images_url, pubdate from channel_%v where status=%v && user_name=%s", i, status_v, usr_name)
				}
			} else {
				if ok_b {
					sql_l = fmt.Sprintf("select uuid, title, Comment_count, Read_count, Like_count, Status, Type, Images_url, pubdate from channel_%v where pubdate>='%s' && pubdate<='%s' && user_name=%s", i, begin_pubdate, end_pubdate, usr_name)
				} else {
					sql_l = fmt.Sprintf("select uuid, title, Comment_count, Read_count, Like_count, Status, Type, Images_url, pubdate from channel_%v where user_name=%s", i, usr_name)
				}
			}
			rows, err := db.Query(sql_l)
			tools.CheckErr(err)
			for rows.Next() {
				//uuid, title, Comment_count, Read_count, Like_count, Status, Type, Iamges_url, pubdate
				var result Article_content
				var image_url, pubdate string

				err = rows.Scan(&result.ID, &result.Title, &result.Comment_count, &result.Read_count, &result.Like_count, &result.Status, &result.Type, &image_url, &pubdate)
				tools.CheckErr(err)
				result.Iamges = []string{image_url}
				result.Pubdate = pubdate[:10]
				results = append(results, result)
			}
		}
	}
	data.Total_count = len(results)
	data.Results = results
	res.Message = "OK"
	res.Data = data
	c.JSON(200, res)
}