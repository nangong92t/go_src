<?php
/**
 * 初始化sphinx客户端
 */

namespace SphinxClient;


class SphinxClientApi {
    private static $connections = array ();
    private static $agentIndexs = array();

    private static $config;

    /**
     *
     * @return SphinxClient
     */
    public static function Connection($name = 'default') {
        if (! isset ( self::$connections [$name] )) {
            self::$connections[$name] = new api\SphinxPhpClient();
            self::_Setup( $name );
        }

        return self::$connections[$name];
    }
    public static function Close($name = null) {
        if ($name) {
            if (isset ( self::$connections [$name] )) {
                self::$connections [$name]->Close ();
                unset ( self::$connections [$name] );
            }
        } else {
            foreach ( self::$connections as $conn )
                $conn->Close ();
            self::$connections = array ();
        }
    }
    private static function _Setup($name) {
        assert ( isset ( self::$connections [$name] ) );
        
        $conn = self::$connections [$name];
        $options = array (
                'ConnectTimeout' => 0,
                'MaxQueryTime' => 0,
                'MatchMode' => SPH_MATCH_EXTENDED,
                'SortMode' => SPH_SORT_RELEVANCE,
                'RankingMode' => SPH_RANK_PROXIMITY_BM25 
        );
        
        // global $CONFIG;
        if(!self::$config)
        {
            self::$config = (array) new \Config\SphinxClient();
        }

        $config = self::$config;
        $host = isset ( $config[$name]['host'] ) ? $config[$name]['host'] : 'localhost';
        $port = isset ( $config[$name]['port'] ) ? $config[$name]['port'] : 9312;
        
        $conn->SetServer ( $host, $port );
        // overwrite options from $CONFIG
        if (isset ( $config[$name]['options'] )) {
            $cfg_opts = $config[$name]['options'];
            
            foreach ( $options as $opt => & $val )
                if (array_key_exists( $opt, $cfg_opts ))
                    $val = self::__getval ( $cfg_opts [$opt] );
        }
        
        // call SetXXX methods
        foreach ( $options as $opt => $val ) {
            $method = 'Set' . $opt;
            if (method_exists ( $conn, $method ))
                call_user_func ( array (
                        $conn,
                        $method 
                ), $val );
        }
    }
    private static function __getval($val) {
        /* if $val is string and has SPH_ constant */
        if (is_string ( $val ) and strpos ( $val, 'SPH_' ) === 0 and defined ( $val ))
            return constant ( $val );
        return $val;
    }
}
