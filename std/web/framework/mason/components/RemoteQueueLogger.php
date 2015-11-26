<?php

class RemoteQueueLogger extends CComponent {
  public $servers;
  private $_pub;

  public function init() {
    $this->_pub = new ZMQSocket(new ZMQContext, ZMQ::SOCKET_PUSH);
    foreach ($this->servers as $server) {
      $this->_pub->connect($server);
    }
  }

  public function log($message, $category = 'cooperation') {
    $data = json_encode(array(
      'created_at' => time(),
      'level' => 'cooperation',
      'category' => $category,
      'message' => $message,
    ));
    $this->_pub->send($data, ZMQ::MODE_NOBLOCK);

    #file_put_contents('/tmp/aaa', $data . PHP_EOL, FILE_APPEND);
  }
}
