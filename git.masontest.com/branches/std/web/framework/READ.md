1. How to cross-domain request ?
    set ningx header configure:
        add_header Access-Control-Allow-Origin *;
        add_header Access-Control-Allow-Method GET,POST,PUT,DELETE,OPTIONS;

