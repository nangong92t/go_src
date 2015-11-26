/**
 * 封装baidu map api
 *
 * 此类目前包含：
 * 创建描点 baiduMapLib.init(element).createMarker({position:[]},callback);第二个参数是回调函数,可不传值
 *
 *
 */
function baiduMapLib()
{
    this.map = null;
    this.zoom = 12;
    this.centerPosition = [104.06,30.67];

    this.markers            = {};

    /**
     * 创建map
     *
     * @param string element Id html element
     * @param array  latLngArr 经纬度[]  可以不传值
     * @param int zoom  可以不传值
     *
     * @return this
     */
    this.createMap = function(element) {
        if (this.map) {
            return this.map;
        }
        this.element = element;
        this.centerPosition = arguments[1] || this.centerPosition;
        this.zoom = arguments[2] || this.zoom;

        this.map = new BMap.Map(this.element);

        this.map.centerAndZoom(
            new BMap.Point(this.centerPosition[0], this.centerPosition[1]),
            this.zoom
        );

        this.map.enableScrollWheelZoom(); // 启动鼠标滚轮缩放地图
        this.map.addControl(new BMap.NavigationControl()); // 地图平移缩放控件，PC端默认位于地图左上方，它包含控制地图的平移和缩放的功能。
        this.map.addControl(new BMap.ScaleControl()); // 比例尺控件，默认位于地图左下方，显示地图的比例关系。
        this.map.addControl(new BMap.OverviewMapControl()); // 缩略地图控件，默认位于地图右下方，是一个可折叠的缩略地图。
        this.map.addControl(new BMap.MapTypeControl()); // 地图类型控件，默认位于地图右上方。

        return this.map;
    }

    /**
     * @param object element 展示地图html element
     * @param array  latLngArr 经纬度[]  可以不传值
     * @param int zoome 默认为16 可以不传值
     *
     * @return this
     */
    this.init = function(element) {
        this.centerPosition = arguments[1] || this.centerPosition;
        this.zoom = arguments[2] || this.zoom;

        this.map = this.createMap(element, this.centerPosition);

        return this;
    }

    /**
     *
     *
     *
     */
    this.testHappen = function() {
        this.map.addControl(new BMap.NavigationControl());
        var mapObj = this;
        this.map.addEventListener('click', function(event) {
            mapObj.createMarker(
                {position: [event.point.lng, event.point.lat]},
                function(marker) {
                    marker.addEventListener(
                        'rightclick',
                        function() {
                            mapObj.map.removeOverlay(this);
                        }    
                    );
                }
            );
        });
        return mapObj;
    }
     /**
     * 创建描点
     *
     * @param object options markOptions {position: [经度, 纬度]}
     * @param function callback 回调函数 
     *
     * @example createMarker({position: [123,324]})或者createMarker({position: []}, function(marker) {});
     * 
     * @return marker
     */
    this.createMarker = function(options) {
        var point = new BMap.Point(options.position[0], options.position[1]);
        delete options.position;

        var marker = new BMap.Marker(point, options);

        var callback = arguments[1] || null;
        if(typeof callback == 'function') {
            callback(marker);
        }
        this.map.addOverlay(marker);
        return marker;
    }

     /**
     * 创建描点
     *
     * @param float lat 纬度
     * @param float lng 经度
     *
     * @return point object
     */
    this.createLatLng = function(lat, lng) {
        return new BMap.Point(lat, lng);
    };

     /**
     * 绘制线
     *
     * @param array flightPlanCoordinates point object object array
     *
     * @return polyline object
     */
    this.flightPath = function(flightPlanCoordinates) {

        var polyline = new BMap.Polyline(flightPlanCoordinates, {strokeColor:"blue", strokeWeight:2, strokeOpacity:0.5});
        this.map.addOverlay(polyline);

        return polyline;
    };

    /**
     * 绑定事物
     *
     * @param string event 事物名
     * @param function callback 函数名
     *
     */
    this.addListener = function(event, callback) {
        return this.map.addEventListener(event, callback);
    };

    /**
     * 解除绑定事物
     *
     * @param string event 事物名
     * @param function callback 函数名
     *
     */
    this.clearListener = function(event, callback) {
        return this.map.removeEventListener(event, callback);
    };

    /**
     * 创建infowindow
     * @param object options 
     *
     */
    this.createInfoWindow = function(options) {
        if (typeof options === 'undefined') {
            options = {};
        }

        var opts            = {
            width : options.width || 250,       // 信息窗口宽度
            height: options.height || 100,      // 信息窗口高度
            title : options.title || ""         // 信息窗口标题
        };
        var content         = options.content || ""  // 信息
        var infoWindow      = new BMap.InfoWindow(content, opts);  // 创建信息窗口对象

        var obj             = options.obj || this.map;

        return obj.openInfoWindow(infoWindow);      // 打开信息窗口
    };

    /**
     * 关闭infowindow
     *
     */
    this.closeInfoWindow = function(options) {
        if (typeof options === 'undefined') {
            options = {};
        }
        var obj = options.obj || this.map;

        return obj.closeInfoWindow();      // 关闭信息窗口
    };

    /**
     * 创建路径
     * @param object options 
     *
     * @return driving;
     */
    this.calcRoute = function(options) {
        if (typeof options === 'undefined') {
            options = {};
        }

        var driving = new BMap.DrivingRoute(this.map, {
            renderOptions: {
                map: this.map
            }
        });

        var begin = options.begin || new BMap.Point(113.331429, 23.14678);
        var end = options.end || new BMap.Point(113.337717, 23.130961);

        driving.search(begin, end);

        return driving;
    };

    /**
     * 获取位置
     * @param pt point
     *
     * @return location;
     */
    this.getLocation = function(pt, next) {
        var gc = new BMap.Geocoder();

        gc.getLocation(pt, function(rs){
            var location = '';
            if (rs) {
                var addComp = rs.addressComponents;
                location = addComp.province + ", " + addComp.city + ", " + addComp.district + ", " + addComp.street + ", " + addComp.streetNumber;
            }
            next && next(location);
        });
    };

    /**
     * 显示在线汽车
     */
    this.showOnlineVehicle  = function(options) {
        var owner = this;
        if (typeof options === 'undefined') {
            options = {};
        }

        var markerArr       = this.markers;
        var curCarId        = options.vehicleId;
        var Point           = {};

        var marker          = typeof markerArr[curCarId] ? markerArr[curCarId] : null;
        if (!marker) {
            marker = this.createMarker( options );
            markerArr[ curCarId ]   = marker;
            Point = marker.getPosition();
        } else {
            Point = this.createLatLng(options.position[0], options.position[1]);

            marker.setPosition(Point);
            this.map.removeOverlay(marker.getLabel());
        }
        this.getLocation(Point, function(location){
            var label = new BMap.Label(location, {"offset":new BMap.Size(9,-15)});
            marker.setLabel(label);
        })

        return marker;
    };

    /**
     * 删除在线汽车图标
     */
    this.removeOnlineVehicle  = function(options) {
        var owner = this;
        if (typeof options === 'undefined') {
            options = {};
        }

        var markerArr       = this.markers;
        var curCarId        = options.vehicleId;

        var marker          = typeof markerArr[curCarId] ? markerArr[curCarId] : null;

        if (marker) {
            marker.remove();
            this.closeInfoWindow({ obj : marker });
            delete(markerArr[curCarId]);
        }
    };
}
