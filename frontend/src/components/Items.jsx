import styles from './Items.module.css';
import Item from './Item';

function ItemsList({ items }) {
	return (
		<div className={styles.items}>
			{items.map((item) => {
				return (
					<div key={item.id}>
						<Item item={item} />
					</div>
				);
			})}
		</div>
	);
}

export default ItemsList;
