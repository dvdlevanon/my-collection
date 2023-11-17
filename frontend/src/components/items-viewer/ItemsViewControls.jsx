import ZoomInIcon from '@mui/icons-material/ZoomIn';
import ZoomOutIcon from '@mui/icons-material/ZoomOut';
import { IconButton, Stack, ToggleButton, ToggleButtonGroup } from '@mui/material';
import React from 'react';
import AspectRatioUtil from '../../utils/aspect-ratio-util';

function ItemsViewControls({ itemsSize, onZoomChanged, aspectRatio, onAspectRatioChanged }) {
	return (
		<Stack flexDirection="row" gap="10px" padding="10px">
			<Stack justifyContent="center" alignContent="center">
				<ToggleButtonGroup
					size="small"
					exclusive
					value={aspectRatio}
					onChange={(e, newValue) => {
						if (newValue != null) {
							onAspectRatioChanged(newValue);
						}
					}}
				>
					<ToggleButton value={AspectRatioUtil.asepctRatio16_9}>
						{AspectRatioUtil.toString(AspectRatioUtil.asepctRatio16_9)}
					</ToggleButton>
					<ToggleButton value={AspectRatioUtil.asepctRatio4_3}>
						{AspectRatioUtil.toString(AspectRatioUtil.asepctRatio4_3)}
					</ToggleButton>
					<ToggleButton value={AspectRatioUtil.asepctRatio4_2}>
						{AspectRatioUtil.toString(AspectRatioUtil.asepctRatio4_2)}
					</ToggleButton>
				</ToggleButtonGroup>
			</Stack>
			<IconButton disabled={itemsSize.width <= 100} onClick={() => onZoomChanged(-50)}>
				<ZoomOutIcon />
			</IconButton>
			<IconButton onClick={() => onZoomChanged(50)}>
				<ZoomInIcon />
			</IconButton>
		</Stack>
	);
}

export default ItemsViewControls;
