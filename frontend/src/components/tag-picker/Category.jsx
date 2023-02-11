import { Link } from '@mui/material';
import { Box } from '@mui/system';

function Category({ category, onClick }) {
	return (
		<Box sx={{ p: 2 }}>
			<Link variant="h6" sx={{ cursor: 'pointer' }} onClick={(e) => onClick(category)}>
				{category.title}
			</Link>
		</Box>
	);
}

export default Category;
