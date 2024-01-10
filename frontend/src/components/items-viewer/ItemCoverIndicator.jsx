import { useTheme } from '@emotion/react';
import { Box } from '@mui/material';
import React, { useState } from 'react';

function ItemCoverIndicator({ item, cover, isHighlighted }) {
	const theme = useTheme();
	const [optionsHidden, setOptionsHidden] = useState(true);

	return (
		<Box
			sx={{
				borderRadius: theme.spacing(0.3),
				height: theme.spacing(0.2),
				width: '100%',
				backgroundColor: isHighlighted ? '#d00' : '#500',
			}}
			key={cover.id}
			onMouseEnter={() => setOptionsHidden(false && isHighlighted)}
			onMouseLeave={() => setOptionsHidden(true && isHighlighted)}
			onClick={(e) => {
				e.stopPropagation();
			}}
		></Box>
	);
}

export default ItemCoverIndicator;
