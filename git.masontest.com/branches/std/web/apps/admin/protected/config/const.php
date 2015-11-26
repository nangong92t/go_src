<?php

class ProjectParams
{
    /**
     * 菜单栏
     */
     public static $menus = array(
        'admin/index' => array(
            'name' => 'Home',
            ),
         'user/index' => array(
            'name' => 'User',
             'subMenus' => array(
                'user/basic' => array(
                    'name' => 'Basic',
                    ),
                 'user/advanced' => array(
                    'name' => 'Advanced',
                    ),
                )
            ),
         'content/index' => array(
            'name' => 'Content',
             'subMenus' => array(
                'content/topics' => array(
                    'name' => 'topics',
                    ),
                 'content/comments' => array(
                    'name' => 'comments',
                    ),
                 'content/labels' => array(
                    'name' => 'Labels',
                    ),
                 'content/background_photos' => array(
                    'name' => 'Background Photos',
                    ),
                 'content/unwanted_words' => array(
                    'name' => 'Unwanted words',
                    ),
                )
            ),
         'ads/index' => array(
            'name' => 'Ads',
             'subMenus' => array(
                'ads/landing' => array(
                    'name' => 'Landing',
                    ),
                 'ads/inner' => array(
                    'name' => 'Inner',
                    ),
                 'ads/words' => array(
                    'name' => 'Words',
                    ),
                )
            ),
            'admin/logout' => array(
                'name'  => 'Logout'
            )
        );
    
}


/**
 * * 定义一些资源常量
 */
class ProfileItems
{
    /**
     * 性别
     */
     const MAN = 1;
    
     const WOMAN = 2;
    
     public static $gender = array(
        self :: MAN => '男性',
         self :: WOMAN => '女性'
        );
    
    /**
     * 艺人身高
     */
     const MINHEIGHT = 50;
     const MAXHEIGHT = 250;
    
    /**
     * 血型
     */
     public static $bloods = array(
        1 => 'A', 2 => 'B', 3 => 'AB', 4 => 'O'
        );
    
    
    }
