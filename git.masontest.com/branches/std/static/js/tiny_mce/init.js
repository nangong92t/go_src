tinyMCE.init({
  // General options
  mode : "textareas",
  theme : "advanced",
  plugins : "autolink,lists,pagebreak,style,layer,table,save,advhr,advimage,advlink,emotions,iespell,inlinepopups,insertdatetime,preview,media,searchreplace,print,contextmenu,paste,directionality,fullscreen,noneditable,visualchars,nonbreaking,xhtmlxtras,template,wordcount,advlist,autosave,visualblocks",
  theme_advanced_buttons1 : "save,newdocument,|,bold,italic,underline,strikethrough,|,justifyleft,justifycenter,justifyright,justifyfull,styleselect,formatselect,fontselect,fontsizeselect",
  theme_advanced_buttons2 : "bullist,numlist,|,outdent,indent,blockquote,|,undo,redo,|,link,unlink,anchor,image,cleanup,help,code,|,insertdate,inserttime,preview,|,forecolor,backcolor",
  theme_advanced_toolbar_location : "top",
  theme_advanced_toolbar_align : "left",
  theme_advanced_statusbar_location : "bottom",
  theme_advanced_resizing : true,

  style_formats : [
    {title : 'Bold text', inline : 'b'},
    {title : 'Red text', inline : 'span', styles : {color : '#ff0000'}},
    {title : 'Red header', block : 'h1', styles : {color : '#ff0000'}},
    {title : 'Example 1', inline : 'span', classes : 'example1'},
    {title : 'Example 2', inline : 'span', classes : 'example2'},
    {title : 'Table styles'},
    {title : 'Table row 1', selector : 'tr', classes : 'tablerow1'}
  ]
});
