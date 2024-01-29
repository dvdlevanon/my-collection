import { useTheme } from '@emotion/react';
import { Link } from '@mui/material';
import { Box } from '@mui/system';

function Category({ isHighlighted, category, onClick }) {
	const theme = useTheme();

	return (
		<Box
			sx={{
				backdropFilter: isHighlighted ? 'brightness(2)' : 'auto',
				padding: theme.multiSpacing(0, 1),
				height: '100%',
				borderRadius: theme.spacing(0.3),
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
