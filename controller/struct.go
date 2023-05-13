package controller

type User_info struct {
	Mobile string `form:"mobile" json:"mobile" uri:"mobile" xml:"mobile" binding:"required"`
	Code   string `form:"code" json:"code" uri:"code" xml:"code" binding:"required"`
}

type Register_captcha struct {
	Mobile string `form:"mobile" json:"mobile" uri:"mobile" xml:"mobile" binding:"required"`
	Email   string `form:"email" json:"email" uri:"email" xml:"email" binding:"required"`
}

type Register_info struct {
	Mobile string `form:"mobile" json:"mobile" uri:"mobile" xml:"mobile" binding:"required"`
	Email   string `form:"email" json:"email" uri:"email" xml:"email" binding:"required"`
	Password string `form:"password" json:"password" uri:"password" xml:"password" binding:"required"`
	Captcha string `form:"captcha" json:"captcha" uri:"captcha" xml:"captcha" binding:"required"`
}

type Article_content struct{
	ID string `json:"id"`
	Title string `json:"title"`
	Comment_count int `json:"comment_count"`
	Status int `json:"status"`
	Pubdate string `json:"pubdate"`
	Like_count int `json:"like_count"`
	Read_count int `json:"read_count"`
	Article_cover `json:"cover"`
}

type Article_cover struct {
	Type string `json:"type"`
	Iamges []string `json:"images"`
}

type Article_data struct {
	Page int `json:"page"`
	Per_page int `json:"per_page"`
	Results []Article_content `json:"results"`
	Total_count int `json:"total_count"`
}

type Article_res struct {
	Data Article_data `json:"data"`
	Message string `json:"message"`
}

/*
{
	"channel_id": 1, 分表
	ID 文章ID 保证各不相同 生成uuid
	"content": "<p>4524545254646</p>", 文章内容
	"title": "assd", 文章标题
	"type": 1, 图片数目
	"cover": {
		"type": 1, 图片数目
		"images": ["http://geek.itheima.net/uploads/1683908806723.jpg"] 图片
	}
	Like_count int `json:"like_count"` 点赞数
	Read_count int `json:"read_count"` 阅读数
	Comment_count int `json:"comment_count"` 评论数
	Status int `json:"status"`					审核状态
	Pubdate string `json:"pubdate"`  发布时间 "2023-05-13 00:26:51"

}
*/