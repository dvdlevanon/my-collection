import styles from './SuperTags.module.css';
import SuperTag from './SuperTag';

function SuperTags({superTags, onSuperTagSelected, onSuperTagDeselected}) {
  return (
    <div className={styles.super_tags}>
      {superTags.map((tag) => {
        return <SuperTag key={tag.id} superTag={tag} onSuperTagSelected={onSuperTagSelected} onSuperTagDeselected={onSuperTagDeselected}/>
      })}
    </div>
  )
}

export default SuperTags