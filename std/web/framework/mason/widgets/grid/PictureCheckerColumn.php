<?php
/**
 * CInputColumn class file.
 *
 * @author TonyXu <tonycbcd@gmail.com>
 * @link http://www.maignet.com/
 * @copyright Copyright &copy; 2012-2013 Maignet LLC
 */


class PictureCheckerColumn extends CheckerColumn
{
  /**
	 * Renders the data cell content.
	 * This method renders a checkbox in the data cell.
	 * @param integer $row the row number (zero-based)
	 * @param mixed $data the data associated with the row
	 */
	protected function renderDataCellContent($row,$data)
	{
    if (is_string($data['data']))
      $data['data'] = json_decode($data['data'], true);

    $pictures = $data['data'][ $this->value ];
    if (!$pictures) return false;

    $html = '<ul class="thumbnails">';
    foreach ($pictures as $i => $one) 
      $html .= '<li><input type="checkbox" checked="1" name="data['.$row.'][photo]['.$i.']" value="'.$one. '"><a href="'.$one.'" target="_checkPic"><img width="100" src="'.$one.'"></a></li>';

    $html .= '</ul>';

    echo $html;
	}
}
