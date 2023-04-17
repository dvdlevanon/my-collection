import KeyboardArrowUpIcon from '@mui/icons-material/KeyboardArrowUp';
import { Fab, Fade, Stack, useScrollTrigger } from '@mui/material';
import { Box } from '@mui/system';
import React from 'react';
import Item from './Item';

function ItemsList({ itemsSize, items, previewMode, itemLinkBuilder }) {
	const trigger = useScrollTrigger({
		disableHysteresis: true,
		threshold: 100,
	});

	const handleClick = (event) => {
		const anchor = document.querySelector('#back-to-top-anchor');

		if (anchor) {
			anchor.scrollIntoView({
				block: 'center',
			});
		}
	};

	const onConvertVideo = (item) => {};

	const onConvertAudio = (item) => {};

	return (
		<React.Fragment>
			<div id="back-to-top-anchor" />
			<Stack flexDirection="row" flexWrap="wrap" gap="20px" padding="20px">
				{items.map((item) => {
					return (
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
					);
				})}
			</Stack>
			<Fade in={trigger}>
				<Box onClick={handleClick} role="presentation" sx={{ position: 'fixed', bottom: 16, right: 16 }}>
					<Fab size="small" aria-label="scroll back to top">
						<KeyboardArrowUpIcon />
					</Fab>
				</Box>
			</Fade>
		</React.Fragment>
	);
}

export default ItemsList;
