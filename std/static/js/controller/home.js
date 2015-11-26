/*
 * Admin Home js Controller
 * v0.0.1-dev for Mason STD Project
 * build by TonyXu<tonycbcd@gmail.com> on 2014-09-19 
 * MIT Licensed
 */
;define(function(require, exports, module) {
    require('lib/highcharts');

    var HighCharts  = Backbone.View.extend({
        container: 0,
        mId: 'hightcharts-container',
        tagName: 'div',
        initialize: function(options) {
            this.container  = options['container'];
            this.container.html( this.render().el );
            this.initHightchats();
        },
        render: function() {
            $(this.el)
                .html('')
                .attr('id', this.mId)
                .attr('style', 'min-width:1100px;height:400px');

            return this;
        },
        initHightchats: function() {
            Highcharts.setOptions({
                global: {
                    useUTC: false
                }
            });

            $('#'+this.mId).highcharts({
                chart: {
                    type: 'spline'
                },
                title: {
                    text: 'STD Data Change Status'
                },
                subtitle: {
                    text: ''
                },
                xAxis: {
                    type: 'datetime',
                    dateTimeLabelFormats: { 
                        second: '%m/%d',
                        minute: '%m/%d',
                        hour: '%m/%d',
                        day: '%m/%d',
                        week: '%m/%d',
                        month: '%m/%d',
                        year: '%Y'
                    }
                },
                yAxis: {
                    title: {
                        text: 'the amount of data.'
                    },
                    min: 0
                },
                /*
                tooltip: {
                    formatter: function() {
                        return '<p style="color:'+this.series.color+';font-weight:bold;">'
                         + this.series.name + 
                         '</p><br /><p style="color:'+this.series.color+';font-weight:bold;">Time: ' + Highcharts.dateFormat('%m/%d', this.x) + 
                         '</p><br /><p style="color:'+this.series.color+';font-weight:bold;">Amount: '+ this.y + '</p>';
                    }
                },
                credits: {
                    enabled: false,
                    text: "stdapp.masontest.com",
                    href: "http://stdapp.masontest.com"
                },*/
                series: [
                    {
                        name: 'User Total Change',
                        data: this.model.dateVal.users,
                        lineWidth: 2,
                        
                    },
                    {
                        name: 'Topic Total Change',
                        data: this.model.dateVal.topics,
                        lineWidth: 2,
                        color : '#9C0D0D'
                    },
                    {
                        name: 'Comment Total Change',
                        data: this.model.dateVal.comments,
                        lineWidth: 2,
                        color : '#C88E54'
                    }            
                ]
            });
        }      
    });

    exports.Init    = Backbone.View.extend({
        dateVal: {},
        initialize: function(options) {
            this.dateVal    = options['date'];
            this.showHightcharts();
        },
        showHightcharts: function() {
            new HighCharts({container: $('#change-status-map'), model: {dateVal: this.dateVal}})
        }
    });
});
