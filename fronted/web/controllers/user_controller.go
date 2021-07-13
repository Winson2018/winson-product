package controllers

import (
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/iris/v12/sessions"
	"strconv"
	"winson-product/datamodels"
	"winson-product/encrypt"
	"winson-product/services"
	"winson-product/tool"
)

type UserController struct {
	Ctx iris.Context
	Service services.IUserService
	Session *sessions.Session
}

func (c *UserController) GetRegister() mvc.View{
	return mvc.View{
		Name: "user/register.html",
	}
}

func (c *UserController) PostRegister(){
	var(
		nickName = c.Ctx.FormValue("nickName")
		userName = c.Ctx.FormValue("userName")
		password = c.Ctx.FormValue("password")
	)

	user := &datamodels.User{
		UserName: userName,
		NickName: nickName,
		HashPassword: password,
	}

	_, err := c.Service.AddUser(user)
	c.Ctx.Application().Logger().Debug(err)
	if err != nil{
		c.Ctx.Redirect("/user/error")
		return
	}
	c.Ctx.Redirect("/user/login")
	return
}

func (c *UserController) GetLogin() mvc.View {
	return mvc.View{
		Name: "user/login.html",
	}
}

func (c *UserController) PostLogin() mvc.Response{
	var (
		userName = c.Ctx.FormValue("userName")
		password = c.Ctx.FormValue("password")
	)

	user, isOk := c.Service.IsPwdSuccess(userName, password)
	if !isOk {

		return mvc.Response{
			Path: "/user/login",
		}
	}

	//写用户ID到cookie中
	tool.GlobalCookie(c.Ctx, "uid", strconv.FormatInt(user.ID, 10))
	uidByte := []byte(strconv.FormatInt(user.ID, 10))
	uidString, err := encrypt.EnPwdCode(uidByte)
	if err != nil {
		fmt.Println(err)
	}
	//写入用户浏览器
	tool.GlobalCookie(c.Ctx, "sign", uidString)
	//c.Session.Set("userID", strconv.FormatInt(user.ID, 10))

	return mvc.Response{
		Path: "/product/",
	}
}