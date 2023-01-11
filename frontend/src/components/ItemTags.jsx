import styles from './ItemTags.module.css';
import ItemTag from './ItemTag'

function ItemTags({item}) {
  return (
    <div className={styles.item_tags}>
        {item.tags.map((tag) => {
            return <ItemTag key={tag.id} tag={tag} />
        })}
    </div>
  )
}

export default ItemTags