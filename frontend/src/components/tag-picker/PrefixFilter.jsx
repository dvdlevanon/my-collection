import { useTheme } from '@emotion/react';
import { Box, Stack } from '@mui/material';
import React from 'react';

function PrefixFilter({ selectedChar, setSelectedChar }) {
	const theme = useTheme();

	const numberToChar = (num) => {
		return String.fromCharCode(num + 'A'.charCodeAt(0));
	};

	const charToNumber = (c) => {
		return c.charCodeAt(0);
	};

	const englishLettersCount = () => {
		return charToNumber('Z') - charToNumber('A') + 1;
	};

	const isSelectedNumber = (i) => {
		return selectedChar == numberToChar(i);
	};

	return (
		<Stack
			padding={theme.spacing(0.3)}
			flexDirection="row"
			justifyContent="space-around"
			alignSelf="center"
			width="50%"
		>
			{Array.from(Array(englishLettersCount()).keys()).map((i) => {
				return (
					<Box
						key={i}
						color="primary"
						onClick={(e) => {
							if (isSelectedNumber(i)) {
								setSelectedChar('');
							} else {
								setSelectedChar(numberToChar(i));
							}
						}}
						sx={{
							color: isSelectedNumber(i) ? 'primary.main' : 'auto',
							cursor: 'pointer',
							userSelect: 'none',
							width: '100%',
							textAlign: 'center',
							'&:hover': {
								textDecoration: 'underline',
							},
						}}
					>
						{numberToChar(i)}
					</Box>
				);
			})}
		</Stack>
	);
}

export default PrefixFilter;
