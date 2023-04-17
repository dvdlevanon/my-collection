import { Box, Stack } from '@mui/material';
import React from 'react';
import AspectRatioUtil from '../../utils/aspect-ratio-util';
import Item from '../items-viewer/Item';

function SubItem({ item, itemWidth, highlighted }) {
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
			{/* <Box
				component="img"
				src={ItemsUtil.getCover(item, 0)}
				alt={item.title}
				loading="lazy"
				sx={{
					width: imageWidth,
					height: AspectRatioUtil.calcHeight(imageWidth, AspectRatioUtil.asepctRatio16_9),
					objectFit: 'contain',
					cursor: 'pointer',
					borderRadius: '10px',
				}}
			/>
			<Stack flexDirection="row" padding="5px" gap="10px">
				<Tooltip title={item.title} arrow followCursor>
					<Typography
						variant="caption"
						sx={{
							textOverflow: 'ellipsis',
							wordWrap: 'break-word',
							overflow: 'hidden',
							maxHeight: '3.6em',
						}}
					>
						{item.title}
					</Typography>
				</Tooltip>
			</Stack> */}
		</Stack>
	);
}

export default SubItem;
