package controller
//
//import (
//	"github.com/gin-gonic/gin"
//	"github.com/jinzhu/gorm"
//	"godev/common"
//	"godev/model"
//	"godev/response"
//	"godev/vo"
//	"strconv"
//)
//
//type ICategoryController interface {
//	RestController
//}
//
//type CategoryController struct {
//	DB *gorm.DB
//}
//
//func NewCategoryController() ICategoryController {
//	db := common.GetDb()
//	db.AutoMigrate(&model.Category{})
//	return CategoryController{DB:db}
//}
//
//func (c CategoryController) Create(ctx *gin.Context)  {
//
//	var requestCategory vo.CreateCategoryRequest
//	if err := ctx.ShouldBind(&requestCategory); err != nil {
//		response.Fail(ctx, "数据验证错误，分类名称必填", nil)
//		return
//	}
//
//	category := model.Category{Name:requestCategory.Name}
//	c.DB.Create(&category)
//
//	response.Success(ctx, gin.H{"category": category}, "")
//}
//
//func (c CategoryController) Update(ctx *gin.Context)  {
//	// 绑定body 中的参数
//	var requestCategory vo.CreateCategoryRequest
//	if err := ctx.ShouldBind(&requestCategory); err != nil {
//		response.Fail(ctx, "数据验证错误，分类名称必填", nil)
//		return
//	}
//
//	// 获取 path 中的参数
//	categoryId, _ := strconv.Atoi(ctx.Params.ByName("id"))
//	var updateCategory model.Category
//	if c.DB.First(&updateCategory, categoryId).RecordNotFound() {
//		response.Fail(ctx, "分类不存在",nil)
//		return
//	}
//
//	//更新分类
//	c.DB.Model(&updateCategory).Update("name", requestCategory.Name)
//
//	response.Success(ctx, gin.H{"category": updateCategory}, "修改成功")
//}
//
//func (c CategoryController) Show(ctx *gin.Context)  {
//	//获取 path 中的参数
//	categoryId, _ := strconv.Atoi(ctx.Params.ByName("id"))
//
//	var infoCategory model.Category
//	if c.DB.First(&infoCategory, categoryId).RecordNotFound() {
//		response.Fail(ctx, "分类不存在", nil)
//		return
//	}
//
//	response.Success(ctx, gin.H{"category":infoCategory}, "")
//}
//
//func (c CategoryController) Delete(ctx *gin.Context)  {
//	//获取 path 中的参数
//	categoryId, _ := strconv.Atoi(ctx.Params.ByName("id"))
//
//	if err := c.DB.Delete(model.Category{}, categoryId).Error; err != nil {
//		response.Fail(ctx, "删除失败，请重试", nil)
//		return
//	}
//
//	response.Success(ctx, nil, "")
//}
