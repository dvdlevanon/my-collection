import styles from './Tags.module.css';
import Tag from "./Tag"

function Tags({ tags, onTagActivated }) {
  return (
    <div className={styles.tags}>
        {tags.map((tag) => {
          return (
            <div key={tag.id}>
                <Tag tag={tag} onTagActivated={onTagActivated}/>
            </div>
        )})}
    </div>
  )
}

export default Tags