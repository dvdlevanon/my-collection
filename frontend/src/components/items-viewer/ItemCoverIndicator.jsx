import { Box } from '@mui/material';
import React, { useState } from 'react';

function ItemCoverIndicator({ item, cover, isHighlighted }) {
	let [optionsHidden, setOptionsHidden] = useState(true);

	return (
		<Box
			sx={{
				borderRadius: '3px',
				height: '2px',
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
