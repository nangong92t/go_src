package controllers

import (
    "std/data-service/std/app/modules"
    "errors"
	r "github.com/revel/revel"
)

type Account struct {
	AbstractController
}

func (c Account) Register(username string, password string) r.Result {
    userModule  := modules.NewUser(c.Db)

    user, err   := userModule.AddNew(username, password)

    return c.Render(user, err)
}

func (c Account) Login(username string, password string) r.Result {
    userModule  := modules.NewUser(c.Db)
    user        := userModule.GetUserByName(username)

    var err error
    if (user.IsDel == 1) { return c.Render(nil, errors.New("Sorry, this account has been deleted!")) }

    if (user.Password == password) {
        response    := map[string]interface{}{
            "session": c.Sess.Build(&user),
            "user": &user,
        }

	    return c.Render(response, nil)
    } else {
        err     = errors.New("Sorry, username or password is wrong.")

	    return c.Render(nil, err)
    }
}

func (c Account) Logout(sid string) r.Result {
    session     := c.Sess.GetSession(sid)
    if session.Id == 0 {
        return c.RenderLogin()
    }

    res := c.Sess.Remove(sid)
	return c.Render(res, nil)
}

func (c Account) ChangePassword(sid string, password string) r.Result {
    session     := c.Sess.GetSession(sid)
    if session.Id == 0 {
        return c.RenderLogin()
    }

    err         := modules.NewUser(c.Db).ChangePassword(session.SessionVal.UserId, password)
	return c.Render(true, err)
}
