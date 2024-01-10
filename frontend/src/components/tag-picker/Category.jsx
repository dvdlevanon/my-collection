import { useTheme } from '@emotion/react';
import { Link } from '@mui/material';
import { Box } from '@mui/system';

function Category({ isHighlighted, category, onClick }) {
	const theme = useTheme();

	return (
		<Box
			backgroundColor={isHighlighted ? 'dark.lighter' : 'auto'}
			sx={{
				padding: theme.multiSpacing(0, 1),
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
