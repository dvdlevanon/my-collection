import React from 'react'
import { IconButton } from '@mui/material';
import CloseIcon from '@mui/icons-material/Close';

function ItemTag({tag}) {
  const onRemoveClicked = (e) => {
    e.stopPropagation()
    console.log("Remove " + tag.id)
  }

  return (
    <div className="item-tag">
      <IconButton onClick={(e) => onRemoveClicked(e)}>
        <CloseIcon/>
      </IconButton>
      {tag.title}
    </div>
  )
}

export default ItemTag
