<?php

Yii::import('zii.widgets.grid.CGridView');
Yii::import('system.myshow.widgets.grid.*');

class MyshowGridView extends CGridView
{
  public $setDataIdOnTr;

  protected $dataCount=0;

  /**
	 * Renders the tr items for the grid view.
	 */
	public function renderTableRow($row)
	{
    if ($this->setDataIdOnTr)
    {
      $data   = $this->dataProvider->data[$row];
      $curId  = $data->_id;
    }

    if($this->rowCssClassExpression!==null)
    {
      $data=$this->dataProvider->data[$row];
      echo '<tr'.($curId ? ' data-id="'.$curId.'"' : '').' class="'.$this->evaluateExpression($this->rowCssClassExpression,array('row'=>$row,'data'=>$data)).'">';
    }
    else if(is_array($this->rowCssClass) && ($n=count($this->rowCssClass))>0)
      echo '<tr'.($curId ? ' data-id="'.$curId.'"' : '').' class="'.$this->rowCssClass[$row%$n].'">';
    else
      echo '<tr'.($curId ? ' data-id="'.$curId.'"' : '').'>';
    foreach($this->columns as $column)
      $column->renderDataCell($row);
    echo "</tr>\n";

	}

	/**
	 * Renders a data cell.
	 * @param integer $row the row number (zero-based)
	 */
	public function renderDataCell($row)
	{
		$data=$this->grid->dataProvider->data[$row];
		$options=$this->htmlOptions;
		if($this->cssClassExpression!==null)
		{
			$class=$this->evaluateExpression($this->cssClassExpression,array('row'=>$row,'data'=>$data));
			if(!empty($class))
			{
				if(isset($options['class']))
					$options['class'].=' '.$class;
				else
					$options['class']=$class;
			}
		}
		echo CHtml::openTag('td',$options);
		$this->renderDataCellContent($row,$data,$this->dataCount++);
		echo '</td>';
	}

  public function registerClientScript()
  {
    return false;
  }

}
