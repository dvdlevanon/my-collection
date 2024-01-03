import KeyboardIcon from '@mui/icons-material/Keyboard';
import { Box, InputAdornment, Popper, TextField } from '@mui/material';
import React, { useRef, useState } from 'react';
import Keyboard from 'react-simple-keyboard';
import 'react-simple-keyboard/build/css/index.css';

function TextFieldWithKeyboard(props) {
	const [showKeyboard, setShowKeyboard] = useState(false);
	const [layoutName, setLayoutName] = useState('default');
	const textEl = useRef(null);
	const [text, setText] = useState('');

	const onChange = (value) => {
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
							<KeyboardIcon />
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
				<Keyboard
					theme={'hg-theme-default dark'}
					layoutName={layoutName}
					value={text}
					onChange={onChange}
					onKeyPress={onKeyPress}
				/>
			</Popper>
		</Box>
	);
}

export default TextFieldWithKeyboard;
