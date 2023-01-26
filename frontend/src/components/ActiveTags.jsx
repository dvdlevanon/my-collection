import { Box } from '@mui/material';
import ActiveTag from './ActiveTag';
import styles from './ActiveTags.module.css';

function ActiveTags({ activeTags, onTagDeactivated, onTagSelected, onTagDeselected }) {
	return (
		<Box className={styles.active_tags}>
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
		</Box>
	);
}

export default ActiveTags;
