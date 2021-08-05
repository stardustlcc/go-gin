package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"godev/common"
	"godev/model"
	"godev/response"
	"godev/vo"
	"log"
	"strconv"
)

type IPostController interface {
	RestController
	PageList(ctx *gin.Context)
}

type PostController struct {
	DB *gorm.DB
}

func NewPostController() IPostController  {
	db := common.GetDb()
	db.AutoMigrate(&model.Post{})
	return PostController{DB:db}
}

func (p PostController) Create(ctx *gin.Context)  {
	var requestPost vo.CreatePostRequest
	if err := ctx.ShouldBind(&requestPost); err != nil {
		log.Println(err.Error())
		response.Fail(ctx, "数据验证错误", nil)
		return
	}

	//获取用户信息
	user, _ := ctx.Get("user")

	//创建post
	post := model.Post{
		UserId:     user.(model.User).ID,
		CategoryId: requestPost.CategoryId,
		Title:      requestPost.Title,
		HeadImg:    requestPost.HeadImg,
		Content:    requestPost.Content,
	}

	if err := p.DB.Create(&post).Error; err != nil {
		panic(err)
		return
	}

	response.Success(ctx, gin.H{"id": post.ID},"创建成功")
}

func (p PostController) Update(ctx *gin.Context)  {
	var requestPost vo.CreatePostRequest
	if err := ctx.ShouldBind(&requestPost); err != nil {
		log.Println(err.Error())
		response.Fail(ctx, "数据验证错误", nil)
		return
	}

	//获取path 中的id
	postId := ctx.Params.ByName("id")
	var post model.Post
	if p.DB.Where("id=?", postId).First(&post).RecordNotFound() {
		response.Fail(ctx, "文章不存在", nil)
		return
	}

	//判断当前用户是否是文字的作者
	//获取用户信息
	user, _ := ctx.Get("user")
	userId := user.(model.User).ID
	if userId != post.UserId {
		response.Fail(ctx, "文章不属于您，请勿非法操作", nil)
		return
	}

	//更新文章
	if err := p.DB.Model(&post).Update(requestPost).Error; err != nil {
		response.Fail(ctx, "更新失败", nil)
		return
	}

	response.Success(ctx, gin.H{"post": post}, "更新成功")

}

func (p PostController) Show(ctx *gin.Context)  {

	//获取path 中的id
	postId := ctx.Params.ByName("id")
	var post model.Post
	if p.DB.Preload("Category").Where("id=?", postId).First(&post).RecordNotFound() {
		response.Fail(ctx, "文章不存在", nil)
		return
	}
	response.Success(ctx, gin.H{"post": post}, "更新成功")
}

func (p PostController) Delete(ctx *gin.Context)  {

	//获取path 中的id
	postId := ctx.Params.ByName("id")
	var post model.Post
	if p.DB.Where("id=?", postId).First(&post).RecordNotFound() {
		response.Fail(ctx, "文章不存在", nil)
		return
	}

	//判断当前用户是否是文字的作者
	//获取用户信息
	user, _ := ctx.Get("user")
	userId := user.(model.User).ID
	if userId != post.UserId {
		response.Fail(ctx, "文章不属于您，请勿非法操作", nil)
		return
	}

	p.DB.Delete(&post)

	response.Success(ctx, gin.H{"res": true}, "删除成功")
}

func (p PostController) PageList(ctx *gin.Context)  {
	//获取分页参数
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum","1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize","20"))

	//分页
	var posts []model.Post
	p.DB.Order("created_at desc").Offset((pageNum-1) * pageSize).Limit(pageSize).Find(&posts)

	//前端渲染需要总数
	var total int
	p.DB.Model(model.Post{}).Count(&total)

	response.Success(ctx, gin.H{"data": posts, "total": total}, "成功")

}
