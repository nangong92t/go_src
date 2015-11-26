package modules

import (
    "fmt"
    "time"
    "strings"
    "errors"
    "github.com/coopernurse/gorp"
    "std/data-service/std/app/models"
    "std/data-service/std/app/helpers"
)

type UnwantWord struct {
    Db *gorp.DbMap
}

func NewUnwantWord(db *gorp.DbMap) *UnwantWord {
    return &UnwantWord {
        Db: db,
    }
}

func (m *UnwantWord) GetOneByWord(word string) *models.UnwantWord {
    unwant := &models.UnwantWord{}
    m.Db.SelectOne(unwant, "select * from unwant_word where word=? limit 1", word)
    return unwant
}

func (m *UnwantWord) Add(creator int64, word string) error {
    if isExists := m.GetOneByWord(word); isExists.Id != 0 {
        return errors.New("Sorry, this unwwant word is exists!")
    }

    newOne  := &models.UnwantWord {
        Word: word,
        Used: 0,
        Creator: creator,
        Created: time.Now().Unix(),
    }

    m.Db.Insert(newOne)

    return nil
}

func (m *UnwantWord) GetList(page, limit int, needTotal bool) (*Pager, error) {
    var unwant []models.UnwantWord

    tableName   := "unwant_word"

    sql     := "select * from " + tableName

    if page == 0 { page = 1 }
    if limit == 0 { limit = 20 }

    offset  := (page - 1) * limit

    sql     = sql + " limit " + fmt.Sprintf("%d",offset) + ", " + fmt.Sprintf("%d", limit)

    _, err  := m.Db.Select(&unwant, sql)

    userIds := []string{}
    unwantCnt := len(unwant)

    // Map 相关子id
    for i:=0; i<unwantCnt; i++ {
        userId  := fmt.Sprintf("%d", unwant[i].Creator)

        if !helpers.StringInSlice(userId, userIds) {
            userIds = append(userIds, userId)
        }
    }

    total   := int64(0)
    if needTotal {
        total, _    = m.Db.SelectInt("select count(*) from " + tableName)
    }

    // 获取相关用户数据.
    var userMap map[int64]models.User
    if len(userIds) > 0 {
        userModule      := NewUser(m.Db)
        userMap, err    = userModule.GetUserMapById(userIds)
        if err != nil { return nil, err }
    }

    // Merge数据.
    for i:=0; i<unwantCnt; i++ {
        unwant[i].Author        = userMap[ unwant[i].Creator ].Username
    }

    pager   := &Pager{
        Total: total,
        List: unwant,
    }

    return pager, err
}

func (m *UnwantWord) Remove(ids []interface{}) (bool, error) {
    valLen      := len(ids)
    condition   := fmt.Sprintf("word_id in (%s)",  strings.TrimRight(strings.Repeat("?,", valLen), ","))
    sql         := "delete from unwant_word where %s"
    _, err      := m.Db.Exec(fmt.Sprintf(sql, condition), ids...)

    return err==nil, err
}
