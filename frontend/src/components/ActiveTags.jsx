import styles from './ActiveTags.module.css';
import ActiveTag from './ActiveTag';

function ActiveTags({ activeTags, onTagDeactivated, onTagSelected, onTagDeselected }) {
	return (
		<div className={styles.active_tags}>
			{activeTags.map((tag) => {
				return (
					<ActiveTag
						key={tag.id}
						tag={tag}
						onTagDeactivated={onTagDeactivated}
						onTagSelected={onTagSelected}
						onTagDeselected={onTagDeselected}
					/>
				);
			})}
		</div>
	);
}

export default ActiveTags;
