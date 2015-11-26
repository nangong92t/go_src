// Copyright 2014, SuccessfumMatch Teq Inc. All rights reserved.
// Author TonyXu<tonycbcd@gmail.com>,
// Build on dev-0.0.1
// MIT Licensed

// for attatch module

package modules

import (
    //"github.com/revel/revel/cache"
    "labix.org/v2/mgo/bson"
    "maintance/staticserver/app/models"
)

type Attach struct {
    AbstractModule
}

// the singleton pattern
var oneAttachIntance *Attach

func NewAttach() *Attach {
    if oneAttachIntance == nil {
        oneAttachIntance  = &Attach{}
    }

    oneAttachIntance.Init()

    return oneAttachIntance
}

func (m *Attach) GetByIds(ids []bson.ObjectId) []models.AttachModel {
    // TODO cache data.

    attaches := []models.AttachModel{}
    models.NewAttach().FindAll(models.M{"_id": models.M{"$in": ids}}, "", 0, 0, false, &attaches)

    return attaches
}

func (m *Attach) GetMapByIds(ids []bson.ObjectId) map[bson.ObjectId]*models.AttachModel {
    // TODO cache data.

    curMap      := map[bson.ObjectId]*models.AttachModel{}
    attaches    := m.GetByIds(ids)
    for i:=0; i<len(attaches); i++ {
        one := attaches[i]
        if _, isOk := curMap[ one.Id ]; !isOk {
            curMap[ one.Id ]  = &one
        }
    }

    return curMap
}
