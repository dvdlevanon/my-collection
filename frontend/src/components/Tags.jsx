import styles from './Tags.module.css';
import Tag from "./Tag"

function Tags({ tags, onTagSelected }) {
  return (
    <div className={styles.tags}>
        {tags.map((tag) => {
          return (
            <div key={tag.id}>
                <Tag tag={tag} onTagSelected={onTagSelected}/>
            </div>
        )})}
    </div>
  )
}

export default Tags