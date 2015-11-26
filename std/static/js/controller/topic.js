/*
 * Admin Home js Controller
 * v0.0.1-dev for Mason STD Project
 * build by Sam<samsheng133@gmail.com> on 2014-09-29
 * MIT Licensed
 */
;
define(function(require, exports, module){

    var Main = Backbone.View.extend({
        options: {},
        components: {},
        container: {},
        apis: {
            deleteTopic: "admin/removetopics"
        },
        bindAction: function(){
            var owner   = this;
            $('#select_all').bind('click', function(){
                $('input[name="topic_id[]"]').attr("checked",this.checked); 
            });
            $('#delete-button').bind('click', function() { owner.deleteTopic(); });
            $('#delete-ban-button').bind('click', function() { owner.deleteTopic(); });
            $('#hide-button').bind('click', function() { owner.hideTopic(); });
        },
        initialize: function(options){
            this.bindAction();
        },
        getSelectedTopic: function() {
            var topics = [];
            $('input[name="topic_id[]"]').each(function() {
                var curObj  = $(this);
                if (curObj.attr("checked")) {
                    topics.push(curObj.val());
                }
            });
            return topics;
        },
        deleteTopic: function() {
            topics  = this.getSelectedTopic();
            if (topics.length == 0) { return false; }

            if (!confirm('Are you sure to delete this topic?')) { return false; }
            window.Store.getData(this.apis.deleteTopic, [topics.join(",")], function(result) {
                if (result) {
                    _.each(topics, function(one) { $("#topic-"+one).remove(); });
                }
            });
        },
        hideTopic: function() {

        }
    });
    exports.Init = Main;
});
