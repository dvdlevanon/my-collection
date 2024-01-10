import { Box } from '@mui/system';
import React from 'react';
import { AutoSizer, Grid } from 'react-virtualized';
import Item from './Item';

function ItemsList({ itemsSize, items, previewMode, itemLinkBuilder, onScroll }) {
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
						padding: '20px',
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
						onScroll={onScroll}
					/>
				)}
			</AutoSizer>
		</React.Fragment>
	);
}

export default ItemsList;
