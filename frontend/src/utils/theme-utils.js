import { createTheme } from '@mui/material';

export default class ThemeUtil {
	static extendTheme(theme) {
		theme.multiSpacing = (...values) => {
			return values
				.map((value) => {
					return theme.spacing(value);
				})
				.join(' ');
		};

		theme.iconSize = (factor) => {
			let val = theme.iconBaseSize * factor;
			return val + 'px';
		};

		theme.border = (factor, style, color) => {
			let val = theme.borderBaseSize * factor;
			return val + 'px ' + style + ' ' + color;
		};

		theme.fontSize = (factor) => {
			let val = theme.typography.fontSize * factor;
			return val + 'px ';
		};

		return theme;
	}

	static createDarkTheme() {
		let theme = createTheme({
			spacing: 10,
			baseSpacing: 10,
			palette: {
				mode: 'dark',
				primary: {
					main: '#ff4400',
					light: '#ff8844',
				},
				secondary: {
					main: '#0D9352',
				},
				bright: {
					main: '#ffddcc',
					darker: '#ddbbaa',
					darker2: '#aaaaaa',
					text: '#ffffff',
				},
				dark: {
					main: '#121212',
					lighter: '#222222',
					lighter2: '#272727',
				},
			},
			typography: {
				fontSize: 16,
			},
			iconBaseSize: 25,
			borderBaseSize: 1,
		});

		return this.extendTheme(theme);
	}
}
