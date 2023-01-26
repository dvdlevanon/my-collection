import KeyboardArrowUpIcon from '@mui/icons-material/KeyboardArrowUp';
import { Fab, Fade, useScrollTrigger } from '@mui/material';
import { Box } from '@mui/system';
import React from 'react';
import Item from './Item';
import styles from './Items.module.css';

function ItemsList({ items, previewMode }) {
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

	return (
		<React.Fragment>
			<div id="back-to-top-anchor" />
			<Box className={styles.items}>
				{items.map((item) => {
					return <Item key={item.id} item={item} preferPreview={previewMode} />;
				})}
			</Box>
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
