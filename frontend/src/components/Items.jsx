import Item from './Item';
import styles from './Items.module.css';

function ItemsList({ items, previewMode }) {
	return (
		<div className={styles.items}>
			{items.map((item) => {
				return (
					<div key={item.id}>
						<Item item={item} preferPreview={previewMode} />
					</div>
				);
			})}
		</div>
	);
}

export default ItemsList;
