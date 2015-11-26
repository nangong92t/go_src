// GENERATED CODE - DO NOT EDIT
package main

import (
	"flag"
	"reflect"
	"github.com/revel/revel"
	controllers0 "github.com/revel/revel/modules/static/app/controllers"
	_ "std/data-service/std/app"
	controllers "std/data-service/std/app/controllers"
)

var (
	runMode    *string = flag.String("runMode", "", "Run mode.")
	port       *int    = flag.Int("port", 0, "By default, read from app.conf")
	importPath *string = flag.String("importPath", "", "Go Import Path for the app.")
	srcPath    *string = flag.String("srcPath", "", "Path to the source root.")

	// So compiler won't complain if the generated code doesn't reference reflect package...
	_ = reflect.Invalid
)

func main() {
	flag.Parse()
	revel.Init(*runMode, *importPath, *srcPath)
	revel.INFO.Println("Running revel server")
	
	revel.RegisterController((*controllers.AbstractController)(nil),
		[]*revel.MethodType{
			&revel.MethodType{
				Name: "Begin",
				Args: []*revel.MethodArg{ 
				},
				RenderArgNames: map[int][]string{ 
				},
			},
			&revel.MethodType{
				Name: "Commit",
				Args: []*revel.MethodArg{ 
				},
				RenderArgNames: map[int][]string{ 
				},
			},
			&revel.MethodType{
				Name: "Rollback",
				Args: []*revel.MethodArg{ 
				},
				RenderArgNames: map[int][]string{ 
				},
			},
			&revel.MethodType{
				Name: "RenderLogin",
				Args: []*revel.MethodArg{ 
				},
				RenderArgNames: map[int][]string{ 
				},
			},
			&revel.MethodType{
				Name: "RenderNoAdmin",
				Args: []*revel.MethodArg{ 
				},
				RenderArgNames: map[int][]string{ 
				},
			},
			
		})
	
	revel.RegisterController((*controllers0.Static)(nil),
		[]*revel.MethodType{
			&revel.MethodType{
				Name: "Serve",
				Args: []*revel.MethodArg{ 
					&revel.MethodArg{Name: "prefix", Type: reflect.TypeOf((*string)(nil)) },
					&revel.MethodArg{Name: "filepath", Type: reflect.TypeOf((*string)(nil)) },
				},
				RenderArgNames: map[int][]string{ 
				},
			},
			&revel.MethodType{
				Name: "ServeModule",
				Args: []*revel.MethodArg{ 
					&revel.MethodArg{Name: "moduleName", Type: reflect.TypeOf((*string)(nil)) },
					&revel.MethodArg{Name: "prefix", Type: reflect.TypeOf((*string)(nil)) },
					&revel.MethodArg{Name: "filepath", Type: reflect.TypeOf((*string)(nil)) },
				},
				RenderArgNames: map[int][]string{ 
				},
			},
			
		})
	
	revel.RegisterController((*controllers.User)(nil),
		[]*revel.MethodType{
			&revel.MethodType{
				Name: "GetUsers",
				Args: []*revel.MethodArg{ 
					&revel.MethodArg{Name: "userId", Type: reflect.TypeOf((*int)(nil)) },
				},
				RenderArgNames: map[int][]string{ 
				},
			},
			&revel.MethodType{
				Name: "GetUserByName",
				Args: []*revel.MethodArg{ 
					&revel.MethodArg{Name: "name", Type: reflect.TypeOf((*string)(nil)) },
				},
				RenderArgNames: map[int][]string{ 
					27: []string{ 
						"user",
						"nil",
					},
				},
			},
			&revel.MethodType{
				Name: "Setting",
				Args: []*revel.MethodArg{ 
					&revel.MethodArg{Name: "sid", Type: reflect.TypeOf((*string)(nil)) },
					&revel.MethodArg{Name: "gender", Type: reflect.TypeOf((*int)(nil)) },
					&revel.MethodArg{Name: "age", Type: reflect.TypeOf((*int)(nil)) },
					&revel.MethodArg{Name: "lllness", Type: reflect.TypeOf((*int)(nil)) },
				},
				RenderArgNames: map[int][]string{ 
					39: []string{ 
						"res",
						"err",
					},
				},
			},
			&revel.MethodType{
				Name: "GetTopicList",
				Args: []*revel.MethodArg{ 
					&revel.MethodArg{Name: "uid", Type: reflect.TypeOf((*int64)(nil)) },
					&revel.MethodArg{Name: "page", Type: reflect.TypeOf((*int)(nil)) },
					&revel.MethodArg{Name: "limit", Type: reflect.TypeOf((*int)(nil)) },
				},
				RenderArgNames: map[int][]string{ 
					51: []string{ 
						"topics",
						"err",
					},
				},
			},
			&revel.MethodType{
				Name: "GetSubscribeTopicList",
				Args: []*revel.MethodArg{ 
					&revel.MethodArg{Name: "uid", Type: reflect.TypeOf((*int64)(nil)) },
					&revel.MethodArg{Name: "page", Type: reflect.TypeOf((*int)(nil)) },
					&revel.MethodArg{Name: "limit", Type: reflect.TypeOf((*int)(nil)) },
				},
				RenderArgNames: map[int][]string{ 
					63: []string{ 
						"topics",
						"err",
					},
				},
			},
			&revel.MethodType{
				Name: "Block",
				Args: []*revel.MethodArg{ 
					&revel.MethodArg{Name: "sid", Type: reflect.TypeOf((*string)(nil)) },
					&revel.MethodArg{Name: "blocked", Type: reflect.TypeOf((*int64)(nil)) },
				},
				RenderArgNames: map[int][]string{ 
					75: []string{ 
						"true",
						"err",
					},
				},
			},
			&revel.MethodType{
				Name: "GetBlockedUsers",
				Args: []*revel.MethodArg{ 
					&revel.MethodArg{Name: "sid", Type: reflect.TypeOf((*string)(nil)) },
					&revel.MethodArg{Name: "page", Type: reflect.TypeOf((*int)(nil)) },
					&revel.MethodArg{Name: "limit", Type: reflect.TypeOf((*int)(nil)) },
				},
				RenderArgNames: map[int][]string{ 
					87: []string{ 
						"blockedUsers",
						"err",
					},
				},
			},
			
		})
	
	revel.RegisterController((*controllers.Account)(nil),
		[]*revel.MethodType{
			&revel.MethodType{
				Name: "Register",
				Args: []*revel.MethodArg{ 
					&revel.MethodArg{Name: "username", Type: reflect.TypeOf((*string)(nil)) },
					&revel.MethodArg{Name: "password", Type: reflect.TypeOf((*string)(nil)) },
				},
				RenderArgNames: map[int][]string{ 
					18: []string{ 
						"user",
						"err",
					},
				},
			},
			&revel.MethodType{
				Name: "Login",
				Args: []*revel.MethodArg{ 
					&revel.MethodArg{Name: "username", Type: reflect.TypeOf((*string)(nil)) },
					&revel.MethodArg{Name: "password", Type: reflect.TypeOf((*string)(nil)) },
				},
				RenderArgNames: map[int][]string{ 
					26: []string{ 
						"nil",
					},
					34: []string{ 
						"response",
						"nil",
					},
					38: []string{ 
						"nil",
						"err",
					},
				},
			},
			&revel.MethodType{
				Name: "Logout",
				Args: []*revel.MethodArg{ 
					&revel.MethodArg{Name: "sid", Type: reflect.TypeOf((*string)(nil)) },
				},
				RenderArgNames: map[int][]string{ 
					49: []string{ 
						"res",
						"nil",
					},
				},
			},
			&revel.MethodType{
				Name: "ChangePassword",
				Args: []*revel.MethodArg{ 
					&revel.MethodArg{Name: "sid", Type: reflect.TypeOf((*string)(nil)) },
					&revel.MethodArg{Name: "password", Type: reflect.TypeOf((*string)(nil)) },
				},
				RenderArgNames: map[int][]string{ 
					59: []string{ 
						"true",
						"err",
					},
				},
			},
			
		})
	
	revel.RegisterController((*controllers.Admin)(nil),
		[]*revel.MethodType{
			&revel.MethodType{
				Name: "GetUsers",
				Args: []*revel.MethodArg{ 
					&revel.MethodArg{Name: "userId", Type: reflect.TypeOf((*int)(nil)) },
				},
				RenderArgNames: map[int][]string{ 
					23: []string{ 
						"users",
						"nil",
					},
				},
			},
			&revel.MethodType{
				Name: "GetAllStat",
				Args: []*revel.MethodArg{ 
					&revel.MethodArg{Name: "sid", Type: reflect.TypeOf((*string)(nil)) },
					&revel.MethodArg{Name: "date", Type: reflect.TypeOf((*string)(nil)) },
				},
				RenderArgNames: map[int][]string{ 
					33: []string{ 
						"nil",
					},
					37: []string{ 
						"nil",
						"err",
					},
					45: []string{ 
						"nil",
						"err",
					},
					49: []string{ 
						"nil",
						"err",
					},
					53: []string{ 
						"nil",
						"err",
					},
					57: []string{ 
						"nil",
						"err",
					},
					61: []string{ 
						"nil",
						"err",
					},
					65: []string{ 
						"nil",
						"err",
					},
					81: []string{ 
						"response",
						"nil",
					},
				},
			},
			&revel.MethodType{
				Name: "GetUserList",
				Args: []*revel.MethodArg{ 
					&revel.MethodArg{Name: "sid", Type: reflect.TypeOf((*string)(nil)) },
					&revel.MethodArg{Name: "page", Type: reflect.TypeOf((*int)(nil)) },
					&revel.MethodArg{Name: "limit", Type: reflect.TypeOf((*int)(nil)) },
					&revel.MethodArg{Name: "utype", Type: reflect.TypeOf((*string)(nil)) },
				},
				RenderArgNames: map[int][]string{ 
					101: []string{ 
						"nil",
					},
					104: []string{ 
						"list",
						"err",
					},
				},
			},
			&revel.MethodType{
				Name: "RemoveUsers",
				Args: []*revel.MethodArg{ 
					&revel.MethodArg{Name: "sid", Type: reflect.TypeOf((*string)(nil)) },
					&revel.MethodArg{Name: "ids", Type: reflect.TypeOf((*string)(nil)) },
					&revel.MethodArg{Name: "withPost", Type: reflect.TypeOf((*bool)(nil)) },
				},
				RenderArgNames: map[int][]string{ 
					118: []string{ 
						"res",
						"err",
					},
				},
			},
			&revel.MethodType{
				Name: "RemoveTopics",
				Args: []*revel.MethodArg{ 
					&revel.MethodArg{Name: "sid", Type: reflect.TypeOf((*string)(nil)) },
					&revel.MethodArg{Name: "ids", Type: reflect.TypeOf((*string)(nil)) },
				},
				RenderArgNames: map[int][]string{ 
					130: []string{ 
						"res",
						"err",
					},
				},
			},
			&revel.MethodType{
				Name: "EditTopic",
				Args: []*revel.MethodArg{ 
					&revel.MethodArg{Name: "sid", Type: reflect.TypeOf((*string)(nil)) },
					&revel.MethodArg{Name: "tid", Type: reflect.TypeOf((*int64)(nil)) },
					&revel.MethodArg{Name: "content", Type: reflect.TypeOf((*string)(nil)) },
					&revel.MethodArg{Name: "keywords", Type: reflect.TypeOf((*string)(nil)) },
				},
				RenderArgNames: map[int][]string{ 
					141: []string{ 
						"res",
						"err",
					},
				},
			},
			
		})
	
	revel.RegisterController((*controllers.Topic)(nil),
		[]*revel.MethodType{
			&revel.MethodType{
				Name: "AddNew",
				Args: []*revel.MethodArg{ 
					&revel.MethodArg{Name: "sid", Type: reflect.TypeOf((*string)(nil)) },
					&revel.MethodArg{Name: "content", Type: reflect.TypeOf((*string)(nil)) },
					&revel.MethodArg{Name: "keywords", Type: reflect.TypeOf((*string)(nil)) },
					&revel.MethodArg{Name: "bg", Type: reflect.TypeOf((*[]byte)(nil)) },
				},
				RenderArgNames: map[int][]string{ 
					28: []string{ 
						"nil",
						"err",
					},
					36: []string{ 
						"topic",
						"err",
					},
				},
			},
			&revel.MethodType{
				Name: "GetTopicList",
				Args: []*revel.MethodArg{ 
					&revel.MethodArg{Name: "sid", Type: reflect.TypeOf((*string)(nil)) },
					&revel.MethodArg{Name: "page", Type: reflect.TypeOf((*int)(nil)) },
					&revel.MethodArg{Name: "order", Type: reflect.TypeOf((*string)(nil)) },
					&revel.MethodArg{Name: "limit", Type: reflect.TypeOf((*int)(nil)) },
					&revel.MethodArg{Name: "maxid", Type: reflect.TypeOf((*int)(nil)) },
					&revel.MethodArg{Name: "needTotal", Type: reflect.TypeOf((*bool)(nil)) },
				},
				RenderArgNames: map[int][]string{ 
					58: []string{ 
						"topics",
						"err",
					},
				},
			},
			&revel.MethodType{
				Name: "GetDetail",
				Args: []*revel.MethodArg{ 
					&revel.MethodArg{Name: "id", Type: reflect.TypeOf((*int64)(nil)) },
				},
				RenderArgNames: map[int][]string{ 
					67: []string{ 
						"topic",
						"nil",
					},
				},
			},
			&revel.MethodType{
				Name: "Subscribe",
				Args: []*revel.MethodArg{ 
					&revel.MethodArg{Name: "sid", Type: reflect.TypeOf((*string)(nil)) },
					&revel.MethodArg{Name: "id", Type: reflect.TypeOf((*int64)(nil)) },
				},
				RenderArgNames: map[int][]string{ 
					78: []string{ 
						"nil",
					},
					85: []string{ 
						"res",
						"err",
					},
				},
			},
			&revel.MethodType{
				Name: "UnSubscribe",
				Args: []*revel.MethodArg{ 
					&revel.MethodArg{Name: "sid", Type: reflect.TypeOf((*string)(nil)) },
					&revel.MethodArg{Name: "id", Type: reflect.TypeOf((*int64)(nil)) },
				},
				RenderArgNames: map[int][]string{ 
					96: []string{ 
						"nil",
					},
					103: []string{ 
						"res",
						"err",
					},
				},
			},
			&revel.MethodType{
				Name: "Remove",
				Args: []*revel.MethodArg{ 
					&revel.MethodArg{Name: "sid", Type: reflect.TypeOf((*string)(nil)) },
					&revel.MethodArg{Name: "id", Type: reflect.TypeOf((*int64)(nil)) },
				},
				RenderArgNames: map[int][]string{ 
					117: []string{ 
						"res",
						"err",
					},
				},
			},
			&revel.MethodType{
				Name: "Like",
				Args: []*revel.MethodArg{ 
					&revel.MethodArg{Name: "sid", Type: reflect.TypeOf((*string)(nil)) },
					&revel.MethodArg{Name: "id", Type: reflect.TypeOf((*int64)(nil)) },
				},
				RenderArgNames: map[int][]string{ 
					129: []string{ 
						"res",
						"err",
					},
				},
			},
			&revel.MethodType{
				Name: "UnLike",
				Args: []*revel.MethodArg{ 
					&revel.MethodArg{Name: "sid", Type: reflect.TypeOf((*string)(nil)) },
					&revel.MethodArg{Name: "id", Type: reflect.TypeOf((*int64)(nil)) },
				},
				RenderArgNames: map[int][]string{ 
					141: []string{ 
						"res",
						"err",
					},
				},
			},
			&revel.MethodType{
				Name: "Comment",
				Args: []*revel.MethodArg{ 
					&revel.MethodArg{Name: "sid", Type: reflect.TypeOf((*string)(nil)) },
					&revel.MethodArg{Name: "id", Type: reflect.TypeOf((*int64)(nil)) },
					&revel.MethodArg{Name: "comment", Type: reflect.TypeOf((*string)(nil)) },
					&revel.MethodArg{Name: "replyto", Type: reflect.TypeOf((*int64)(nil)) },
				},
				RenderArgNames: map[int][]string{ 
					152: []string{ 
						"nil",
					},
					156: []string{ 
						"nil",
					},
					169: []string{ 
						"err",
					},
				},
			},
			&revel.MethodType{
				Name: "GetCommentList",
				Args: []*revel.MethodArg{ 
					&revel.MethodArg{Name: "id", Type: reflect.TypeOf((*int64)(nil)) },
					&revel.MethodArg{Name: "page", Type: reflect.TypeOf((*int)(nil)) },
					&revel.MethodArg{Name: "limit", Type: reflect.TypeOf((*int)(nil)) },
				},
				RenderArgNames: map[int][]string{ 
					177: []string{ 
						"list",
						"err",
					},
				},
			},
			
		})
	
	revel.RegisterController((*controllers.Label)(nil),
		[]*revel.MethodType{
			&revel.MethodType{
				Name: "GetList",
				Args: []*revel.MethodArg{ 
					&revel.MethodArg{Name: "page", Type: reflect.TypeOf((*int)(nil)) },
					&revel.MethodArg{Name: "limit", Type: reflect.TypeOf((*int)(nil)) },
				},
				RenderArgNames: map[int][]string{ 
					17: []string{ 
						"list",
						"err",
					},
				},
			},
			&revel.MethodType{
				Name: "Like",
				Args: []*revel.MethodArg{ 
					&revel.MethodArg{Name: "sid", Type: reflect.TypeOf((*string)(nil)) },
					&revel.MethodArg{Name: "ids", Type: reflect.TypeOf((*string)(nil)) },
				},
				RenderArgNames: map[int][]string{ 
					30: []string{ 
						"isOk",
						"nil",
					},
				},
			},
			
		})
	
	revel.RegisterController((*controllers.Setting)(nil),
		[]*revel.MethodType{
			&revel.MethodType{
				Name: "Notification",
				Args: []*revel.MethodArg{ 
					&revel.MethodArg{Name: "sid", Type: reflect.TypeOf((*string)(nil)) },
					&revel.MethodArg{Name: "mine", Type: reflect.TypeOf((*int)(nil)) },
					&revel.MethodArg{Name: "other", Type: reflect.TypeOf((*int)(nil)) },
				},
				RenderArgNames: map[int][]string{ 
					19: []string{ 
						"true",
						"nil",
					},
				},
			},
			
		})
	
	revel.DefaultValidationKeys = map[string]map[int]string{ 
		"std/data-service/std/app/models.(*Attach).Validate": { 
			22: "m.Size",
			23: "m.Extention",
			24: "m.Name",
			26: "m.SavePath",
			27: "m.SaveName",
		},
		"std/data-service/std/app/models.(*Comment).Validate": { 
			19: "m.TopicId",
			20: "m.Creator",
			21: "m.Content",
		},
		"std/data-service/std/app/models.(*Keyword).Validate": { 
			14: "m.Keyword",
		},
		"std/data-service/std/app/models.(*Label).Validate": { 
			14: "m.Name",
		},
		"std/data-service/std/app/models.(*LabelHasKeyword).Validate": { 
			14: "m.LabelId",
			15: "m.KeywordId",
		},
		"std/data-service/std/app/models.(*Like).Validate": { 
			21: "m.UserId",
			22: "m.Type",
			23: "m.TypeId",
		},
		"std/data-service/std/app/models.(*Notification).Validate": { 
			15: "m.NoticeTo",
			16: "m.Message",
		},
		"std/data-service/std/app/models.(*Profile).Validate": { 
			40: "m.Id",
			41: "m.Gender",
			42: "m.Age",
			43: "m.Lllness",
		},
		"std/data-service/std/app/models.(*Session).Validate": { 
			24: "m.Session",
		},
		"std/data-service/std/app/models.(*Subscribe).Validate": { 
			21: "m.UserId",
			22: "m.Type",
			23: "m.TypeId",
			24: "m.Created",
		},
		"std/data-service/std/app/models.(*Topic).Validate": { 
			23: "m.UserId",
			24: "m.Content",
		},
		"std/data-service/std/app/models.(*TopicHasKeyword).Validate": { 
			13: "m.TopicId",
			14: "m.KeywordId",
		},
		"std/data-service/std/app/models.(*User).Validate": { 
			32: "m.Username",
			33: "m.Password",
		},
		"std/data-service/std/app/models.(*UserBlocked).Validate": { 
			15: "m.UserId",
			16: "m.Blocked",
			17: "m.Created",
		},
		"std/data-service/std/app/models.(*UserLikeKeywordRate).Validate": { 
			14: "m.UserId",
			15: "m.KeywordId",
			16: "m.Rate",
		},
		"std/data-service/std/app/models.(*UserLoadedMaxTopic).Validate": { 
			13: "m.UserId",
			14: "m.MaxTopicId",
		},
	}
	revel.TestSuites = []interface{}{ 
	}

	revel.Run(*port)
}
