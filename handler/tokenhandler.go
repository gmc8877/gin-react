package handler

import (
	"DEMO01/tools"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// 用户信息类，作为生成token的参数
type UserClaims struct {
	ID       string `json:"userId"`
	Name     string `json:"name"`
	Password string `json:"password"`
	//jwt-go提供的标准claim
	jwt.StandardClaims
}

var (
	//自定义的token秘钥
	secret = []byte("464fsd895sfdf48569")
	//该路由下不校验token
	noVerify = []string{"/login","/register/captcha","/register", "/root/login"}
	//token有效时间（纳秒）
	effectTime = 24 * time.Hour
)

// 生成token
func GenerateToken(claims *UserClaims) string {
	//设置token有效期，也可不设置有效期，采用redis的方式
	//   1)将token存储在redis中，设置过期时间，token如没过期，则自动刷新redis过期时间，
	//   2)通过这种方式，可以很方便的为token续期，而且也可以实现长时间不登录的话，强制登录
	//本例只是简单采用 设置token有效期的方式，只是提供了刷新token的方法，并没有做续期处理的逻辑
	claims.ExpiresAt = time.Now().Add(effectTime).Unix()
	//生成token
	sign, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(secret)

	//这里因为项目接入了统一异常处理，所以使用panic并不会使程序终止，如不接入，可使用原始方式处理错误
	//接入统一异常可参考 https://blog.csdn.net/u014155085/article/details/106733391
	tools.CheckErr(err)

	return sign
}

// 验证token
func JwtVerify(c *gin.Context) {
	//过滤是否验证token
  for _,url:= range noVerify {
    if url==c.Request.RequestURI {
      return
    }
  }
	token := c.GetHeader("token")
	
	if token == "" {
		c.JSON(401, gin.H{
			"message":"token not exist",
		})
		panic("token not exist !")
	}
	//验证token，并存储在请求中
	claims := parseToken(token)
	var password, user_name string
	//数据库查询用户名和密码
	err := DB.QueryRow("select password,user_name from user_info where user_id=?",claims.ID).Scan(&password, &user_name)
	tools.CheckErr(err)
	if claims.Name==user_name&&claims.Password==password {
		c.Set("user", claims)
	}else {
		c.JSON(401, gin.H{
			"message":"token not right",
		})
		panic("token not right !")
	}
	
}

func RootVerify(c *gin.Context) {
	//过滤是否验证token
  for _,url:= range noVerify {
    if url==c.Request.RequestURI {
      return
    }
  }
	token := c.GetHeader("token")
	
	if token == "" {
		c.JSON(401, gin.H{
			"message":"token not exist",
		})
		panic("token not exist !")
	}
	//验证token，并存储在请求中
	claims := parseToken(token)
	var password, user_name string
	//数据库查询用户名和密码
	err := DB.QueryRow("select password,user_name from root_info where user_id=?",claims.ID).Scan(&password, &user_name)
	tools.CheckErr(err)
	if claims.Name==user_name&&claims.Password==password {
		c.Set("user", claims)
	}else {
		c.JSON(401, gin.H{
			"message":"token not right",
		})
		panic("token not right !")
	}
	
}



// 解析Token
func parseToken(tokenString string) *UserClaims {
	//解析token
	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	tools.CheckErr(err)
	claims, ok := token.Claims.(*UserClaims)
	if !ok {
		panic("token is valid")
	}
	return claims
}

// 更新token
func Refresh(tokenString string) string {
	jwt.TimeFunc = func() time.Time {
		return time.Unix(0, 0)
	}
	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	tools.CheckErr(err)
	claims, ok := token.Claims.(*UserClaims)
	if !ok {
		panic("token is valid")
	}
	jwt.TimeFunc = time.Now
	claims.StandardClaims.ExpiresAt = time.Now().Add(24 * time.Hour).Unix()
	return GenerateToken(claims)
}
