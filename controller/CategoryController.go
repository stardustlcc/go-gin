package controller

import (
	"github.com/gin-gonic/gin"
	"godev/model"
	"godev/repository"
	"godev/response"
	"godev/vo"
	"strconv"
)

type ICategoryController interface {
	RestController
}

type CategoryController struct {
	Repository repository.CategoryRepository
}

func NewCategoryController() ICategoryController {
	repository := repository.NewCategoryRepository()
	repository.DB.AutoMigrate(&model.Category{})
	return CategoryController{Repository:repository}
}

func (c CategoryController) Create(ctx *gin.Context)  {

	var requestCategory vo.CreateCategoryRequest
	if err := ctx.ShouldBind(&requestCategory); err != nil {
		response.Fail(ctx, "数据验证错误，分类名称必填", nil)
		return
	}

	category, err := c.Repository.Create(requestCategory.Name)
	if err != nil {
		panic(err)
		return
	}

	response.Success(ctx, gin.H{"category": category}, "")

}

func (c CategoryController) Update(ctx *gin.Context)  {
	// 绑定body 中的参数
	var requestCategory vo.CreateCategoryRequest
	if err := ctx.ShouldBind(&requestCategory); err != nil {
		response.Fail(ctx, "数据验证错误，分类名称必填", nil)
		return
	}

	// 获取 path 中的参数
	categoryId, _ := strconv.Atoi(ctx.Params.ByName("id"))
	updateCategory, err := c.Repository.SelectById(categoryId)

	if err != nil {
		response.Fail(ctx, "分类不存在",nil)
		return
	}

	//更新分类
	category, err := c.Repository.Update(*updateCategory, requestCategory.Name)
	if err != nil {
		panic(err)
	}

	response.Success(ctx, gin.H{"category": category}, "修改成功")
}

func (c CategoryController) Show(ctx *gin.Context)  {
	//获取 path 中的参数
	categoryId, _ := strconv.Atoi(ctx.Params.ByName("id"))

	infoCategory, err := c.Repository.SelectById(categoryId)

	if err != nil {
		response.Fail(ctx, "分类不存在",nil)
		return
	}

	response.Success(ctx, gin.H{"category":infoCategory}, "")
}

func (c CategoryController) Delete(ctx *gin.Context)  {
	//获取 path 中的参数
	categoryId, _ := strconv.Atoi(ctx.Params.ByName("id"))
	err := c.Repository.DeleteById(categoryId)
	if err != nil {
		response.Fail(ctx, "删除失败，请稍后再试",nil)
		return
	}
	response.Success(ctx, nil, "")
}
