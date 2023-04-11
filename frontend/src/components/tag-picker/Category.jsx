import { Link } from '@mui/material';
import { Box } from '@mui/system';

function Category({ isHighlighted, category, onClick }) {
	return (
		<Box
			backgroundColor={isHighlighted ? 'dark.lighter' : 'auto'}
			sx={{
				padding: '0px 10px',
				height: '100%',
			}}
		>
			<Link
				variant="h6"
				onClick={(e) => onClick(category)}
				sx={{
					cursor: 'pointer',
				}}
			>
				{category.title}
			</Link>
		</Box>
	);
}

export default Category;
