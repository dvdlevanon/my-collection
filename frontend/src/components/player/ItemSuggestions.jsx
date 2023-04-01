import { Box, Grid, Stack } from '@mui/material';
import React from 'react';
import AspectRatioUtil from '../../utils/aspect-ratio-util';
import Item from '../items-viewer/Item';

function ItemSuggestions({ suggestedItems, width, onBackgroundClicked }) {
	const getItemSize = () => {
		let itemWidth = width / 5;
		return {
			width: itemWidth,
			height: AspectRatioUtil.calcHeight(itemWidth, AspectRatioUtil.asepctRatio16_9),
		};
	};

	return (
		<Stack
			flexDirection="column"
			onClick={onBackgroundClicked}
			sx={{
				position: 'absolute',
				left: 50,
				top: 50,
				right: 50,
				bottom: 100,
			}}
		>
			<Grid container height="100%" width="100%">
				{suggestedItems.map((item) => {
					return (
						<Grid item xs={3}>
							<Box
								height="100%"
								width="100%"
								display="flex"
								justifyContent="center"
								alignItems="center"
								sx={{
									padding: '10px',
								}}
							>
								<Box
									display="flex"
									justifyContent="center"
									alignItems="center"
									sx={{
										padding: '10px',
										background: 'rgba(0, 0, 0, 0.4)',
										borderRadius: '10px',
									}}
								>
									<Item
										key={item.id}
										item={item}
										preferPreview={true}
										itemWidth={getItemSize().width}
										itemHeight={getItemSize().height}
									/>
								</Box>
							</Box>
						</Grid>
					);
				})}
			</Grid>
		</Stack>
	);
}

export default ItemSuggestions;
