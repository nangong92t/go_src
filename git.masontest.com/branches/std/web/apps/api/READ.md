1. How to cross-domain request ?
    set ningx header configure:
        add_header Access-Control-Allow-Origin *;
        add_header Access-Control-Allow-Method GET,POST,PUT,DELETE,OPTIONS;



        工作汇报：

        １. 完成Api webservice开发，实现前端与PHP后端的数据中转；如本地： http://api.ly.com/rest/asdfasdsadf/Vehicle_Main/getUserInfoByUid/%5B123,234%5D



        1.1  采用标准的REST规范： 
    对于所有 item 列表调用： GET /rest/item
对于某个 item 信息调用： GET /rest/item/(\d+)
    创建一个 item： POST /rest/item
    更新一个 item： PUT /rest/item/(\d+)
删除一个 item： DELETE /rest/item/(\d+)

    api uri规范：
    /rest/[token]/[PHPServer Handle Name]/[PHPServer Handle Action Name]/[Params Json string]

    它实际是PHPServer的前端WebService的Proxy和Authen认证。


    1.2 数据输出为标准的HTTP CODE：
    200 => 'OK',
    400 => 'Bad Request',
    401 => 'Unauthorized',
    402 => 'Payment Required',
    403 => 'Forbidden',
    404 => 'Not Found',
    500 => 'Internal Server Error',
    501 => 'Not Implemented',

    1.3 关于api认证， 为了节约时间，暂时做了很简单的token认证，oauth规范以后有精力再加了。


    Q总， 你哪儿记得把 api.hravel.com 加上.

    详细代码在： git.18.tl:branches/ly/api,  详细文挡没时间写了， 以后再补

