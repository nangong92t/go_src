<?php
namespace Model;

/**
 * Base classes of models which base on database.
 */
abstract class DbBase extends Base
{
    /**
     * 模型的字段名称.
     */
    protected $fields   = array();

    /**
     * 是否为新数据模型.
     */
    protected $isNewRecord_ = true;

    /**
     * 单个模型数据数组.
     */
    private $attributes_    = array();

    /**
     * 被修改过的模型字段.
     */
    private $hasChangedAttributeField_  = array();

    /**
     * Abstract get model table name function.
     * 
     * @return object
     */
    abstract function getTableName();

    /**
     * Returns a model field data. eg: $value=$model->propertyName; $handlers=$model->eventName.
     *
     * @param string $name The property name or event name.
     *
     * @return mixed the property value, event handlers attached to the event, or the named behavior
     */
    public function __get($name)
    {
        if (isset($this->attributes_[$name])) {
            return $this->attributes_[$name];
        } else {
            return null;
        }
    }

    /**
     * Sets value of a model field.
     *
     * @param string $name  The property name.
     * @param mixed  $value The property value.
     *
     * @return mixed
     */
    public function __set($name, $value)
    {
        if (
            !$this->isNewRecord_ && (
                !isset($this->attributes_[$name]) || (
                    $this->attributes_[ $name] != $value &&
                    !isset($this->hasChangedAttributeField_[ $name ])
                )
            )
        ) {
            $this->hasChangedAttributeField_[ $name ]  = $value;
        }
        $this->attributes_[ $name ]  = $value;
    }

    /**
     * Returns current model all attribute values.
     *
     * @return mixed.
     */
    public function getAttributes()
    {
        return $this->attributes_;
    }

    /**
     * Sets current model all attribute values.
     *
     * @param mixed $values The attribute values (name=>value) to be set.
     *
     * @return void
     */
    public function setAttributes($values = array())
    {
        $this->attributes_  = $values;
    }

    /**
     * Return the last sql.
     *
     * @return string
     */
    public function getLastSql()
    {
        return $this->getReadDb()->getLastSql();
    }

    /**
     * Performs the validation.
     *
     * @return void
     */
    public function validate()
    {
        // TODO.
    }

    /**
     * Returns the primary key of the associated database table.Default set primary key is 'uid'. 
     * 
     * @return string
     */
    protected function primaryKey()
    {
        return 'user_id';
    }

    /**
     * Returns the  of the current connection database name.Default set database name is 'kb'.
     * 
     * @return string
     */
    protected function db()
    {
        return 'std';
    }
    
    /**
     * 得到读数据库.
     * 
     * @return \Db\Connection
     */
    public function getReadDb()
    {
        return \Db\Connection::instance()->read($this->db());
    }
    
    /**
     * 得到写数据库.
     * 
     * @return \Db\Connection
     */
    public function getWriteDb()
    {
        return \Db\Connection::instance()->write($this->db());
    }

    /**
     * Add a new data to model.
     *
     * @param mixed $data Array('column_name'=>'value').
     * 
     * @return integer
     */
    public function addNew($data)
    {
        $id         = $this->getWriteDb()->insert($this->getTableName(), $data);
        return $id;
    }

    /**
     * Finds a single record data with the specified condition.
     * 
     * @param mixed  $condition Query condition.
     * @param string $findField For find fields.
     *
     * @return single record data
     */
    public function find($condition = array(), $findField = '*')
    {
        $data   = $this->getReadDb()->select($findField)->from($this->getTableName())->where($condition)->limit(1)->queryRow();

        $this->setAttributes($data);

        if ($data) {
            $this->isNewRecord_ = false;
        }

        return $data;
    }

    /**
     * Finds more record data with the specified condition.
     * 
     * @param mixed   $condition Query condition.
     * @param integer $limit     Item limit number.
     * @param integer $offset    Item offset number.
     * @param string  $findField For select fields.
     * @param string  $order     Order string.
     *
     * @return many record data
     */
    public function findAll($condition = array(), $limit = '', $offset = 0, $findField = '*', $order = null)
    {

        $finder     = $this->getReadDb()
            ->select($findField)
            ->from($this->getTableName());

        if ($condition) {
             $finder->where($condition);
        }
        if ($order) {
            $finder->order($order);
        }
        
        if ($limit) {
            $finder->limit($offset, $limit);
        }

        return $finder->queryAll();
    }

    /**
     * Finds all record data total number with the specified condition.
     * 
     * @param mixed $condition Query condition.
     *
     * @return all record data total number.
     */
    public function findCount($condition = array())
    {
        $data   = $this->getReadDb()->count($this->getTableName(), $condition, 1);

        return $data;
    }

    /**
     * Finds a single record with the specified primary key.
     *
     * @param integer $id        Primary key id.
     * @param string  $findField Select field name string, e.g. 'f1,f2', default '*'.
     * 
     * @return resource
     */
    public function findByPk($id, $findField = '*')
    {
        $data   = $this->getReadDb()
            ->select($findField)
            ->from($this->getTableName())
            ->where(array($this->primaryKey() => $id ))
            ->queryRow();

        $this->setAttributes($data);

        if ($data) {
            $this->isNewRecord_ = false;
        }

        return $data;
    }

    /**
     * Update the data by primary key.
     *
     * @param integer $id       User id.
     * @param mixed   $updateDb Update user data array('column_name'=>'value').
     * 
     * @return void
     */
    public function updateByPk($id, $updateDb = array())
    {
        $this->getWriteDb()->update(
            $this->getTableName(),
            $updateDb,
            array($this->primaryKey() => $id)
        );
    }

    /**
     * Update the data by other condition.
     *
     * @param mixed $condition Update condition array.
     * @param mixed $updateDb  Update user data array('column_name'=>'value').
     * 
     * @return void
     */
    public function update($condition = array(), $updateDb = array())
    {
        $this->getWriteDb()->update(
            $this->getTableName(),
            $updateDb,
            $condition
        );
    }

    /**
     * Save or Update the current model data.
     *
     * @return boolean
     */
    public function save()
    {
        // TODO $this->validate();
        $attributes = $this->getAttributes();
        if (count($attributes) === 0) {
            return false;
        }

        $primaryKey = $this->primaryKey();

        if ($this->isNewRecord_) {
            $id     = $this->getWriteDb()->insert($this->getTableName(), $attributes);
            $this->$primaryKey  = $id;
        } elseif (!empty($this->hasChangedAttributeField_)) {
            $updateData = array();
            foreach ($this->hasChangedAttributeField_ as $k => $v) {
                $updateData[ $k ]   = $v;
            }
            $this->updateByPk($this->$primaryKey, $updateData);
        }

        $this->isNewRecord_ = false;

        return true;
    }

    /**
     * Delete one data by primary key.
     *
     * @param integer $id The primary key id.
     *
     * @return mixed
     */
    public function deleteByPk($id)
    {
        $primaryKey = $this->primaryKey();
        if (!$id && $this->$primaryKey) {
            $id = $this->$primaryKey;
        }
        return $this->deleteAll(array($primaryKey => $id));
    }

    /**
     * Delete all data by conditions.
     *
     * @param mixed $condition The delete condition array.
     *
     * @return mixed
     */
    public function deleteAll($condition = array())
    {
        return $this->getWriteDb()->delete($this->getTableName(), $condition);
    }

}
