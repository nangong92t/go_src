/*
 * Location js 地区Modules
 * v0.0.1-dev | for ly project
 * build by tony(tonycbcd@gmail.com) on 2014-07-07
 * MIT Licensed
 */
define(function (require, exports) {

    var Location = Backbone.View.extend({
        options: 0,
        provinces: 0,
        tagName: 'div',
        template: '<select class="input-sm col-xs-6 col-md-2" id="province"></select><select class="hide input-sm col-xs-6 col-md-2" id="city"></select><select class="hide input-sm col-xs-6 col-md-2" id="area"></select>',
        events: {
            'change #province': 'provinceChanger',
            'change #city': 'cityChanger',
            'change #area': 'areaChanger'
        },
        apis: {
            'getProvinces': 'BasicService_Location/getProvinces',
            'getCitiesByProId': 'BasicService_Location/getCitiesByProId',
            'getAreaByCityId': 'BasicService_Location/getAreaByCityId'
        },
        initialize: function (options) {
            this.options = options;
        },
        render: function () {
            var owner = this;
            $(this.el).append(this.template);

            this.initData();
            return this;
        },
        initData: function () {
            var owner = this;
            this.setAreaId();
            window.Store.getData(this.apis.getProvinces, [], function (data) {
                var proHtml = (function () {
                    var html = '<option value="0">--请选择省份--</option>';
                    for (var i = 0; i < data.length; i++) {
                        var one = data[i];
                        html += '<option value="' + one['provinceID'] + '">' + one['province'] + '</option>';
                    }
                    return html;
                })();
                $('#province', owner.el).append(proHtml);
                if (owner.areaIds[0]) {
                    $('#province', owner.el).find('option[value="' + owner.areaIds[0] + '"]').attr('selected', true).change();
                }
            });
        },
        setAreaId: function () {
            this.areaIds = [];
            var areaId = this.options['defaultAraaId'];
            if (!areaId) {
                return;
            }
            var areaId = parseInt(areaId);
            var province, city, area;
            var pn = areaId % 10000;
            var cn = areaId % 100;
            province = areaId - pn;
            if (cn > 0) {
                city = areaId - cn;
                area = areaId;
            }
            this.areaIds = [province, city, area];
        },
        changeCityByProId: function (proId, next) {
            window.Store.getData(this.apis.getCitiesByProId, [proId], function (data) {
                next && next(data);
            });
        },
        changeAreaByCityId: function (cityId, next) {
            window.Store.getData(this.apis.getAreaByCityId, [cityId], function (data) {
                next && next(data);
            });
        },
        provinceChanger: function (e) {
            var owner = this;
            var curObj = $(e.currentTarget);
            var proId = $(curObj).val();

            if (typeof this.options['changedProvinceHandle'] === 'function') {
                this.options['changedProvinceHandle'](proId);
            }

            var next = function (data) {
                var options = '<option value="0">--请选择城市--</option>';
                if (typeof(data) !== 'undefined' && data.length > 0) {
                    for (var i = 0; i < data.length; i++) {
                        options += '<option value="' + data[i]['cityID'] + '">' + data[i]['city'] + '</option>';
                    }
                    $('#city', owner.el).html(options).removeClass('hide');
                    if (owner.areaIds[1]) {
                        $('#city', owner.el).find('option[value="' + owner.areaIds[1] + '"]').attr('selected', true).change();
                    }
                }
            };

            owner.changeCityByProId(proId, next);
        },
        cityChanger: function (e) {
            var owner = this;
            var curObj = $(e.currentTarget);
            var cityId = $(curObj).val();

            if (typeof this.options['changedCityHandle'] === 'function') {
                this.options['changedCityHandle'](cityId);
            }

            var next = function (data) {
                var options = '<option value="0">--请选择地区--</option>';
                if (typeof(data) !== 'undefined' && data.length > 0) {
                    for (var i = 0; i < data.length; i++) {
                        options += '<option value="' + data[i]['areaID'] + '">' + data[i]['area'] + '</option>';
                    }
                    $('#area', owner.el).html(options).removeClass('hide');
                    if (owner.areaIds[2]) {
                        $('#area', owner.el).find('option[value="' + owner.areaIds[2] + '"]').attr('selected', true).change();
                    }
                }
            };

            owner.changeAreaByCityId(cityId, next);
        },
        areaChanger: function (e) {
            var owner = this;
            var curObj = $(e.currentTarget);
            var areaId = $(curObj).val();

            if (typeof this.options['changedAreaHandle'] === 'function') {
                this.options['changedAreaHandle'](areaId);
            }
        },
        parseLocationByGoogleMap: function (data, nt) {
            var mapLeng = data.length;
            var findColumns = ['postal_code', 'country', 'administrative_area_level_1', 'locality'];
            var address = {};

            var isTypeOf = function (type, types) {
                for (var i = 0; i < types.length; i++) {
                    if (type === types[i])
                        return true;
                }
                return false;
            };

            var toParse = function (addressComponents, address) {
                var compLeng = addressComponents.length;
                for (var i = 0; i < compLeng; i++) {
                    var curComp = addressComponents[i];
                    for (var j = 0; j < findColumns.length; j++) {
                        var curColumn = findColumns[j];
                        if (typeof address[ curColumn ] !== 'undefined')
                            break;

                        if (isTypeOf(curColumn, curComp['types'])) {
                            address[ curColumn ] = curComp['long_name'];
                        }
                    }
                }

                return address;
            };

            for (var i = 0; i < mapLeng; i++) {
                if (typeof data[i].address_components !== 'undefined') {
                    address = toParse(data[i].address_components, address);
                    if (typeof address['postal_code'] !== 'undefined' && typeof address['administrative_area_level_1'] !== 'undefined' &&
                        typeof address['country'] !== 'undefined' && typeof address['locality'] !== 'undefined')
                        break;
                }
            }

            return address;
        }
    });

    exports.location = Location;
});
