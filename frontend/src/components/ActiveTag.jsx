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

  return (
    <div className={tag.selected ? "active-tag selected" : "active-tag"} onClick={(e) => onTagClicked(e)}>
      <IconButton onClick={() => onTagDeactivated(tag)}>
        <CloseIcon/>
      </IconButton>
      {tag.title}
    </div>
  )
}

export default ActiveTag