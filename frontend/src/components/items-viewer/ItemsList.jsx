import { useTheme } from '@emotion/react';
import KeyboardArrowUpIcon from '@mui/icons-material/KeyboardArrowUp';
import { Fab, Fade } from '@mui/material';
import { Box } from '@mui/system';
import React, { useState } from 'react';
import { AutoSizer, Grid } from 'react-virtualized';
import Item from './Item';

function ItemsList({ itemsSize, items, previewMode, itemLinkBuilder, onScroll }) {
	const [scrollTop, setScrollTop] = useState(0);
	const theme = useTheme();

	const onConvertVideo = (item) => {};

	const onConvertAudio = (item) => {};

	const calcColumnsCount = (width) => {
		return Math.floor(width / (itemsSize.width + 20));
	};

	const cellRenderer = (width) => {
		return ({ columnIndex, key, rowIndex, style }) => {
			let itemIndex = rowIndex * calcColumnsCount(width) + columnIndex;
			if (itemIndex >= items.length) {
				return <div key={key}></div>;
			}

			let item = items[itemIndex];
			return (
				<Box
					key={item.id}
					style={style}
					sx={{
						padding: theme.spacing(2),
						width: itemsSize.width,
						height: itemsSize.height,
					}}
				>
					<Item
						key={item.id}
						item={item}
						preferPreview={previewMode}
						onConvertAudio={onConvertAudio}
						onConvertVideo={onConvertVideo}
						itemWidth={itemsSize.width}
						itemHeight={itemsSize.height}
						direction="column"
						itemLinkBuilder={itemLinkBuilder}
						withItemTitleMenu={false}
					/>
				</Box>
			);
		};
	};

	return (
		<React.Fragment>
			<AutoSizer>
				{({ height, width }) => (
					<Grid
						cellRenderer={cellRenderer(width)}
						columnCount={calcColumnsCount(width)}
						rowCount={Math.floor(items.length / calcColumnsCount(width) + 1)}
						columnWidth={itemsSize.width + 40}
						rowHeight={itemsSize.height + 100}
						height={height}
						width={width}
						scrollTop={scrollTop}
						onScroll={(e) => {
							setScrollTop(e.scrollTop);
							onScroll(e);
						}}
					/>
				)}
			</AutoSizer>
			<Fade in={scrollTop > 100}>
				<Box
					onClick={() => {
						setScrollTop(0);
					}}
					role="presentation"
					sx={{ position: 'fixed', bottom: 16, right: 16 }}
				>
					<Fab size="small">
						<KeyboardArrowUpIcon sx={{ fontSize: theme.iconSize(1) }} />
					</Fab>
				</Box>
			</Fade>
		</React.Fragment>
	);
}

export default ItemsList;
