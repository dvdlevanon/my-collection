import { useTheme } from '@emotion/react';
import KeyboardIcon from '@mui/icons-material/Keyboard';
import { Box, ClickAwayListener, InputAdornment, Popper, TextField } from '@mui/material';
import React, { useEffect, useRef, useState } from 'react';
import Keyboard from 'react-simple-keyboard';
import 'react-simple-keyboard/build/css/index.css';

function TextFieldWithKeyboard(props) {
	const [showKeyboard, setShowKeyboard] = useState(false);
	const [layoutName, setLayoutName] = useState('default');
	const [text, setText] = useState('');
	const textEl = useRef(null);
	const keyboard = useRef(null);
	const theme = useTheme();

	useEffect(() => {
		if (keyboard.current) {
			keyboard.current.setInput(text);
		}
	}, [text, showKeyboard, keyboard.current]);

	const onChange = (value) => {
		if (value == text) {
			return;
		}

		setText(value);
		props.onChange(value);
	};

	const onKeyPress = (e) => {
		if (e === '{shift}' || e === '{lock}') {
			setLayoutName(layoutName === 'default' ? 'shift' : 'default');
		}
	};

	return (
		<Box>
			<TextField
				ref={textEl}
				{...props}
				value={text}
				onChange={(e) => onChange(e.target.value)}
				InputProps={{
					endAdornment: (
						<InputAdornment
							position="end"
							sx={{
								cursor: 'pointer',
							}}
							onClick={() => setShowKeyboard(!showKeyboard)}
						>
							<KeyboardIcon sx={{ fontSize: theme.iconSize(1) }} />
						</InputAdornment>
					),
				}}
			></TextField>
			<Popper
				open={showKeyboard}
				anchorEl={textEl.current}
				sx={{
					zIndex: 10000,
				}}
			>
				<ClickAwayListener onClickAway={(e) => setShowKeyboard(false)}>
					<Box>
						<Keyboard
							keyboardRef={(r) => (keyboard.current = r)}
							theme={'hg-theme-default dark'}
							layoutName={layoutName}
							value={text}
							onChange={onChange}
							onKeyPress={onKeyPress}
						/>
					</Box>
				</ClickAwayListener>
			</Popper>
		</Box>
	);
}

export default TextFieldWithKeyboard;
