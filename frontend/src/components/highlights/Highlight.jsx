import { Box, Stack } from '@mui/material';
import React from 'react';
import AspectRatioUtil from '../../utils/aspect-ratio-util';
import Item from '../items-viewer/Item';

function Highlight({ item, itemWidth, highlighted }) {
	return (
		<Stack
			flexDirection="row"
			gap="10px"
			padding="10px"
			sx={{
				cursor: 'pointer',
				borderRadius: '10px',
				backgroundColor: highlighted ? 'dark.lighter2' : 'unset',
				'&:hover': {
					backgroundColor: 'dark.lighter',
				},
			}}
		>
			<Box>
				<Item
					item={item}
					preferPreview={true}
					itemWidth={itemWidth}
					itemHeight={AspectRatioUtil.calcHeight(itemWidth, AspectRatioUtil.asepctRatio16_9)}
					direction="row"
					showOffests={true}
					titleSx={{
						whiteSpace: 'normal',
						lineHeight: '1.5em',
						maxHeight: '3em',
						textAlign: 'start',
					}}
					itemLinkBuilder={(item) => {
						return '/spa/item/' + item.id;
					}}
				/>
			</Box>
		</Stack>
	);
}

export default Highlight;
