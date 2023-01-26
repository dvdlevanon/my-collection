import { Link } from '@mui/material';
import { Box } from '@mui/system';
import styles from './SuperTag.module.css';

function SuperTag({ superTag, onSuperTagClicked }) {
	return (
		<Box sx={{ p: 2 }}>
			<Link variant="h6" className={styles.super_tag} onClick={(e) => onSuperTagClicked(superTag)}>
				{superTag.title}
			</Link>
		</Box>
	);
}

export default SuperTag;
