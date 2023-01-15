import { Switch } from '@mui/material';
import ActiveTags from './ActiveTags';
import styles from './GalleryFilters.module.css';

function SidePanel({ activeTags, onTagDeactivated, onTagSelected, onTagDeselected, onChangeCondition }) {
	const onConditionChanged = (e) => {
		onChangeCondition(e.target.checked ? '&&' : '||');
	};

	return (
		<div className={styles.side_panel}>
			{activeTags.length > 0 ? (
				<ActiveTags
					activeTags={activeTags}
					onTagDeactivated={onTagDeactivated}
					onTagSelected={onTagSelected}
					onTagDeselected={onTagDeselected}
				/>
			) : (
				''
			)}
			{activeTags.length > 1 ? (
				<div className={styles.condition_switch}>
					<span>||</span>
					<Switch onChange={(e) => onConditionChanged(e)} />
					<span>&&</span>
				</div>
			) : (
				''
			)}
		</div>
	);
}

export default SidePanel;
