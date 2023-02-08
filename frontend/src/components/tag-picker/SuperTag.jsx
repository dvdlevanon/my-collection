import { Link } from '@mui/material';
import { Box } from '@mui/system';

function SuperTag({ superTag, onSuperTagClicked }) {
	return (
		<Box sx={{ p: 2 }}>
			<Link variant="h6" sx={{ cursor: 'pointer' }} onClick={(e) => onSuperTagClicked(superTag)}>
				{superTag.title}
			</Link>
		</Box>
	);
}

export default SuperTag;
