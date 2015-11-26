package modules

import (
    "std/data-service/std/app/models"
    "std/data-service/std/app/helpers"
    "github.com/coopernurse/gorp"
    "errors"
    "time"
    "strings"
    "fmt"
    "strconv"
)

type User struct {
    Db *gorp.DbMap
}

func NewUser(db *gorp.DbMap) *User {
    return &User{
        Db: db,
    }
}

// 获取用户列表
func (m *User) GetUsers() []models.User {
    // 暂时没有加缓存.
    var users []models.User
	_, err  := m.Db.Select(&users, `select * from user`)

	if err != nil {
		panic(err)
	}

    return users
}

// 添加新用户
func (m *User) AddNew(username string, password string) (models.User, error) {
    // 查找当前用户名是否存在
    var user models.User
    var profile models.Profile
    m.Db.SelectOne(&user, "select * from user where username=?", username)

    if user.Id != 0 && user.IsDel == 0 {
        return user, errors.New("Sorry, this user has exists!")
    }

    if user.IsDel == 1 {
        user.Username   = user.Username + "_del" + strconv.FormatInt(time.Now().Unix(), 10)
        m.Db.Update(&user)
    }

    // add account
    user = models.User {
        Username: username,
        Password: password,
        Created: time.Now().Unix(),
    }
    err := m.Db.Insert(&user)
    if err != nil { return user, err }

    // add profile
    profile = models.Profile {
        Id: user.Id,
        Gender: 1,      // default the man.
        Age: 0,
        Lllness: 0,
        Avator: "",
    }
    err = m.Db.Insert(&profile)
    if err != nil { return user, err }

    return user, nil
}

// 通过用户名获取用户数据.
func (m *User) GetUserByName(username string) models.User {
    var user models.User

    m.Db.SelectOne(&user, "select * from user where username=?", username)

    return user
}

// 通过User Id 夺取对应Map数据
func (m *User) GetUserMapById(userIds []string) (map[int64]models.User, error) {
    idLen   := len(userIds)
    userMap := make(map[int64]models.User, idLen)

    var users []models.User
    _, err  := m.Db.Select(&users, "select * from user where user_id in (" + strings.Join(userIds, ", ") + ")")

    if err != nil { return nil, err }

    for _, one := range users {
        userMap[ one.Id ] = one
    }

    return userMap, nil
}

// 通过User Id获取单个用户.
func (m *User) GetById(id int64) *models.User {
    user := &models.User{}

    if id == 0 { return user }

    m.Db.SelectOne(user, "select * from user where user_id=? limit 1", id)
    return user
}

// 修改用户密码
func (m *User) ChangePassword(uid int64, password string) error {
    _, err  := m.Db.Exec("update user set password=? where user_id=?", password, uid)

    return err
}

// 拉黑某用户
func (m *User) Block(uid int64, blocked int64, isAdmin bool) error {
    if isAdmin && uid == 0  { return errors.New("sorry, just admin role can set blocked a user for All!") }

    block   := &models.UserBlocked{
        Id: 0,
        UserId: uid,
        Blocked: blocked,
        Created: time.Now().Unix(),
    }
    err := m.Db.Insert(block)

    return err
}

// 获取我的黑名单列表
func (m *User) GetBlockedUsers(uid int64, page int, limit int) ([]models.UserProfile, error) {
    var users []models.UserProfile

    _, err  := m.Db.Select(&users, "select u.user_id, u.username, u.created from user u left join user_blocked b on b.blocked=u.user_id where b.user_id = ?", uid)

    return users, err
}

// Get the user stat in a date.
func (m *User) GetStat(date int64) (map[string]int64, error) {
    result  := map[string]int64 {
        "new": 0,
        "reported": 0,
        "blocked": 0,
    }

    maxDate := date + 86400
    var err error
    result["new"], err    = m.Db.SelectInt("select count(1) from user where created >= ? and created < ? limit 1", date, maxDate)
    if err != nil { return result, err }

    result["blocked"], err = m.Db.SelectInt("select count(1) from user_blocked where created >= ? and created < ? limit 1", date, maxDate)

    if err != nil { return result, err }

    return result, nil
}

//  get user date news count.
func (m *User) GetOneMonthDateNewCnt(toDate int64) ([]interface{}, error) {
    // golang the default time is 08:00:00, so need to add to next day time
    toDate  += (24 - 8) * 3600

    days    := int64(31)
    daySeconds  := int64(86400)
    minDay  := toDate - days * daySeconds

    var res []models.DateData

    sql     := "select unix_timestamp(date_format(FROM_UNIXTIME( `created`),'%Y-%m-%d')) * 1000 as date, count(1) as val from user where created > ? and created < ? group by date"
    _, err  := m.Db.Select(&res, sql, minDay, toDate)
    if err != nil { return nil, err }

    result  := helpers.MergeDateData(res, days, toDate)

    return result, nil
}

// set user self's profile.
func (m *User) SettingProfile(uid int64, gender, age, lllness int) (bool, error) {
    // to check the profile weather is exists.
    var profile models.Profile
    err := m.Db.SelectOne(&profile, "select * from profile where user_id = ?",uid)
    if err != nil { return false, err }
    if profile.Id == 0 { return false, errors.New("Sorry no this user") }

    if gender != 0 { profile.Gender  = gender }
    if age != 0 { profile.Age     = age }
    if lllness != 0 { profile.Lllness = lllness }

    _, err  = m.Db.Update(&profile)

    return err==nil, err
}

// get user detail by uids
func (m *User) GetUserWithProfileByIds(uids []interface{}) ([]models.UserProfile, error) {
    sql := `select p.*, u.user_id, u.username, u.created from user u left join profile p on u.user_id=p.user_id where u.user_id in (%s)`
    var users []models.UserProfile
    userCnt := len(uids)
    if userCnt <= 0 { return users, nil }

    _, err  := m.Db.Select(
        &users,
        fmt.Sprintf(sql, strings.TrimRight(strings.Repeat("?,", userCnt), ",")),
        uids...
    )

    for i, one := range users {
        if one.Lllness == 0 { continue }
        one.LllnessStr   = strings.Join(helpers.RestoreMoreDataInOneColumn(int64(one.Lllness), models.LllnessTypes), ", ")
        users[i]    = one
    }

    if err != nil { return users, err }

    return users, nil
}

func (m *User) GetMostPostUserList(page, limit, maxCnt int) (Pager, error) {
    if page == 0 { page = 1}
    if limit == 0 { limit = 20 }
    if maxCnt == 0 { maxCnt = 2 }
    offset  := (page - 1) * limit

    userCnt := []struct{
        UserId int64 `db:"user_id"`
        Cnt    int64 `db:"cnt"`
    }{}

    userCntMap  := map[int64]int64{}
    userIds     := []interface{}{}

    // get post more maxCnt topic's user.
    uSql := `select %s from (
                select user_id, count(*) cnt 
                from topic 
                where is_del=0 and created >= ? 
                group by user_id
                order by created desc) a 
            where a.cnt > ? %s`
    prevDay     := time.Now().Unix() - 86400
    _, err      := m.Db.Select(&userCnt, fmt.Sprintf(uSql, "a.user_id, cnt", "limit ? offset ?"), prevDay, maxCnt, limit, offset)
    if err != nil { return Pager{}, err }

    total, err  := m.Db.SelectInt(fmt.Sprintf(uSql, "count(1)", ""), prevDay, maxCnt)
    if err != nil { return Pager{}, err }

    // map user
    for _, one := range userCnt {
        userCntMap[ one.UserId ]    = one.Cnt
        userIds = append(userIds, strconv.FormatInt(one.UserId, 10))
    }

    // get users
    users, err  := m.GetUserWithProfileByIds(userIds)
    if err != nil { return Pager{}, err }

    // merge total

    pager       := Pager{
        Total: total,
        List: users,
    }

    return pager, nil
}

func (m *User) GetBlockUserList(page, limit int, isBlocked bool) (Pager, error) {
    if page == 0 { page = 1}
    if limit == 0 { limit = 20 }
    offset  := (page - 1) * limit

    block   := []struct{
        Block int64 `db:"block"`
    }{}
    selected    := "distinct user_id"
    if isBlocked { selected = "distinct blocked" }

    userIds     := []interface{}{}

    // get user ids.
    uSql        := `select %s from user_blocked %s`
    _, err      := m.Db.Select(&block, fmt.Sprintf(uSql, selected + " block", "order by created limit ? offset ?"), limit, offset)
    if err != nil { return Pager{}, err }

    total, err  := m.Db.SelectInt(fmt.Sprintf(uSql, "count("+selected+") cnt", ""))
    if err != nil { return Pager{}, err }

    // map user
    for _, one := range block {
        userIds = append(userIds, strconv.FormatInt(one.Block, 10))
    }

    // get users
    users, err  := m.GetUserWithProfileByIds(userIds)
    if err != nil { return Pager{}, err }

    // merge total

    pager       := Pager{
        Total: total,
        List: users,
    }

    return pager, nil
}

func (m *User) RemoveUsers(uid []interface{}) (bool, error) {
    uidLen  := len(uid)
    curTime := time.Now().Unix()
    sql := "update user set is_del=1, deleted=%d where is_del=0 and user_id in (%s)"
    _, err  := m.Db.Exec(fmt.Sprintf(sql, curTime, strings.TrimRight(strings.Repeat("?,", uidLen), ",")), uid...)

    return err==nil, err
}

