/**
 * 封装google map api
 *
 * 此类目前包含：
 * 创建描点 googleMapLib.init(element).createMarker({position:''},callback);第二个参数是回调函数,可不传值
 *
 *
 */
function googleMapLib()
{
    this.map = null;
    this.directionService = null;
    this.directionDisplay = null;

    this.infoWindows        = {};
    this.markers            = {};
    this.onlineCarMarker    = null;

    //html element
    this.element = null;
    this.pannel = null;
    this.zoom = 4;
    //默认中心位置重庆
    this.centerPosition = [41.87, 102.95];

    // 绘制线条路径
    this.flightPlanCoordinates  = [];

    // 线条对象
    this.polyline   = null;
    this.lastPolylineMarker = null;
    
    this.request = {
        zoom: this.zoom,
        center: {},
        zoomControl: true,
        panControl: true,
        scaleControl: true,
        mapTypeId: google.maps.MapTypeId.ROADMAP
    };
    //默认是否计算
    this._defaultCaled = true;

    this._origin = '';
    this._destination = '';
    this._clickNum = 0;
    this._wayPoints = [];

    this._pi = 3.14159265358979324;
    this._a = 6378245.0;
    this._ee = 0.00669342162296594323;
    
    /**
     * 创建map
     *
     * @param object element html element
     * @param array  latLngArr 经纬度[] 
     *
     * @return this
     */
    this.createMap = function(element) {
        if (this.map) {
            return this.map;
        }

        this.element = element;
        this.centerPosition = arguments[1] || [41.87, 102.95];

        this.request.center = new google.maps.LatLng(this.centerPosition[0], this.centerPosition[1]);
        return this.map = new google.maps.Map(
            this.element,        
            this.request
        );
    };

    this._initExpand = function() {
        google.maps.LatLng.prototype.distanceFrom = function(latlng) {
            var lat = [this.lat(), latlng.lat()]
            var lng = [this.lng(), latlng.lng()] //var R = 6371; // km (change this constant to get miles)
            var R = 6378137; // In meters
            var dLat = (lat[1] - lat[0]) * Math.PI / 180;
            var dLng = (lng[1] - lng[0]) * Math.PI / 180;
            var a = Math.sin(dLat / 2) * Math.sin(dLat / 2) + Math.cos(lat[0] * Math.PI / 180) * Math.cos(lat[1] * Math.PI / 180) * Math.sin(dLng / 2) * Math.sin(dLng / 2);
            var c = 2 * Math.atan2(Math.sqrt(a), Math.sqrt(1 - a));
            var d = R * c;
            return Math.round(d);
        }

        google.maps.Marker.prototype.distanceFrom = function(marker) {
            return this.getPosition().distanceFrom(marker.getPosition());
        }

        google.maps.Polyline.prototype.getLength = function() {
            var d = 0;
            var path = this.getPath();
            var latlng;
            for (var i = 0; i < path.getLength() - 1; i++) {
                latlng = [path.getAt(i), path.getAt(i + 1)];
                d += latlng[0].distanceFrom(latlng[1]);
            }
            return d;
        }
    };

    /**
     * 创建DirectionsService
     *
     * @return google.maps.DirectionsService;
     */
    this.createDirectionService = function() {
        if (this.directionService) {
            return this.directionService;
        }
        return this.directionService = new google.maps.DirectionsService();
    };

    /**
     * 创建展示DirectionsRender
     * 
     * @param object obj 默认为空,可以不传
     *
     * @return google.maps.DirectionsRenderer
     */
    this.createDirectionDisplay = function() {
        if (this.directionDisplay) {
            return this.directionDisplay;
        }
        var options = arguments[0] || {draggable: true};
        return this.directionDisplay = new google.maps.DirectionsRenderer(options);
    };

    /**
     * @param object element 展示地图html element
     * @param object options 初始化参数
     *
     * @return this
     */
    this.init = function(element, options) {
        if (typeof options === 'undefined') {
            options = {};
        }

        this.centerPosition = options.center || this.centerPosition;
        this.pannel = options.pannel || this.pannel;
        this._initExpand();

        this.map = this.createMap(element, this.centerPosition);
        this.directionDisplay = this.createDirectionDisplay();
        this.directionDisplay.setMap(this.map);
        if (!this.pannel) {
            this.directionDisplay.setPanel(this.pannel);
        }

        return this;
    };

    /**
     * 测试触发事物
     *
     *
     */
    this.testHappen = function() {
        var googleMapObj = this;
        var isCaled = this._defaultCaled;
        
        this.addListener(this.map, 'click', function(event){
            var marker = googleMapObj.createMarker({position: event.latLng});
            googleMapObj.addListener(marker, 'rightclick', function(event){
                this.setMap(null);
            });
            googleMapObj._wayPoints = [];

            if (isCaled) {
                //默认计算动作，今后要修改的
                ++googleMapObj._clickNum%2 ? (googleMapObj._origin = marker) : (googleMapObj._destination = marker);
                if (googleMapObj._origin && googleMapObj._destination) {
                    googleMapObj._origin.setTitle('起点');
                    googleMapObj._destination.setTitle('终点');
                    googleMapObj.calcRoute({
                        origin: googleMapObj._origin.getPosition(), 
                        destination: googleMapObj._destination.getPosition(),
                    });
                }
            }
        });

        this.addListener(this.map, 'rightclick', function(event) {
            googleMapObj.createMarker({position: event.latLng});
            googleMapObj._wayPoints.push({location: event.latLng, stopover: true});
            googleMapObj.calcRoute({
                origin: googleMapObj._origin.getPosition(), 
                destination: googleMapObj._destination.getPosition(),
            });
        });
    };

    /**
     * 绑定事物
     *
     * @param object 监听对象
     * @param string event 事物名
     * @param function callback 函数名
     *
     * @return void;
     */
    this.addListener = function(obj, event, callback) {
        return google.maps.event.addListener(obj, event, callback);
    };

    /**
     * 解除绑定事物
     *
     * @param object 被监听对象
     * @param string event 事物名
     * @param function callback 函数名
     *
     * @return void;
     */
    this.clearListeners = function(obj, event, callback) {
        return google.maps.event.clearListeners(obj, event, callback);
    };

    /**
     * 创建描点
     *
     * @param object options markOptions
     * @param function callback 回调函数 
     *
     * @example createMarker({position: ''})或者createMarker({position: ''}, function(marker) {});
     *
     * @return marker
     */
    this.createMarker = function(options) {
        options.map = options.map || this.map;
        options.icon = options.icon || 'img/car-stoping.png';
        if (options.noicon) { options.icon = ''; }

        var marker = new google.maps.Marker(options);
        
        var callback = arguments[1] || null;
        var infoWindow = arguments[2] || null;

        if(typeof callback == 'function') {
            callback(marker, infoWindow);
        } else if(infoWindow) {
            var self = this;
            if (typeof infoWindow == 'object') {
                this.addListener(marker, callback, function() {
                    infoWindow.open(self.map, marker);
                });
            } else {
                this.addListener(marker, callback, function() {
                    self.createInfoWindow({content: infoWindow.toString()}).open(self.map, marker);
                });
            }
        } 
        return marker;
    };

    /**
     * @param object requestOptions DirectionsRequest
     *
     * @return this;
     */
    this.calcRoute = function(requestOptions) {
        var travelMode =  requestOptions.travelMode || google.maps.TravelMode.DRIVING;
        var alternatives = requestOptions.provideRouteAlternatives || false;
        requestOptions.travelMode = travelMode;
        requestOptions.provideRouteAlternatives = alternatives;

        requestOptions.waypoints = this._wayPoints;

        this.directionService = this.createDirectionService();

        var googleMapObj = this;
        var callback = arguments[1] || null;
        this.directionService.route(requestOptions, function(result, status) {
            if (status == google.maps.DirectionsStatus.OK) {
                googleMapObj.directionDisplay.setDirections(result);
            }       
            if (typeof callback == 'function') {
                callback(result);
            }
        });

        return this;
    };

    /**
     * 清除路径
     *
     */
    this.clearRoute = function(controlDiv) {
        var self = this;
        if (self.directionDisplay) {
            self.directionDisplay.setMap(null);
            self.directionDisplay = null;
            self.directionDisplay = self.createDirectionDisplay();
            self.directionDisplay.setMap(self.map);
        }
    };


    this.flightPath = function(flightPlanCoordinates) {

        var flightPath = new google.maps.Polyline({
            path: flightPlanCoordinates,
            geodesic: true,
            strokeColor: '#FF0000',
            strokeOpacity: 1.0,
            strokeWeight: 2
        });

        flightPath.setMap(this.map);
        return flightPath;
    };

    /**
     * 创建infowindow
     * @param object options 
     *
     */
    this.createInfoWindow = function(options) {
        return new google.maps.InfoWindow(options);
    };

    /**
     * GPS偏差纠正
     *
     */
    this.transform = function(wgLat, wgLon) {
        var a = this._a;
        var pi = this._pi;
        var ee = this._ee;

        var latlng = new Array();
        if (this._outOfChina(wgLat, wgLon)) {  
            latlng[0] = wgLat;  
            latlng[1] = wgLon;  
            return latlng;  
        }  
        var dLat = this._transformLat(wgLon - 105.0, wgLat - 35.0);  
        var dLon = this._transformLon(wgLon - 105.0, wgLat - 35.0);  
        var radLat = wgLat / 180.0 * pi;  
        var magic = Math.sin(radLat);  
        magic = 1 - ee * magic * magic;  
        var sqrtMagic = Math.sqrt(magic);  
        dLat = (dLat * 180.0) / ((a * (1 - ee)) / (magic * sqrtMagic) * pi);  
        dLon = (dLon * 180.0) / (a / sqrtMagic * Math.cos(radLat) * pi);  
        latlng[0] = wgLat + dLat;  
        latlng[1] = wgLon + dLon;   

        return latlng;
    };

    this._outOfChina = function(lat, lon) {
        if (lon < 72.004 || lon > 137.8347)  
            return true;  
        if (lat < 0.8293 || lat > 55.8271)  
            return true;  
        return false;  
    };

    this._transformLat = function(x, y) {
        var a = this._a;
        var pi = this._pi;
        var ee = this._ee;

        var ret = -100.0 + 2.0 * x + 3.0 * y + 0.2 * y * y + 0.1 * x * y + 0.2 * Math.sqrt(Math.abs(x));  
        ret += (20.0 * Math.sin(6.0 * x * pi) + 20.0 * Math.sin(2.0 * x * pi)) * 2.0 / 3.0;  
        ret += (20.0 * Math.sin(y * pi) + 40.0 * Math.sin(y / 3.0 * pi)) * 2.0 / 3.0;  
        ret += (160.0 * Math.sin(y / 12.0 * pi) + 320 * Math.sin(y * pi / 30.0)) * 2.0 / 3.0;  
        return ret;  
    };

    this._transformLon = function(x, y) {
        var a = this._a;
        var pi = this._pi;
        var ee = this._ee;
 
        var ret = 300.0 + x + 2.0 * y + 0.1 * x * x + 0.1 * x * y + 0.1 * Math.sqrt(Math.abs(x));  
        ret += (20.0 * Math.sin(6.0 * x * pi) + 20.0 * Math.sin(2.0 * x * pi)) * 2.0 / 3.0;  
        ret += (20.0 * Math.sin(x * pi) + 40.0 * Math.sin(x / 3.0 * pi)) * 2.0 / 3.0;  
        ret += (150.0 * Math.sin(x / 12.0 * pi) + 300.0 * Math.sin(x / 30.0 * pi)) * 2.0 / 3.0;  
        return ret;  
    };

    this.createLatLng = function(lat, lon) {
        var _arr = this.transform(lat, lon);
        return new google.maps.LatLng(_arr[0], _arr[1]);
    };

    this.closeAllInfoWindow = function() {
        for (var carId in this.infoWindows) {
            this.infoWindows[carId].close();
        }
    };

    /**
     * 显示在线汽车列表与动作监听
     */
    this.toShow  = function(car, onCheck, isPush, drivings, isChecked, isSetCenter, options) {
        if (typeof options === 'undefined') {
            options = {};
        }

        var infoWindowArr   = this.infoWindows;
        var markerArr       = this.markers;
        var curCarId        = car.VehicleID;
        if (options.isOnline) {
            this.map.setZoom(15);

            if (this.onlineCarMarker) {
                this.onlineCarMarker.setMap(null);
            }
        }

        var marker = typeof markerArr[curCarId] ? markerArr[curCarId] : null;
        if (onCheck && isChecked === 0) { delete(markerArr[curCarId]); }

        var infoWindow = typeof infoWindowArr[curCarId] !== 'undefined' ? infoWindowArr[curCarId] : null;
        if (onCheck && isChecked === 0) { delete(infoWindowArr[curCarId]); }

        var gps = car["LastPosition"],
            lat = gps["Latitude"],
            lon = gps["Longitude"];

        /*** 或取显示icon类型 ***/
        var isOpenAcc   = car['LastPosition'].AccState; 

        var latLng = this.createLatLng(lat, lon);
        
        var firstFlag = !markerArr.length ? true : false;

        if (!marker) {
            var mOption = {position: latLng , title: car['LicensePlateNo']};
            if (isOpenAcc) { mOption['icon'] = 'img/car-moving.png'; }
            marker = this.createMarker( mOption );
            if (options.isOnline) {
                this.onlineCarMarker = marker;
            }
            markerArr[ curCarId ]   = marker;
        } else {
            marker.setPosition(latLng);
            marker.setMap(this.map);
            marker.setIcon(null);
            if (isOpenAcc) {
                marker.setIcon('img/car-moving.png');
            } else {
                marker.setIcon('img/car-stoping.png');
            }
        }

        if (onCheck && isChecked === 0) {
            marker.setMap(null);
            return false;
        }

        var infoWindowHtml = car.LicensePlateNo+'-'+car.LicensePlateColor+'<br/>';
        infoWindowHtml += '里程:'+car.LastPosition.Mileage+'公里<br/>';
        infoWindowHtml += '时间:'+car.LastPosition.GpsTime+'<br/>';
        infoWindowHtml += '车速:'+car.LastPosition.Speed+'公里/小时<br/>';
        if (typeof car.LastPosition['location'] !== 'undefined') {
            infoWindowHtml += '位置:'+car.LastPosition['location']+'';
        }

        //infoWindowHtml += '经纬度:['+car.LastPosition.Latitude+','+car.LastPosition.Longitude+']';
        
        if (!infoWindow) {
            infoWindow = this.createInfoWindow({content:infoWindowHtml});
            infoWindowArr[ curCarId ]   = infoWindow;
        } else {
            infoWindow.setContent(infoWindowHtml);
        }

        this.addListener(marker, 'click', function() {
            infoWindow.open(this.map, marker); 
        });


        //在地图可视区域显示所有描点
        if (!isPush) { 
            if (isSetCenter || !onCheck || _.keys(markerArr).length<=1) {
                this.map.setCenter(latLng);
                if (isSetCenter) {
                    this.closeAllInfoWindow();
                    infoWindow.open(this.map, marker);
                }
            }

            var bounds = new google.maps.LatLngBounds();
            for (var i in this.markers) {
                bounds.extend(this.markers[i].getPosition());
            }
            if (_.keys(this.markers).length <= 1 || !onCheck) {
                this.map.setZoom(15);
            } else if (!isSetCenter || (onCheck && isChecked===1)) {
                this.map.fitBounds(bounds);
            }
        }

        //是否轨迹监听
        if (typeof drivings !== 'undefined') {
            for(var i=0; i<drivings.length; i++) {
                if (car.VehicleID == drivings[i].VehicleID) {
                    if (!drivings[i].driving) {
                        drivings[i].driving = this.flightPath([latLng]);
                    } else {
                        var pathArr = drivings[i].driving.getPath().getArray();
                        pathArr.push(latLng);
                        drivings[i].driving.setPath(pathArr);
                    }
                    break;
                }
            }
        }
        return marker;
    };

    this.drawPolyline   = function(latlng) {
        var owner   = this;
        this.flightPlanCoordinates.push(latlng); 

        if (!this.polyline) {
            this.polyline = this.flightPath(this.flightPlanCoordinates);
        } else {
            this.polyline.setPath(this.flightPlanCoordinates);
        }
        this.polyline.setVisible(true);

        this.calcRoute({
            origin: this.flightPlanCoordinates[0],
            destination: this.flightPlanCoordinates[this.flightPlanCoordinates.length-1]
        },function(ret) {
            var legs = ret['routes'][0]['legs'];
            var zuiyou = 0;
            for(var i=0;i<legs.length; i++) {
                zuiyou += legs[i]['distance']['value'];
            }
            
            //计算
            var mym = list[list.length-1]['Mileage']-list[0]['Mileage'];
            var _htmlStr = "<p>历史轨迹："+mym.toFixed(3)+"公里</p>";
            _htmlStr += "<p>最优路径："+ (zuiyou/1000).toFixed(4)+"公里</p>";
           // _htmlStr += "差值里程："+((mym*1000-zuiyou)/1000).toFixed(4)+"公里</p>";

            if (!owner.lastPolylineMarker) {
                owner.lastPolylineMarker = this.createMarker(
                    {position: owner.flightPlanCoordinates[owner.flightPlanCoordinates.length-1]},
                    'click',
                    _htmlStr
                );
            } else {
                owner.lastPolylineMarker.setPosition(owner.flightPlanCoordinates[owner.flightPlanCoordinates.length-1]);
            }
            owner.lastPolylineMarker.setVisible(true);
        });

        this.directionDisplay.setMap(this.map);

        this.map.setCenter(this.flightPlanCoordinates.pop());
    }; 

}
