package modules

import (
    "std/data-service/std/app/models"
    "github.com/coopernurse/gorp"
    "time"
    "fmt"
    "strings"
)

type Keyword struct {
    Db *gorp.DbMap
}

func NewKeyword(db *gorp.DbMap) *Keyword {
    return &Keyword {
        Db: db,
    }
}

// 保存未过滤过的keyword字符串
func (c *Keyword) SaveString(keys string) []models.Keyword {
    keywords    := strings.Split(keys, " ")
    keyCnt      := len(keywords)
    newKeywords := make([]models.Keyword, 0)

    for i:=0; i<keyCnt; i++ {
        // 过滤多余的.
        if keywords[i] == "" {
            continue
        }
        keyword, _ := c.AddNew(keywords[i])
        newKeywords = append(newKeywords, *keyword)
    }

    return newKeywords
}

func (c *Keyword) AddNew(keyName string) (*models.Keyword, error) {
    isExists, _ := c.GetByName(keyName)

    if isExists.Id != 0 {
        return isExists, nil
    }

    keyword := &models.Keyword{
        Id: 0,
        Keyword: keyName,
        Created: time.Now().Unix(),
    }

    err := c.Db.Insert(keyword)

    return keyword, err
}

func (c *Keyword) GetByName(keyName string) (*models.Keyword, error) {
    var keyword models.Keyword

    err := c.Db.SelectOne(&keyword, "select * from keyword where keyword = ?", keyName)

    return &keyword, err
}

// 更改某用户对某些keywords的喜欢程度.
func (m *Keyword) ChangeUserLikeRate(uid int64, keywords []models.Keyword, way string, rate int) bool {
    var likeRates []models.UserLikeKeywordRate

    keyIds  := []interface{}{}
    for _, one  := range keywords {
        keyIds  = append(keyIds, one.Id)
    }

    // 查找到现有的.
    sql     := "select * from user_like_keyword_rate where user_id="+fmt.Sprintf("%d", uid) + " and keyword_id in (" + strings.TrimRight(strings.Repeat("?,", len(keyIds)), ",") + ")"
    m.Db.Select(&likeRates, sql, keyIds...)

    // 生成现有喜欢数据map。
    likedMap := map[int64]models.UserLikeKeywordRate{}
    for _, one := range likeRates {
        likedMap[ one.KeywordId ] = one
    }

    // 正式增加数据.
    for _, one := range keywords {
        fmt.Printf("%#v, %#v\n", one, likedMap[ one.Id ].UserId)
        if likedMap[ one.Id ].UserId != 0  {          // 已存在数据，更新
            oneLike := likedMap[ one.Id ]
            curRate := 0
            switch way {
                case "+":
                    curRate = oneLike.Rate + rate
                case "-":
                    curRate = oneLike.Rate - rate
            }
            m.Db.Exec("update user_like_keyword_rate set rate = ? where user_id=? and keyword_id=?", curRate, oneLike.UserId, oneLike.KeywordId)
        } else {                                // 数据不存在，插入
            switch way {
                case "+":
                case "-":
                    rate = 0 - rate
            }

            newLike := &models.UserLikeKeywordRate{
                UserId: uid,
                KeywordId: one.Id,
                Rate: rate,
            }
            m.Db.Insert(newLike)
        }
    }

    return true
}


