;
define(function (require, exports, module) {
    var Store = {
        Cache: {},
        messages: {
            error: '抱歉，数据读取失败，已通知技术人员，我们尽快处理，请稍后再试.'
        },
        alert: function (mesg) {
            if (utils)
                utils.alert(mesg);
            else
                alert(mesg);
        },
        getData: function (url, data, next, otherData, way, isCache) {
            this.clearExpiredCache();

            var owner = this;
            !otherData && (otherData = {});
            !way && (way = 'get');

            var accessToken = typeof token !== 'undefined' && token ? token : '';
            if (accessToken) url = apiUrl + '/v1/' + accessToken + '/' + url;

            switch (way) {
                case "get":
                    if (data)
                        url += '/' + JSON.stringify(data);
                    else
                        url += '/[]';

                    break;

                case "post":
                    break;
            }

            var cacheKey = hex_md5(url + JSON.stringify(data) + JSON.stringify(otherData));

            // 因为数据随时可以改变，情况比较复杂，暂时不考虑缓存
            if (isCache && typeof this.Cache[ cacheKey ] !== 'undefined') {
                var cacheTime = this.Cache[ cacheKey ]['t'];
                var nowTime = new Date().getTime();
                if (nowTime < cacheTime) {
                    next && next(this.Cache[ cacheKey ]['d']);
                    return false;
                }
            }

            if (_.isArray(otherData)) {
                otherData = {data: otherData};
            }
            $.ajax({
                url: url,
                type: way,
                data: otherData,
                dataType: "json",
                success: function (result) {
                    if (typeof result['errcode'] !== 'undefined') {
                        owner.alert(result['errcode']);
                    } else {
                        var nowTime = new Date().getTime();
                        if (isCache && !_.isEmpty(result.recordset)) owner.Cache[ cacheKey ] = {d: result.recordset, t: isCache + nowTime};
                        next && next(result);
                    }
                },
                error: function (XMLHttpRequest, textStatus, errorThrown) {
                    if (XMLHttpRequest.responseText != '') {
                        var result = $.parseJSON(XMLHttpRequest.responseText);
                        if (typeof result['errcode'] !== 'undefined') {
                            owner.alert(result.errmsg);
                            return;
                        }
                    }
                    //alert( url );
                },
                complete: function (XMLHttpRequest, textStatus) {
                    this; // 调用本次AJAX请求时传递的options参数
                }
            });

        },

        clearExpiredCache: function () {
            var curTime = new Date().getTime();
            for (var key in this.Cache) {
                var curData = this.Cache[ key ];
                if (curTime > curData['t']) {
                    delete(this.Cache[ key ]);
                }
            }
        }

    };

    window.Store = Store;
});


