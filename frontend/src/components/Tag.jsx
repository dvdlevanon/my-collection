import styles from './Tag.module.css';

function Tag({ tag, onTagSelected }) {
	return (
		<>
			<div
				className={styles.tag + ' ' + (tag.selected ? styles.selected : styles.unselected)}
				onClick={() => onTagSelected(tag)}
			>
				{tag.title}
			</div>
		</>
	);
}

export default Tag;
