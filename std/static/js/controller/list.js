/*
 * Admin abstruct list js Controller
 * v0.0.1-dev for Mason STD Project
 * build by TonyXu<tonycbcd@gmail.com> on 2014-10-07
 * MIT Licensed
 */
;define(function(require, exports, module){
    var LabelBuilder    = Backbone.View.extend({
        tagName: 'div',
        parentObj: {},
        template: require('view/label_builder.html'),
        trTmplate: require('view/label_item.html'),
        apis: {
            add: "label/add"
        },
        events: {
            'click #add-btn': 'toAdd'
        },
        initialize: function(options) {
            this.parentObj  = options['parentObj'];    
        },
        render: function() {
            $(this.el).html( this.template );
            return this;
        },
        toAdd: function() {
            var label   = $('#new-label').val();
            if (label == "") {
                utils.alert("Please input the label name");
                return false;
            }

            var owner   = this;
            window.Store.getData(this.apis.add, [label], function(result) {
                result['created']   = moment.unix(Number(result['created'])).format('YYYY-MM-DD'); 
                result['author']    = 'tony'; 
                $('#label-container').prepend(_.template(owner.trTmplate, result)); 
                utils.closeFloatDialog();        
            });
        }
    });
    var UnwantedWordsBuilder = Backbone.View.extend({
        tagName: 'div',
        parentObj: {},
        template: require('view/UnwantedWords_builder.html'),
        trTmplate: require('view/UnwantedWords_item.html'),
        apis: {
            add: "admin/AddUnwantWord"
        },
        events: {
            'click #add-btn': 'toAdd'
        },
        initialize: function(options) {
            this.parentObj  = options['parentObj'];
        },
        render: function() {
            $(this.el).html( this.template );
            return this;
        },
        toAdd: function() {
            var unwantedwords   = $('#new-unwantedwords').val();
            if (unwantedwords == "") {
                utils.alert("Please input the unwanted words");
                return false;
            }

            var owner   = this;
            window.Store.getData(this.apis.add, [unwantedwords], function(result) {
                //result['created']   = moment.unix(Number(result['created'])).format('YYYY-MM-DD'); 
                //result['creator']    = 'tony'; 
                //$('#unwanted_words-container').prepend(_.template(owner.trTmplate, result)); 
                //utils.closeFloatDialog();
                window.location.reload();
            });
        }
    });

    var Main = Backbone.View.extend({
        options: {},
        type: '',
        components: {},
        container: {},
        apis: {
            deleteTopic: "admin/removetopics"
        },
        bindAction: function(){
            var owner   = this;
            $('#select_all').bind('click', function(){
                $('input[name="'+owner.type+'_id[]"]').attr("checked",this.checked); 
            });
            $('#delete-button').bind('click', function() { owner.deleteItem(); });
            $('#delete-ban-button').bind('click', function() { owner.deleteItem(); });
            $('#hide-button').bind('click', function() { owner.hideItem(); });
            $('#add-button').bind('click', function() { owner.addItem(); });
            $('#add-image-button-container').bind('click', function() { owner.addItem(); });
        },
        initialize: function(options){
            this.apis   = options.apis;
            this.type   = options.type;
            this.bindAction();
            this.initUploader();
        },
        getSelectedItem: function() {
            var items = [];
            $('input[name="'+this.type+'_id[]"]').each(function() {
                var curObj  = $(this);
                if (curObj.attr("checked")) {
                    items.push(curObj.val());
                }
            });
            return items;
        },
        deleteItem: function() {
            var owner   = this;
            var items  = this.getSelectedItem();
            if (items.length == 0) { return false; }

            if (!confirm('Are you sure to delete this '+this.type+'?')) { return false; }
            window.Store.getData(this.apis.delete, [items.join(",")], function(result) {
                if (result) {
                    _.each(items, function(one) { $("#"+owner.type+"-"+one).remove(); });
                }
            });
        },
        hideItem: function() {

        },
        addItem: function() {
            switch (this.type) {
                case 'label':
                    this.components['builder']    = new LabelBuilder({parentObj:this});
                    break;
                case 'unwanted_words':
                    this.components['builder'] = new UnwantedWordsBuilder({parentObj:this});
                    break;
                default:
                    utils.alert("Sorry, no find this type builder class");
                    return false;
            }

            utils.openFloatDialog( this.components['builder'].render().el, {
              title: 'add new '+ this.type,
              modal: true,
              autoOpen: true,
              resizable: false,
              autoResize: true,
              position: {my:'center', at:'center'},
              buttons: {}
            });
        },
        initUploader: function() {
            var owner   = this;
            require.async('lib/single-uploader', function(loader) {
                loader.singleuploader({
                    container: '#add-image-button-container',
                    afterHandle: function(attach) { window.location = window.location; }
                });
            });
        }
    });
    exports.Init = Main;
});
