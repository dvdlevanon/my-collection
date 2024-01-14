import { useTheme } from '@emotion/react';
import { Box, Stack } from '@mui/material';
import React from 'react';
import AspectRatioUtil from '../../utils/aspect-ratio-util';
import Item from '../items-viewer/Item';

function Highlight({ item, itemWidth, highlighted }) {
	const theme = useTheme();

	return (
		<Stack
			flexDirection="row"
			gap={theme.spacing(1)}
			padding={theme.spacing(1)}
			borderRadius={theme.spacing(1)}
			sx={{
				cursor: 'pointer',
				backgroundColor: highlighted ? theme.palette.primary.light : 'unset',
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
					showOffests={true}
					itemLinkBuilder={(item) => {
						return '/spa/item/' + item.id;
					}}
				/>
			</Box>
		</Stack>
	);
}

export default Highlight;
