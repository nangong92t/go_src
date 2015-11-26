package runtime

import (
    "reflect"
    "../apps/std"
)

func ProcessSTD(className, method string, params interface{}) (interface{}, error) {
    return std.controller.(className).(method)(params)
}
