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
	Type int `json:"type"`
	Iamges []string `json:"images"`
}

type Article_data struct {
	Page int `json:"page"`
	Per_page int `json:"per_page"`
	Results []Article_content `json:"results"`
	Total_count int `json:"total_count"`
}

type Response struct {
	Data interface{} `json:"data"`
	Message string `json:"message"`
}

type Article_upload struct {
	Title string `json:"title"`
	Channel_id int `json:"channel_id"`
	Article_cover `json:"cover"`
	Type int `json:"type"`
	Content string `json:"content"`
	Id string `json:"id"`
}

type Article_update struct {
	Title string `json:"title"`
	Channel_id int `json:"channel_id"`
	Article_cover `json:"cover"`
	Content string `json:"content"`
	Id string `json:"id"`
	Pubdate string `json:"pubdate"`
}

type Article_update_put struct {
	Title string `json:"title"`
	Channel_id int `json:"channel_id"`
	Article_cover `json:"cover"`
	Content string `json:"content"`
}


/*
 const params = {
      channel_id,
      content,
      title,
      type,
      cover: {
        type: type,
        images: fileList.map(item => item.url)
      }
    }
*/