import styles from './ItemTag.module.css';
import { IconButton } from '@mui/material';
import CloseIcon from '@mui/icons-material/Close';

function ItemTag({tag, onRemoveClicked}) {
  const removeHandler = (e) => {
    e.stopPropagation()
    onRemoveClicked(tag);
  }

  return (
    <div className={styles.item_tag}>
      <IconButton onClick={(e) => removeHandler(e)}>
        <CloseIcon/>
      </IconButton>
      {tag.title}
    </div>
  )
}

export default ItemTag
