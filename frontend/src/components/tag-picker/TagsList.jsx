import { useTheme } from '@emotion/react';
import KeyboardArrowUpIcon from '@mui/icons-material/KeyboardArrowUp';
import { Box, Fab, Fade } from '@mui/material';
import React, { useState } from 'react';
import { AutoSizer, Grid } from 'react-virtualized';
import Tag from './Tag';

function TagsList({ tags, tagsSize, onTagClicked, tagLinkBuilder, tit, parent, onScroll }) {
	const [scrollTop, setScrollTop] = useState(0);
	const theme = useTheme();

	const calcColumnsCount = (width) => {
		return Math.floor(width / (tagsSize.width + 20));
	};

	const cellRenderer = (width) => {
		return ({ columnIndex, key, rowIndex, style }) => {
			let itemIndex = rowIndex * calcColumnsCount(width) + columnIndex;
			if (itemIndex >= tags.length) {
				return <div key={key}></div>;
			}

			let tag = tags[itemIndex];
			return (
				<Box
					key={tag.id}
					style={style}
					sx={{
						padding: theme.spacing(2),
						width: tagsSize.width,
						height: tagsSize.height,
					}}
				>
					<Tag
						key={tag.id}
						tag={tag}
						parent={parent}
						tagDimension={tagsSize}
						selectedTit={tit}
						tagLinkBuilder={tagLinkBuilder}
						onTagClicked={onTagClicked}
					/>
				</Box>
			);
		};
	};

	return (
		<>
			<AutoSizer>
				{({ height, width }) => (
					<Grid
						cellRenderer={cellRenderer(width)}
						columnCount={calcColumnsCount(width)}
						rowCount={Math.floor(tags.length / calcColumnsCount(width) + 1)}
						columnWidth={tagsSize.width + 20}
						rowHeight={tagsSize.height + 20}
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
		</>
	);
}

export default TagsList;
