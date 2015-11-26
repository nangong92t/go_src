<?php
/**
 * CheckerColumn class file.
 *
 * @author TonyXu <tonycbcd@gmail.com>
 * @link http://www.maignet.com/
 * @copyright Copyright &copy; 2012-2013 Maignet LLC
 */

Yii::import('zii.widgets.grid.CGridColumn');

/**
 * CInputColumn represents a grid view column of input.
 *
 */
class CheckerColumn extends CGridColumn
{
	/**
	 * @var string the attribute name of the data model. The corresponding attribute value will be rendered
	 * in each data cell as the checkbox value. Note that if {@link value} is specified, this property will be ignored.
	 * @see value
	 */
	public $name;

	public $value;

  public $baseScriptUrl;

  public $columns = array();
	/**
	 * @var array the HTML options for the data cell tags.
	 */
	public $htmlOptions=array('class'=>'checker-column');
	/**
	 * @var array the HTML options for the header cell tag.
	 */
	public $headerHtmlOptions=array('class'=>'checker-column');
	/**
	 * @var array the HTML options for the footer cell tag.
	 */
	public $footerHtmlOptions=array('class'=>'checker-column');


	/**
	 * Initializes the column.
	 * This method registers necessary client script for the checkbox column.
	 */
	public function init()
	{
    if($this->baseScriptUrl===null)
      $this->baseScriptUrl=Yii::app()->getAssetManager()->publish(Yii::getPathOfAlias('zii.widgets.assets')).'/gridview';

    if ($this->columns)
      $this->initColumns();
	}

	/**
	 * Creates column objects and initializes them.
	 */
	protected function initColumns()
	{
		if($this->columns===array())
		{
			if($this->dataProvider instanceof CActiveDataProvider)
				$this->columns=$this->dataProvider->model->attributeNames();
			elseif($this->dataProvider instanceof IDataProvider)
			{
				// use the keys of the first row of data as the default columns
				$data=$this->dataProvider->getData();
				if(isset($data[0]) && is_array($data[0]))
					$this->columns=array_keys($data[0]);
			}
		}
		$id = 'CheckerColumn';
		foreach($this->columns as $i=>$column)
		{
			if(is_string($column))
				$column=$this->createDataColumn($column);
			else
			{
				if(!isset($column['class']))
					$column['class']='CDataColumn';
				$column=Yii::createComponent($column, $this);
			}
			if(!$column->visible)
			{
				unset($this->columns[$i]);
				continue;
			}
			if($column->id===null)
				$column->id=$id.'_c'.$i;
			$this->columns[$i]=$column;
		}

		foreach($this->columns as $column)
			$column->init();
	}


	/**
	 * Renders the data cell content.
	 * This method renders a checkbox in the data cell.
	 * @param integer $row the row number (zero-based)
	 * @param mixed $data the data associated with the row
	 */
	protected function renderDataCellContent($row,$data)
	{
    $html = '<fieldset>';
    if (isset($data['data']) && is_string($data['data']))
      $data['data'] = json_decode($data['data'], true);



    foreach ($this->columns as $column)
    {
      $html .= '<div class="clearfix">';
      $html .= '<label for="' . $column->id. '">';
      $html .= $column->header . '</label>';
      $html .= '<div class="input'.($column->htmlOptions['class'] ? ' '.$column->htmlOptions['class'] : '').'">';

      $columnType = $column->type ? strtolower($column->type) : 'text';
 
      $methodName = $columnType . 'InputElement';
      if (method_exists($this, $methodName))
        $html .= $this->$methodName($column, $data, $row);
      else
      {
        $options  = $this->getSubHtmlOptions($column, $data, $row);
        $html .= $options['value'];
      }

      if ($column->htmlOptions['unit']) 
        $html .= ' '.$column->htmlOptions['unit'];

      if ($column->htmlOptions['tip'])
        $html .= '<span class="help-block">'. $column->htmlOptions['tip'] .'</span>';

      $html .= '</div>';
      $html .= '</div>';
    }

    $html .= '</fieldset>';

    echo $html;
	}

  private function getSubHtmlOptions($column, $data, $row)
  {
    if (isset($data['data']))
    {
      $id     = $column->id;
      $value  = $data['data'][ $column->value ];
      $colName= $column->value;
    } else {
      $id     = $column->id;
			$value=$this->evaluateExpression($column->value,array('data'=>$data,'row'=>$row));
      $colName  = str_replace('$data->', '', $column->value);
    }

    return array('id'=>$id, 'value'=>$value, 'colName'=>$colName);
  }

  protected function textInputElement($column,$data,$row)
  {
    $options  = $this->getSubHtmlOptions($column, $data, $row);
    $html     = '<input id="'. $options['id'] . '"' .
                ($column->htmlOptions['data-provide'] ?  ' data-provide="'.$column->htmlOptions['data-provide'].'"' : '') .
                ' class="xlarge '.$options['colName'].($column->htmlOptions['data-provide'] ? ' typeahead' : '').'" type="text" name="data['.$row.']['.$options['colName'].']" value="' . $options['value'] . '">';

    return $html;
  }

  protected function textareaInputElement($column, $data, $row)
  {
    $options  = $this->getSubHtmlOptions($column, $data, $row);

    $html   = '<textarea id="'.$options['id'].'" name="data['.$row.']['.$options['colName'].']" rows="10" cols="30" class="span6">'.$options['value'].'</textarea>';
    return $html;
  }

  protected function gpsInputElement($column,$data,$row)
  {
    $valueKey = explode('-', $column->value);
    $lng      = $data['data'][ $valueKey[0] ];
    $lat      = $data['data'][ $valueKey[1] ];

    $html     = '经度:<input id="'.$column->id.'" class="span2 lng" type="text" name="data['.$row.']['.$valueKey[0].']" value="'.$lng.'">';
    $html     .= ' - 纬度:<input class="span2 lat" type="text" name="data['.$row.']['.$valueKey[1].']" value="'.$lat.'">';
    
    return $html;
  }

  protected function linkInputElement($column, $data, $row)
  {
    $value  = $data['data'][ $column->value ];
    $html   = '<a href="'.$value.'" target="_blank">'.($column->htmlOptions['label'] ? $column->htmlOptions['label'] : '查看数据源地址').'</a>';
    return $html;
  }

  protected function locationInputElement($column, $data, $row)
  {
    $cacheName  = 'country';
    $country    = Yii::app()->cache->get($cacheName);
    if (!$country)
    {
      $country  = Country::model()->findAll();
      foreach ($country as $i => $one) { $country[$i] = $one->attributes; }
      Yii::app()->cache->set($cacheName, $country);
    }

    $html   = '<div class="span3">';
    $html   .= '<select name="data['.$row.'][country]" data-change="province" class="cou medium">';
    $html   .= '<option value="0">--请选择国家--</option>';
    foreach ($country as $one) { $html .= '<option value="'.$one['_id'].'">'.$one['name'].'</option>'; }
    $html   .= '</select>';
    $html   .= '</div><div class="span3">';

    $html   .= '<select name="data['.$row.'][province]" data-change="city" class="province medium"><option value="1">--请选择省份--</option></select>';
    $html   .= '</div><div class="span3">';

    $html   .= '<select name="data['.$row.'][city]" class="city medium"><option value="1">--请选择城市--</option></select>';
    $html   .= '</span3>';

    /*
    $cs=Yii::app()->getClientScript();
    $cs->registerScript('', 'customApp.unshift("lib/location");');
     */

    return $html;
  }

  protected function checkerbuttonInputElement($column, $data, $row)
  {
    $html = '<input type="hidden" name="data['.$row.'][roid]" value="'.$data['_id'].'">';
    $html .= '<label class="button"><input type="radio" name="data['.$row.'][action]" value="1" checked="1">通过</label>';
    $html .= '<label class="button"><input type="radio" name="data['.$row.'][action]" value="-1">删除</label>';
    $html .= '<label class="button"><input type="radio" name="data['.$row.'][action]" value="2">保留</label>';
    return $html;
  }

  protected function selectInputElement($column, $data, $row)
  {
    $options  = $this->getSubHtmlOptions($column, $data, $row);
    $value  = $options['value'];

    $html   = '<select name="data['.$row.']['.$options['colName'].']" class="'.$options['colName'].'">';
    if (!is_array($column->htmlOptions['options']) && is_callable($column->htmlOptions['options']))
      $column->htmlOptions['options'] = $column->htmlOptions['options']();

    foreach ($column->htmlOptions['options'] as $id => $one)
      $html .= '<option value="'.$id.'"'.($id==$value ? ' selected' : '').'>'.$one.'</option>';
    $html   .= '</select>';
    return $html;
  }

}
