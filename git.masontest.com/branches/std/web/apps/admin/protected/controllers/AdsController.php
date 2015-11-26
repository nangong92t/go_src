<?php
    class AdsController extends Controller
    {
        public function actionLanding() {
            $this->render('landing');
        }
        public function actionInner() {
            $this->render('inner');
        }
        public function actionWords() {
            $this->render('Words');
        }
    }
    
?>
