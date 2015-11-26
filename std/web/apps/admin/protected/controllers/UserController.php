<?php
    class UserController extends Controller {
        public function actionBasic($page = 1, $utype = 'active') {
            $limit  = 20;
            $token  = $this->getToken();
            $res    = RpcClient_STD_Admin::Instance()->getUserList($token, $page, $limit, $utype);

            $subTypes   = array(
                'active'    => 'Users have 20+ posts in 24H',
                'blocked'   => 'Blocked by others',
                'block'     => 'Blocked users',
            );

            $data   = $res['data'];
            $total  = isset($data['total']) ? $data['total'] : 0;
            $pages              = new CPagination($total);
            $pages->pageSize    = $limit;
            $pages->applyLimit(new CDbCriteria());
            
            $params['list']     = $data['list'] ? $data['list'] : array();
            $params['pages']    = $pages;
            $params['total']    = $total;
            $params['gender']   = array(1=>'Man', 2=>'Woman');
            $params['subTypes'] = $subTypes;
            $params['curType']  = $utype;
            
            $this->jsApp    = array('controller/list', 'Init', array("type"=>"user", "apis"=>array(
                "delete" => "admin/removeusers"
            )));
            $this->render('basic', $params);
        }

        public function actionAdvanced($page = 1, $utype = 'sameip') {
            $limit  = 20;
            $token  = $this->getToken();
            $res    = RpcClient_STD_Admin::Instance()->getUserList($token, $page, $limit, $utype);

            $subTypes   = array(
                'sameip'    => 'Users from same IP',
                'sameudid'  => 'Users from same UDID',
            );

            $data   = $res['data'];
            $total  = isset($data['total']) ? $data['total'] : 0;
            $pages              = new CPagination($total);
            $pages->pageSize    = $limit;
            $pages->applyLimit(new CDbCriteria());
            
            $params['list']     = $data['list'] ? $data['list'] : array();
            $params['pages']    = $pages;
            $params['total']    = $total;
            $params['gender']   = array(1=>'Man', 2=>'Woman');
            $params['subTypes'] = $subTypes;
            $params['curType']  = $utype;
            
            $this->jsApp    = array('controller/list', 'Init', array("type"=>"user", "apis"=>array(
                "delete" => "admin/RemoveUnwantWord"
            )));

            $this->render('basic', $params);
        }
    }
?>
