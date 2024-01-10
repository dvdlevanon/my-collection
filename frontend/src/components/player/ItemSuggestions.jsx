import { useTheme } from '@emotion/react';
import { Box, Grid, Stack } from '@mui/material';
import React from 'react';
import AspectRatioUtil from '../../utils/aspect-ratio-util';
import Item from '../items-viewer/Item';

function ItemSuggestions({ suggestedItems, width, onBackgroundClick, onBackgroundDoubleClick }) {
	const theme = useTheme();
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
			onClick={onBackgroundClick}
			onDoubleClick={onBackgroundDoubleClick}
			sx={{
				position: 'absolute',
				left: theme.spacing(5),
				top: theme.spacing(5),
				right: theme.spacing(5),
				bottom: theme.spacing(10),
			}}
		>
			<Grid container height="100%" width="100%">
				{suggestedItems.map((item) => {
					return (
						<Grid item xs={3} key={item.id}>
							<Box
								height="100%"
								width="100%"
								display="flex"
								justifyContent="center"
								alignItems="center"
								onClick={onBackgroundClick}
								onDoubleClick={onBackgroundDoubleClick}
								padding={theme.spacing(1)}
							>
								<Box
									display="flex"
									justifyContent="center"
									alignItems="center"
									onClick={onBackgroundClick}
									onDoubleClick={onBackgroundDoubleClick}
									padding={theme.spacing(1)}
									sx={{
										background: 'rgba(0, 0, 0, 0.4)',
										borderRadius: theme.spacing(1),
									}}
								>
									<Item
										key={item.id}
										item={item}
										preferPreview={true}
										itemWidth={getItemSize().width}
										itemHeight={getItemSize().height}
										direction="column"
										withItemTitleMenu={false}
										itemLinkBuilder={(item) => {
											return '/spa/item/' + item.id + window.location.search;
										}}
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
