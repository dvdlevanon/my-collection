import styles from './SuperTag.module.css';

function SuperTag({ superTag, onSuperTagClicked }) {
	return (
		<div className={styles.super_tag} onClick={(e) => onSuperTagClicked(superTag)}>
			{superTag.title}
		</div>
	);
}

export default SuperTag;
