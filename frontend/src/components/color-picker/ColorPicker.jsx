import DoneIcon from '@mui/icons-material/Done';
import ColorIcon from '@mui/icons-material/Palette';
import { Avatar, Box, Divider, IconButton, Popover, Stack, useTheme } from '@mui/material';
import { useRef, useState } from 'react';

import { Sketch } from '@uiw/react-color';
import ColorUtil from '../../utils/color-utils';

function ColorPicker({ color, onChange }) {
	const theme = useTheme();
	const [pickerOpenend, setPickerOpened] = useState(false);
	const anochorEl = useRef();

	return (
		<Box ref={anochorEl}>
			<IconButton onClick={() => setPickerOpened(true)}>
				<Avatar
					sx={{
						bgcolor: color,
						border: '1px solid rgba(255, 255, 255, 0.6)',
						borderRadius: '50%',
					}}
				>
					<ColorIcon sx={{ color: ColorUtil.getInvertedColor(color) }} />
				</Avatar>
			</IconButton>
			<Popover
				container={anochorEl.current}
				anchorEl={anochorEl.current}
				open={pickerOpenend}
				onClose={() => {
					setPickerOpened(false);
				}}
				sx={{
					zIndex: 10000,
				}}
			>
				<Stack padding={theme.spacing(1)}>
					<Sketch
						color={color}
						onChange={(c) => {
							onChange(c.hex);
						}}
					/>
					<Divider />
					<Stack
						flexDirection={'row'}
						alignItems={'center'}
						padding={theme.spacing(1)}
						gap={theme.spacing(1)}
						justifyContent={'center'}
					>
						<IconButton
							onClick={() => {
								setPickerOpened(false);
							}}
						>
							<DoneIcon sx={{ fontSize: theme.iconSize(1) }} />
						</IconButton>
					</Stack>
				</Stack>
			</Popover>
		</Box>
	);
}

export default ColorPicker;
