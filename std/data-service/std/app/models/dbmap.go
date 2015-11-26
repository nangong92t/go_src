package models

import (
    "github.com/coopernurse/gorp"
)

func InitDbTableMap(dbmap *gorp.DbMap) {
    dbmap.AddTableWithName(User{}, "user").SetKeys(true, "Id")
    dbmap.AddTableWithName(Comment{}, "comment").SetKeys(true, "Id")
    dbmap.AddTableWithName(Keyword{}, "keyword").SetKeys(true, "Id")
    dbmap.AddTableWithName(Label{}, "label").SetKeys(true, "Id")
    dbmap.AddTableWithName(LabelHasKeyword{}, "label_has_keyword").SetKeys(true, "Id")
    dbmap.AddTableWithName(Notification{}, "notification").SetKeys(true, "Id")
    dbmap.AddTableWithName(Profile{}, "profile").SetKeys(false, "Id")
    dbmap.AddTableWithName(Topic{}, "topic").SetKeys(true, "Id")
    dbmap.AddTableWithName(TopicHasKeyword{}, "topic_has_keyword")
    dbmap.AddTableWithName(UserLikeKeywordRate{}, "user_like_keyword_rate")
    dbmap.AddTableWithName(UserLoadedMaxTopic{}, "user_loaded_max_topic")
    dbmap.AddTableWithName(Session{}, "user_session")
    dbmap.AddTableWithName(Attach{}, "attach").SetKeys(true, "Id")
    dbmap.AddTableWithName(UserBlocked{}, "user_blocked").SetKeys(true, "Id")
    dbmap.AddTableWithName(UnwantWord{}, "unwant_word").SetKeys(true, "Id")
}
