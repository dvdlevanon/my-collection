import styles from './ItemTags.module.css';
import ItemTag from './ItemTag';
import { Fab } from '@mui/material';
import AddIcon from '@mui/icons-material/Add';

function ItemTags({ item, onAddTag, onTagRemoved }) {
	return (
		<div className={styles.item_tags_editor}>
			<div className={styles.item_tags}>
				{item.tags
					? item.tags.map((tag) => {
							return <ItemTag key={tag.id} tag={tag} onRemoveClicked={onTagRemoved} />;
					  })
					: ''}
			</div>
			<div className={styles.add_tag_button}>
				<Fab
					onClick={(e) => {
						onAddTag();
					}}
				>
					<AddIcon />
				</Fab>
			</div>
		</div>
	);
}

export default ItemTags;
