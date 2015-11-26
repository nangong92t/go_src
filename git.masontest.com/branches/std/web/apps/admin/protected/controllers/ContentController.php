<?php
    /**
     * Topic Administrator controller
     *
     * @author sam<@gmail.com>
     * @date 2014-10-06
     */
    class ContentController extends Controller {
        public function actionTopics($page = 1)
        {
            $limit  = 20;
            $token  = $this->getToken();
            $res    = RpcClient_STD_Topic::Instance()->getTopicList($token, $page, "", $limit, 0, true);
            $data   = $res['data'];
            $total  = isset($data['total']) ? $data['total'] : 0;

            $pages              = new CPagination($total);
            $pages->pageSize    = $limit;
            $pages->applyLimit(new CDbCriteria());
            
            $params['topics']   = $data['list'] ? $data['list'] : array();
            $params['pages']    = $pages;
            $params['total']    = $total;
            $params['current_page'] = $pages->getOffset();
            
            
            $this->jsApp    = array('controller/list', 'Init', array("type"=>"topic", "apis"=>array(
                "delete" => "admin/removetopics"
            )));
            $params['current_time'] = date('y-m-d h:i:s',time());
            $params['current_from'] = $pages->getCurrentPage() * $limit + 1;
            $params['current_to'] = $params['current_from'] + count($params['topics']) - 1;
            
            $this->render('topics', $params);
        }

        public function actionComments($page = 1)
        {
            $limit  = 20;
            $res    = RpcClient_STD_Topic::Instance()->getAllCommentList($page, $limit);
            $data   = $res['data'];
            $total  = isset($data['total']) ? $data['total'] : 0;
            $pages              = new CPagination($total);
            $pages->pageSize    = $limit;
            $pages->applyLimit(new CDbCriteria());
            
            $params['comments'] = $data['list'];
            $params['pages']    = $pages;
            $params['total']    = $total;
            
            $this->jsApp    = array('controller/list', 'Init', array(
                "type"  => "comment",
                "apis"  => array(
                    "delete" => "admin/removecomments"
                )
            ));
            $params['current_time'] = date('y-m-d h:i:s',time());
            $params['current_from'] = $pages->getCurrentPage() * $limit + 1;
            $params['current_to'] = $params['current_from'] + count($params['comments']) - 1;
            $this->render('comments', $params);
        }

        public function actionLabels($page = 1)
        {
            $limit  = 20;

            $token  = $this->getToken();
            $res    = RpcClient_STD_Admin::Instance()->getLabelList($token, $page, $limit);
            $data   = $res['data'];
            $total  = isset($data['total']) ? $data['total'] : 0;

            $pages              = new CPagination($total);
            $pages->pageSize    = $limit;
            $pages->applyLimit(new CDbCriteria());
            
            $params['list']     = $data['list'] ? $data['list'] : array();
            $params['pages']    = $pages;
            $params['total']    = $total;
            
            $this->jsApp    = array('controller/list', 'Init', array(
                "type"  => "label",
                "apis"  => array(
                    "delete" => "admin/removelabels"
                )
            ));
            $params['current_time'] = date('y-m-d h:i:s',time());
            $params['current_from'] = $pages->getCurrentPage() * $limit + 1;
            $params['current_to'] = $params['current_from'] + count($params['list']) - 1;
            $this->render('labels', $params);
        }

        public function actionBackground_photos($page = 1)
        {
            $limit  = 20;

            $token  = $this->getToken();
            $res    = RpcClient_STD_Admin::Instance()->getBackgroundPhotoList($token, $page, $limit);
            $data   = $res['data'];
            $total  = isset($data['total']) ? $data['total'] : 0;

            $pages              = new CPagination($total);
            $pages->pageSize    = $limit;
            $pages->applyLimit(new CDbCriteria());
            
            $params['list']     = $data['list'] ? $data['list'] : array();
            $params['pages']    = $pages;
            $params['total']    = $total;
            
            $this->jsApp    = array('controller/list', 'Init', array(
                "type"  => "background",
                "apis"  => array(
                    "delete" => "admin/removebgs"
                )
            ));
            $params['current_time'] = date('y-m-d h:i:s',time());
            $params['current_from'] = $pages->getCurrentPage() * $limit + 1;
            $params['current_to'] = $params['current_from'] + count($params['list']) - 1;
            $this->render('background_photos', $params);
        }

        public function actionUnwanted_words($page = 1)
        {
            $limit  = 20;
            
            $token  = $this->getToken();
            $res    = RpcClient_STD_Admin::Instance()->GetUnwantWordList($token, $page, $limit);
            $data   = $res['data'];
            $total  = isset($data['total']) ? $data['total'] : 0;
            
            $pages              = new CPagination($total);
            $pages->pageSize    = $limit;
            $pages->applyLimit(new CDbCriteria());
            
            $params['list']     = $data['list'] ? $data['list'] : array();
            $params['pages']    = $pages;
            $params['total']    = $total;
            
            $this->jsApp    = array('controller/list', 'Init', array(
                    "type"  => "unwanted_words",
                    "apis"  => array(
                            "delete" => "admin/RemoveUnwantWord"
                    )
            ));
            //print $pages->getOffset();
            //print $pages->getPageSize();
            $params['current_from'] = $pages->getCurrentPage() * $limit + 1; 
            $params['current_to'] = $params['current_from'] + count($params['list']) - 1;
//             echo '<pre>';
//             var_dump($params);die;
//             echo '</pre>';
			$params['current_time'] = date('y-m-d h:i:s',time());
            $this->render('unwanted_words', $params);
        }
    }
    
?>
