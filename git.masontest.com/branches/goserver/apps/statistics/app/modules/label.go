package modules

import (
    "strings"
    "github.com/coopernurse/gorp"
    "std/data-service/std/app/models"
)

type Label struct {
    Db *gorp.DbMap
}

func NewLable(db *gorp.DbMap) *Label {
    return &Label {
        Db: db,
    }
}

// 喜欢多个Lable.
func (m *Label) Like(uid int64, ids string) bool {
    idsArr := strings.Split(ids, ",")

    // 整理label id.
    idsArr2 := []interface{}{}
    for _, lid := range idsArr {
        if lid == "" { continue }
        idsArr2 = append(idsArr2, lid)
    }
    if len(idsArr2) == 0 { return false }

    // 查找Label下有那些Keywords.
    keywords    := m.GetKeywordsInLabels(idsArr2)

    // 保存用户喜欢的Keywords关联.
    sql         := `INSERT INTO user_like_keyword_rate
        (user_id, keyword_id, rate)
        SELECT ?, ?, 100 FROM dual
        WHERE not exists (
            select * from user_like_keyword_rate a
            where a.user_id = ? and a.keyword_id=?
        )`

    for _, one := range keywords {
        m.Db.Exec(sql, uid, one.Id, uid, one.Id)
    }

    return true
}

// 获取Label下有那些Keywords.
func (m *Label) GetKeywordsInLabels(ids []interface{}) []models.Keyword {
    var keywords []models.Keyword

    m.Db.Select(&keywords, "select distinct k.keyword_id, k.keyword, k.created from keyword k left join label_has_keyword l on l.keyword_id=k.keyword_id where l.label_id in (" + strings.TrimRight(strings.Repeat("?,", len(ids)), ",") + ")", ids...)

    return keywords
}

// 获取label列表.
func (m *Label) GetList(page int, limit int) (*Pager, error) {
    if page == 0 { page = 1 }
    if limit == 0 { limit = 20 }
    offset := (page-1) * limit

    sql     := `Select * from label limit ?, ?`

    var labels []models.Label
    _, err  := m.Db.Select(&labels, sql, offset, limit)

    pager   := &Pager {
        List: labels,
    }

    return pager, err
}

