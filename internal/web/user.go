package web

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"example.com/basic-gin/webook/internal/domain"
	"example.com/basic-gin/webook/internal/service"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

const (
	emailRegexPattern = "^\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*$"
	//go正统的正则包，不支持“?=”这种写法，需要引入“github.com/dlclark/regexp2”包
	passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,}$`
	phoneRegexPattern    = `^1[3-9]\d{9}$`
	//日期正则校验的格式是：yyyy-MM-DD
	dateRegexPattern = `^(?:(?!0000)[0-9]{4}\-(?:(?:0[13578]|1[02])(?:\-0[1-9]|\-[12][0-9]|\-3[01])|(?:0[469]|11)(?:\-0[1-9]|\-[12][0-9]|\-30)|02(?:\-0[1-9]|\-1[0-9]|\-2[0-8]))|(?:(((\d{2})(0[48]|[2468][048]|[13579][26])|(([02468][048])|([13579][26]))00))\-02\-29))$`
)

type UserHandler struct {
	emailRegexExp    *regexp.Regexp //预编译正则表达式来提高校验速度
	passwordRegexExp *regexp.Regexp
	phoneRegexExp    *regexp.Regexp
	dateRegexExp     *regexp.Regexp
	svc              *service.UserService
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{
		emailRegexExp:    regexp.MustCompile(emailRegexPattern, regexp.None),
		passwordRegexExp: regexp.MustCompile(passwordRegexPattern, regexp.None),
		phoneRegexExp:    regexp.MustCompile(phoneRegexPattern, regexp.None),
		dateRegexExp:     regexp.MustCompile(dateRegexPattern, regexp.None),
		svc:              svc,
	}
}

// （分散式）注册路由
func (h *UserHandler) RegisterRoutes(server *gin.Engine) {
	ug := server.Group("/users") //分组路由
	ug.POST("/signup", h.Signup)
	//server.POST("/users/signup", h.Signup) //未分组的情况

	ug.POST("/login", h.Login)
	ug.POST("/edit", h.Edit)
	ug.GET("/profile", h.Profile)
}

// 注册
func (h *UserHandler) Signup(ctx *gin.Context) {
	//内部类，除了方法Signup，谁都不能用
	type SignupReq struct {
		Email           string `json:"email"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirmPassword"`
	}
	var req SignupReq
	if err := ctx.Bind(&req); err != nil {
		fmt.Printf("error: %v", err)
		return //报错
	}

	isEmail, err := h.emailRegexExp.MatchString(req.Email)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	if !isEmail {
		ctx.String(http.StatusOK, "非法有想格式")
		return
	}

	isPassword, err := h.passwordRegexExp.MatchString(req.Password)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	if !isPassword {
		ctx.String(http.StatusOK, "密码必须包含数字、特殊字符，并且长度不能小于8位")
		return
	}

	if req.Password != req.ConfirmPassword {
		ctx.String(http.StatusOK, "确认密码不一致")
		return
	}

	//验证完成之后，就开始实际的“注册”逻辑
	err = h.svc.Signup(ctx, domain.User{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		ctx.String(http.StatusOK, "注册失败")
	}

	ctx.String(http.StatusOK, "注册成功")
}

// 登陆
func (h *UserHandler) Login(ctx *gin.Context) {
	type LoginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req LoginReq
	if err := ctx.Bind(&req); err != nil {
		ctx.String(http.StatusOK, "")
		return
	}

	u, err := h.svc.Login(ctx, req.Email, req.Password)
	switch err {
	case nil:
		//登陆成功
		sess := sessions.Default(ctx)
		sess.Set("userId", u.Id)
		sess.Options(sessions.Options{
			MaxAge: 15 * int(time.Minute),
		})
		err = sess.Save()
		if err != nil {
			ctx.String(http.StatusOK, "系统错误")
			return
		}
		ctx.String(http.StatusOK, "登陆成功")
	case service.ErrInvalidUserOrPassword:
		ctx.String(http.StatusOK, "用户名或者密码不正确")
	default:
		ctx.String(http.StatusOK, "系统错误")
	}
}

// 个人信息
func (h *UserHandler) Profile(ctx *gin.Context) {
	uid, err := h.getUidBySession(ctx)
	if err != nil {
		ctx.String(http.StatusInternalServerError, "系统错误") //通过了登陆中间件的校验，但是咩有获取到uid>>系统错误
		return
	}

	u, err := h.svc.Profile(ctx, uid)
	if err != nil {
		ctx.String(http.StatusInternalServerError, "没找到用户")
		return
	}

	type Profile struct {
		Email    string `json:"email"`
		Nickname string `json:"nickname"`
		Birthday string `json:"birthday"`
		Gender   int8   `json:"gender"`
		Phone    string `json:"phone"`
	}

	ctx.JSON(http.StatusOK, Profile{
		Email:    u.Email,
		Nickname: u.Nickname,
		Birthday: u.Birthday,
		Gender:   u.Gender,
		Phone:    u.Phone,
	})
}

// 修改个人信息
func (h *UserHandler) Edit(ctx *gin.Context) {
	uid, err := h.getUidBySession(ctx)
	if err != nil {
		ctx.String(http.StatusInternalServerError, "系统错误") //通过了登陆中间件的校验，但是咩有获取到uid>>系统错误
		return
	}

	type EditReq struct {
		Nickname string `json:"nickname"`
		Birthday string `json:"birthday"`
		Gender   int8   `json:"gender"`
		Phone    string `json:"phone"`
	}
	var req EditReq
	if err := ctx.Bind(&req); err != nil {
		ctx.String(http.StatusOK, "参数不正确")
		return
	}

	if req.Phone != "" {
		isPhone, err := h.phoneRegexExp.MatchString(req.Phone)
		if err != nil {
			ctx.String(http.StatusOK, "系统错误")
			return
		}
		if !isPhone {
			ctx.String(http.StatusOK, "手机格式不正确")
			return
		}
	}

	if req.Birthday != "" {
		isDate, err := h.dateRegexExp.MatchString(req.Birthday)
		if err != nil {
			ctx.String(http.StatusOK, "系统错误")
			return
		}
		if !isDate {
			ctx.String(http.StatusOK, "生日格式不正确")
			return
		}
	}

	if req.Gender != 0 {
		gender := int8(req.Gender)
		if gender != 1 && gender != 2 {
			ctx.String(http.StatusOK, "性别数据非法")
			return
		}
	}

	err = h.svc.Edit(ctx, domain.User{
		Id:       uid,
		Nickname: req.Nickname,
		Birthday: req.Birthday,
		Gender:   req.Gender,
		Phone:    req.Phone,
	})
	if err != nil {
		ctx.String(http.StatusOK, "更新个人信息失败")
	}

	ctx.String(http.StatusOK, "修改成功")
}

func (h *UserHandler) getUidBySession(ctx *gin.Context) (int64, error) {
	sess := sessions.Default(ctx)
	userId := sess.Get("userId")
	uid, ok := userId.(int64)
	if !ok {
		return int64(0), errors.New("获取Uid失败")
	}

	return uid, nil
}
