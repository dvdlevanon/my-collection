import React from 'react'
import { IconButton } from '@mui/material';
import CloseIcon from '@mui/icons-material/Close';

function ActiveTag({tag, onTagDeactivated, onTagSelected, onTagDeselected }) {

  const onTagClicked = (e) => {
    if (tag.selected) {
      onTagDeselected(tag);
    } else {
      onTagSelected(tag);
    }
  }
  
  const onCloseClicked = (e) => {
    e.stopPropagation()
    onTagDeactivated(tag)
  }

  return (
    <div className={tag.selected ? "active-tag selected" : "active-tag"} onClick={(e) => onTagClicked(e)}>
      <IconButton onClick={(e) => onCloseClicked(e)}>
        <CloseIcon/>
      </IconButton>
      {tag.title}
    </div>
  )
}

export default ActiveTag