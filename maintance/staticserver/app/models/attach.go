// Copyright 2014, SuccessfumMatch Teq Inc. All rights reserved.
// Author TonyXu<tonycbcd@gmail.com>,
// Build on dev-0.0.1
// MIT Licensed

// 文件附件数据模型

package models

import (
    "errors"
    "time"
    "labix.org/v2/mgo/bson"
)

type AttachModel struct {
    AbstractModel               `bson:"-", json:"-"`
    Id          bson.ObjectId   `json:"id" bson:"_id"`
    Name        string          `json:"name"`
    Ext         string          `json:"ext"`
    Size        int             `json:"size"`
    Hash        string          `json:"hash"`
    SaveHost    string          `json:"save_host" bson:"save_host"`
    SavePath    string          `json:"save_path" bson:"save_path"`
    SaveName    string          `json:"save_name" bson:"save_name"`
    Created     time.Time       `json:"created"`
}

// the singleton pattern
var oneAttach *AttachModel

func NewAttach() *AttachModel {
    if oneAttach == nil {
        oneAttach  = &AttachModel{}
        oneAttach.Init(oneAttach)
    }

    return oneAttach
}

// override parent function
func (m *AttachModel) GetDbName() string {
    return "static"
}

// override parent function
func (m *AttachModel) GetCollectionName() string {
    return "attaches"
}

// override parent function
func (m *AttachModel) BeforeSave() (error, bool) {
    return nil, true
}

// override parent function
func (m *AttachModel) AfterSave() {
}

// override parent function
func (m *AttachModel) Validate() error {
    if !m.Id.Valid() { m.SetDefault() }
    if m.Size == 0 { return errors.New("20108") }
    return nil
}

// override parent function
func (m *AttachModel) SetDefault() {
    m.Id        = m.GetID()
    m.Created   = time.Now()
    m.IsNew     = true
}

func (m *AttachModel) GetId() bson.ObjectId { return m.Id }
func (m *AttachModel) BeforeFind() {}
func (m *AttachModel) AfterFind() {}
func (m *AttachModel) BeforeUpdate() (error, bool) { return nil, true }
func (m *AttachModel) AfterUpdate() {}
func (m *AttachModel) BeforeDelete() {}
func (m *AttachModel) AfterDelete() {}

