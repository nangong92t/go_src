<?php

/**
 * innerCForm class for use custom the form view
 */
class innerCForm extends CForm
{

  /**
   * Overloaded the parent function
   */
  public function renderElement($element)
  {
    if(is_string($element))
    {
      if(($e=$this[$element])===null && ($e=$this->getButtons()->itemAt($element))===null)
        return $element;
      else
        $element=$e;
    }
    if($element->getVisible())
    {
      if($element instanceof CFormInputElement)
      {
        $model          = $this->getModel();
        $modelClassName = get_class($model);
        $elData         = $element->attributes;
        if($element->type==='hidden')
          return "<div style=\"visibility:hidden\">\n".$element->render()."</div>\n";
        else
          return "<div class=\"clearfix field_{$element->name}\">\n".
                  "<label for=\"".
                    (isset($elData['id']) ? $elData['id'] :
                    $modelClassName.'_'.$element->name) .
                    "\">". ($element->getRequired()? '<span class="required">*</span>' : '').$element->getLabel() ."：</label>\n".
                  "<div class=\"input\">\n".
                  $element->renderInput().
                  "<span class=\"help-block join_add_textfield".($model->getError($element->name) ? ' warn' : '')."\">".($model->getError($element->name) ? $model->getError($element->name) : $element->hint)."</span>\n".
                  "</div>\n".
                  "</div>\n";
      }
      else if($element instanceof CFormButtonElement)
        return $element->render()."\n";
      else
        return $element->render();
    }
    return '';
  }


	/**
	 * Renders the {@link buttons} in this form.
	 * @return string the rendering result
	 */
	public function renderButtons()
	{
		$output='';
		foreach($this->getButtons() as $button)
			$output.=$this->renderElement($button);
		return $output!=='' ? "<div class=\"row actions\">".$output."</div>\n" : '';
	}

}
