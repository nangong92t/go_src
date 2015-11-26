<?php
/**
 * Redis的配置文件.
 * 
 * @author yongf <yongf@jumei.com>
 */

namespace Config;

/**
 * Redis的配置文件.
 */
class Redis
{
    /**
     * Configs of Redis.
     * @var array
     */
    public $default = array(
        'nodes' => array(
            array('master' => "127.0.0.1:6379", 'slave' => "127.0.0.1:6379"),
        ),
        'db' => 0
    );
    public $storage = array(
        'nodes' => array(
            array('master' => "127.0.0.1:6379", 'slave' => "127.0.0.1:6379"),
        ),
        'db' => 2
    );

    public $RedisKeysTypeList   = array (
        //    "xiaomei_uids" => "cache",
        //    "wxfree_detail_index_" => "cache",
        //    "WashData" => "cache",
        "virtual_bundle_set" => "storage",
        'user_number' => 'storage',
        'user_click' => 'storage',
        'unfollow_user_click' => 'storage',
        'unfollow_user_number' => 'storage',
        //    "user_written_report_on_products_" => "cache",
        "user_reports_image_info_detail_" => "storage",
        "user_reports_image_info_" => "storage",
        "user_normal_fans" => "storage",
        "user_normal_attention" => "storage",
        //    "user_for_product_unreported_" => "cache",
        "user_following_" => "storage",
        "user_followed_by_" => "storage",
        "user_fans" => "storage",
        "user_daren_fans" => "storage",
        "user_daren_attention" => "storage",
        "user_attention" => "storage",
        //    "TransferFallImg" => "cache",
        //    "TransferActivities" => "cache",
        "thumb_file_" => "storage",
        "thumb_file_" => "storage",
        "site_map_cache_review" => "storage",
        "site_map_cache_report_list" => "storage",
        "site_map_cache_product_reviews" => "storage",
        "select_koubei_data" => "storage",
        //    "sec_id_" => "cache",
        "report_user_feel_useful_" => "storage",
        "report_time_number_" => "storage",
        "report_thumb_key_id_" => "storage",
        "report_thumb_info_key_id_" => "storage",
        "report_draft_" => "storage",
        "random_valuable_user_info" => "storage",
        "random_valuable_report" => "storage",
        "random_valuable_product" => "storage",
        "random_valuable_daren_info" => "storage",
        "product_with_uid_report" => "storage",
        "product_with_uid_dianping" => "storage",
        //    "product_user_function_info" => "cache",
        //    "Product_Stats_Data_" => "cache",
        "product_statis_users_count_" => "storage",
        "product_statis_users_count_" => "storage",
        //    "product_function_key" => "cache",
        //    "product_function_detail_key_" => "cache",
        "product_category_avg_score_" => "storage",
        "product_blocked_set_report_list" => "storage",
        //    "prod_cate_id" => "cache",
        //    "moveddata_" => "cache",
        "moved_user_" => "storage",
        "luckybox_products_like" => "storage",
        "luckybox_coupon_" => "storage",
        "last_view_product" => "storage",
        //    "Koubei_user_product_number_" => "cache",
        "Koubei_User_Latest_View_" => "storage",
        //    "Koubei_User_Info_ " => "cache",
        "koubei_user_event_" => "storage",
        "koubei_user_event_" => "storage",
        "Koubei_user_data_check_" => "storage",
        "koubei_system_user_count" => "storage",
        "koubei_system_user_" => "storage",
        "koubei_system_GLOBAL" => "storage",
        //    "koubei_product_map_" => "cache",
        "koubei_group_" => "storage",
        //    "koubei_function_map_" => "cache",
        "koubei_event_center_message_log" => "storage",
        //    "koubei_correct_user_event_" => "cache",
        //    "Koubei_correct_user_data_" => "cache",
        "koubei_comment_count_" => "storage",
        //    "koubei_category_map_" => "cache",
        //    "koubei_brand_map_" => "cache",
        "Koubei_bayesian_word_total_num_spam" => "storage",
        "Koubei_bayesian_word_total_num_ham" => "storage",
        "Koubei_bayesian_word_items_spam" => "storage",
        "Koubei_bayesian_word_items_ham" => "storage",
        //    "Koubei_All_Products" => "cache",
        //    "jumei_product" => "cache",
        //    "jumei_function" => "cache",
        //    "jumei_category" => "cache",
        //    "jumei_brand" => "cache",
        "ip_black_lists" => "storage",
        //    "ImgThumb" => "cache",
        //    "iccenter_product" => "cache",
        //    "iccenter_function" => "cache",
        //    "iccenter_category" => "cache",
        //    "iccenter_brand" => "cache",
        "how_hot_brand_list" => "storage",
        //    "fresh_detail_index_" => "cache",
        //    "free_detail_index_" => "cache",
        //    "free_apply_term_" => "cache",
        //    "fourth_third_category_map_" => "cache",
        //    "event_user_last_view_time_" => "cache",
        "event_sys_last_time_USER_" => "storage",
        "event_sys_last_time_GLOBAL" => "storage",
        "each_day_keword_count" => "storage",
        "DisposalDataIncremental" => "storage",
        "delay_insert_user_product_report" => "storage",
        "day_keyword" => "storage",
        "daily_unfollow_uids_" => "storage",
        "daily_follow_uids_" => "storage",
        "daily_follow_stats_" => "storage",
        //    "cover_choice_List" => "cache",
        "cover_choice_List" => "storage",
        //    "comment_time_number" => "cache",
        //    "comment_time_number" => "cache",
        "center_for_product" => "storage",
        "center_for_category" => "storage",
        "center_for_brand" => "storage",
        //    "CalculateInit" => "cache",
        //    "Calculate" => "cache",
        //    "cache_key_temp_report_view_count" => "cache",
        //    "brand_statistic_items_brand_statistic_items_" => "cache",
        //    "brand_statistic_items_" => "cache",
        "boms_relations" => "storage",
        "blocked_products_id_hash" => "storage",
        //    "all_product_categories" => "cache",
        "all_buyer_distribution" => "storage",
        //    "all_brands" => "cache",
        //    "all_admin_privilege_list" => "cache",
        //    "address_list_cache" => "cache",
        //    "activity_clinique" => "cache",
        "user_new_followed_count" => "storage",
        "user_new_free_numbers" => "storage",
        "current_dynamic_num" => "storage",
    );



}
