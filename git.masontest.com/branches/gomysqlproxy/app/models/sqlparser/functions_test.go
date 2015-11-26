// Copyright 2014, Successfulmatch Inc. All rights reserved.
// Author TonyXu<tonycbcd@gmail.com>, 
// Build on dev-0.0.1
// MIT Licensed

// parse the sql submodel.

package sqlparser

import (
    "fmt"
    "testing"
)

func Test_ParseInsertPart(t *testing.T) {
    testData    := []string{
        "insert into test   (id, name) values(1, ?)",
        "insert into test values('', 'test')",
        "insert into test select a, b, c from test1",
        "insert into TEST     (id, name)     select a, b, c from test1 a, test2 b",
    }

    for _, insert := range testData {
        cols, vals, other, typ  := ParseInsertPart(insert, "test")
        fmt.Printf("%s:\n%#v, %#v, %#v, %#v\n\n", insert, typ, cols, vals, other)
    }
}

func Test_ParseAlterPart(t *testing.T) {
    alterSql    :=  [][]string {
        {"ALTER TABLE t2 ADD  INDEX (d), ADD UNIQUE (a)", "t2"},
        {"create   index   idx on tbl(id)", "tbl"},
        {"   ALTER table tbl add name varchar(20) not null", "tbl", "name"},
        {"alter table tbl drop column name", "tbl", "name"},
        {"ALTER TABLE document MODIFY COLUMN document_id INT auto_increment", "document", "document_id"},
        {"ALTER TABLE document MODIFY document_id INT auto_increment", "document", "document_id"},
        {"alter table tbl change name name1 varchar (20) not null", "tbl", "name"},
        {"ALTER TABLE t2 MODIFY a TINYINT NOT NULL, CHANGE b c CHAR(20)", "t2", "a"},
        {"ALTER TABLE t2 ADD c INT UNSIGNED NOT NULL AUTO_INCREMENT, ADD PRIMARY KEY (c)", "t2"},
        {"alter table tbl modify name vahrchar(20) not null", "tbl"},
        {"alter table   table1 add id int unsigned not Null auto_increment primary key", "table1"},
    }

    for _, alter    := range alterSql {
        cols    := ParseAlterPart(alter[0], alter[1])
        fmt.Printf("%#v\n", cols)
    }

}
