import DarkModeIcon from '@mui/icons-material/DarkMode';
import LightModeIcon from '@mui/icons-material/LightMode';
import { Stack, ToggleButton, ToggleButtonGroup, Tooltip } from '@mui/material';
import React from 'react';
import ThemeUtil from '../../utils/theme-utils';

function ThemeSelector({ theme, setTheme }) {
	return (
		<Stack flexDirection="row" gap={theme.spacing(1)}>
			<ToggleButtonGroup
				size="small"
				exclusive
				value={theme.name}
				onChange={(e, newValue) => {
					if (newValue != null) {
						setTheme(ThemeUtil.themeByName(newValue));
					}
				}}
			>
				<Tooltip title="Orange Dark Mode" value="dark-orange">
					<ToggleButton value="dark-orange">
						<DarkModeIcon sx={{ fontSize: theme.iconSize(1), color: '#ff4400' }} />
					</ToggleButton>
				</Tooltip>
				<Tooltip title="Purple Light Mode" value="dark-purple">
					<ToggleButton value="dark-purple">
						<DarkModeIcon sx={{ fontSize: theme.iconSize(1), color: '#BB00FF' }} />
					</ToggleButton>
				</Tooltip>
				<Tooltip title="Blue Light Mode" value="light-blue">
					<ToggleButton value="light-blue">
						<LightModeIcon sx={{ fontSize: theme.iconSize(1), color: '#009DFF' }} />
					</ToggleButton>
				</Tooltip>
				<Tooltip title="Blue Light Mode" value="light-pink">
					<ToggleButton value="light-pink">
						<LightModeIcon sx={{ fontSize: theme.iconSize(1), color: '#FF0072' }} />
					</ToggleButton>
				</Tooltip>
			</ToggleButtonGroup>
		</Stack>
	);
}

export default ThemeSelector;
