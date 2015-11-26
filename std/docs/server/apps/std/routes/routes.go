// GENERATED CODE - DO NOT EDIT
package routes

import "github.com/revel/revel"


type tAbstractController struct {}
var AbstractController tAbstractController


func (_ tAbstractController) Begin(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("AbstractController.Begin", args).Url
}

func (_ tAbstractController) Commit(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("AbstractController.Commit", args).Url
}

func (_ tAbstractController) Rollback(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("AbstractController.Rollback", args).Url
}

func (_ tAbstractController) RenderLogin(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("AbstractController.RenderLogin", args).Url
}

func (_ tAbstractController) RenderNoAdmin(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("AbstractController.RenderNoAdmin", args).Url
}


type tStatic struct {}
var Static tStatic


func (_ tStatic) Serve(
		prefix string,
		filepath string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "prefix", prefix)
	revel.Unbind(args, "filepath", filepath)
	return revel.MainRouter.Reverse("Static.Serve", args).Url
}

func (_ tStatic) ServeModule(
		moduleName string,
		prefix string,
		filepath string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "moduleName", moduleName)
	revel.Unbind(args, "prefix", prefix)
	revel.Unbind(args, "filepath", filepath)
	return revel.MainRouter.Reverse("Static.ServeModule", args).Url
}


type tUser struct {}
var User tUser


func (_ tUser) GetUsers(
		userId int,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "userId", userId)
	return revel.MainRouter.Reverse("User.GetUsers", args).Url
}

func (_ tUser) GetUserByName(
		name string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "name", name)
	return revel.MainRouter.Reverse("User.GetUserByName", args).Url
}

func (_ tUser) Setting(
		sid string,
		gender int,
		age int,
		lllness int,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "sid", sid)
	revel.Unbind(args, "gender", gender)
	revel.Unbind(args, "age", age)
	revel.Unbind(args, "lllness", lllness)
	return revel.MainRouter.Reverse("User.Setting", args).Url
}

func (_ tUser) GetTopicList(
		uid int64,
		page int,
		limit int,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "uid", uid)
	revel.Unbind(args, "page", page)
	revel.Unbind(args, "limit", limit)
	return revel.MainRouter.Reverse("User.GetTopicList", args).Url
}

func (_ tUser) GetSubscribeTopicList(
		uid int64,
		page int,
		limit int,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "uid", uid)
	revel.Unbind(args, "page", page)
	revel.Unbind(args, "limit", limit)
	return revel.MainRouter.Reverse("User.GetSubscribeTopicList", args).Url
}

func (_ tUser) Block(
		sid string,
		blocked int64,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "sid", sid)
	revel.Unbind(args, "blocked", blocked)
	return revel.MainRouter.Reverse("User.Block", args).Url
}

func (_ tUser) GetBlockedUsers(
		sid string,
		page int,
		limit int,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "sid", sid)
	revel.Unbind(args, "page", page)
	revel.Unbind(args, "limit", limit)
	return revel.MainRouter.Reverse("User.GetBlockedUsers", args).Url
}


type tAccount struct {}
var Account tAccount


func (_ tAccount) Register(
		username string,
		password string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "username", username)
	revel.Unbind(args, "password", password)
	return revel.MainRouter.Reverse("Account.Register", args).Url
}

func (_ tAccount) Login(
		username string,
		password string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "username", username)
	revel.Unbind(args, "password", password)
	return revel.MainRouter.Reverse("Account.Login", args).Url
}

func (_ tAccount) Logout(
		sid string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "sid", sid)
	return revel.MainRouter.Reverse("Account.Logout", args).Url
}

func (_ tAccount) ChangePassword(
		sid string,
		password string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "sid", sid)
	revel.Unbind(args, "password", password)
	return revel.MainRouter.Reverse("Account.ChangePassword", args).Url
}


type tAdmin struct {}
var Admin tAdmin


func (_ tAdmin) GetUsers(
		userId int,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "userId", userId)
	return revel.MainRouter.Reverse("Admin.GetUsers", args).Url
}

func (_ tAdmin) GetAllStat(
		sid string,
		date string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "sid", sid)
	revel.Unbind(args, "date", date)
	return revel.MainRouter.Reverse("Admin.GetAllStat", args).Url
}

func (_ tAdmin) GetUserList(
		sid string,
		page int,
		limit int,
		utype string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "sid", sid)
	revel.Unbind(args, "page", page)
	revel.Unbind(args, "limit", limit)
	revel.Unbind(args, "utype", utype)
	return revel.MainRouter.Reverse("Admin.GetUserList", args).Url
}

func (_ tAdmin) RemoveUsers(
		sid string,
		ids string,
		withPost bool,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "sid", sid)
	revel.Unbind(args, "ids", ids)
	revel.Unbind(args, "withPost", withPost)
	return revel.MainRouter.Reverse("Admin.RemoveUsers", args).Url
}

func (_ tAdmin) RemoveTopics(
		sid string,
		ids string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "sid", sid)
	revel.Unbind(args, "ids", ids)
	return revel.MainRouter.Reverse("Admin.RemoveTopics", args).Url
}

func (_ tAdmin) EditTopic(
		sid string,
		tid int64,
		content string,
		keywords string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "sid", sid)
	revel.Unbind(args, "tid", tid)
	revel.Unbind(args, "content", content)
	revel.Unbind(args, "keywords", keywords)
	return revel.MainRouter.Reverse("Admin.EditTopic", args).Url
}


type tTopic struct {}
var Topic tTopic


func (_ tTopic) AddNew(
		sid string,
		content string,
		keywords string,
		bg []byte,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "sid", sid)
	revel.Unbind(args, "content", content)
	revel.Unbind(args, "keywords", keywords)
	revel.Unbind(args, "bg", bg)
	return revel.MainRouter.Reverse("Topic.AddNew", args).Url
}

func (_ tTopic) GetTopicList(
		sid string,
		page int,
		order string,
		limit int,
		maxid int,
		needTotal bool,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "sid", sid)
	revel.Unbind(args, "page", page)
	revel.Unbind(args, "order", order)
	revel.Unbind(args, "limit", limit)
	revel.Unbind(args, "maxid", maxid)
	revel.Unbind(args, "needTotal", needTotal)
	return revel.MainRouter.Reverse("Topic.GetTopicList", args).Url
}

func (_ tTopic) GetDetail(
		id int64,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "id", id)
	return revel.MainRouter.Reverse("Topic.GetDetail", args).Url
}

func (_ tTopic) Subscribe(
		sid string,
		id int64,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "sid", sid)
	revel.Unbind(args, "id", id)
	return revel.MainRouter.Reverse("Topic.Subscribe", args).Url
}

func (_ tTopic) UnSubscribe(
		sid string,
		id int64,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "sid", sid)
	revel.Unbind(args, "id", id)
	return revel.MainRouter.Reverse("Topic.UnSubscribe", args).Url
}

func (_ tTopic) Remove(
		sid string,
		id int64,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "sid", sid)
	revel.Unbind(args, "id", id)
	return revel.MainRouter.Reverse("Topic.Remove", args).Url
}

func (_ tTopic) Like(
		sid string,
		id int64,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "sid", sid)
	revel.Unbind(args, "id", id)
	return revel.MainRouter.Reverse("Topic.Like", args).Url
}

func (_ tTopic) UnLike(
		sid string,
		id int64,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "sid", sid)
	revel.Unbind(args, "id", id)
	return revel.MainRouter.Reverse("Topic.UnLike", args).Url
}

func (_ tTopic) Comment(
		sid string,
		id int64,
		comment string,
		replyto int64,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "sid", sid)
	revel.Unbind(args, "id", id)
	revel.Unbind(args, "comment", comment)
	revel.Unbind(args, "replyto", replyto)
	return revel.MainRouter.Reverse("Topic.Comment", args).Url
}

func (_ tTopic) GetCommentList(
		id int64,
		page int,
		limit int,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "id", id)
	revel.Unbind(args, "page", page)
	revel.Unbind(args, "limit", limit)
	return revel.MainRouter.Reverse("Topic.GetCommentList", args).Url
}


type tLabel struct {}
var Label tLabel


func (_ tLabel) GetList(
		page int,
		limit int,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "page", page)
	revel.Unbind(args, "limit", limit)
	return revel.MainRouter.Reverse("Label.GetList", args).Url
}

func (_ tLabel) Like(
		sid string,
		ids string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "sid", sid)
	revel.Unbind(args, "ids", ids)
	return revel.MainRouter.Reverse("Label.Like", args).Url
}


type tSetting struct {}
var Setting tSetting


func (_ tSetting) Notification(
		sid string,
		mine int,
		other int,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "sid", sid)
	revel.Unbind(args, "mine", mine)
	revel.Unbind(args, "other", other)
	return revel.MainRouter.Reverse("Setting.Notification", args).Url
}


