import CloseIcon from '@mui/icons-material/Close';
import { IconButton, Stack, Tab, Tabs, useTheme } from '@mui/material';
import { useState } from 'react';
import SubtitlesAppearanceControls from './SubtitlesAppearanceControls';
import SubtitlesFinder from './SubtitlesFinder';
import SubtitlesSettings from './SubtitlesSettings';
import { useSubtitleStore } from './SubtitlesStore';

function SubtitlesControls() {
	const theme = useTheme();
	const subtitleStore = useSubtitleStore();
	const [selectedTab, setSelectedTab] = useState('settings');

	const handleChange = (event, newValue) => {
		setSelectedTab(newValue);
	};

	return (
		subtitleStore.controlsShown && (
			<Stack
				flexDirection="column"
				sx={{
					gap: theme.spacing(1),
					background: theme.palette.gradient.color2,
					padding: theme.multiSpacing(0.5, 1),
					opacity: '0.7',
					borderRadius: theme.spacing(1),
					position: 'absolute',
					right: theme.spacing(2),
					bottom: '100px',
				}}
			>
				<Stack flexDirection={'row'}>
					<IconButton
						onClick={(e) => {
							e.preventDefault();
							e.stopPropagation();
							subtitleStore.hideSubtitlesControls();
						}}
					>
						<CloseIcon sx={{ fontSize: theme.iconSize(1) }} />
					</IconButton>
					<Tabs value={selectedTab} onChange={handleChange}>
						<Tab label="Settings" value="settings" />
						<Tab label="Finder" value="finder" />
						<Tab label="Appearnace" value="appearnace" />
					</Tabs>
				</Stack>
				{selectedTab == 'settings' && <SubtitlesSettings />}
				{selectedTab == 'finder' && <SubtitlesFinder />}
				{selectedTab == 'appearnace' && <SubtitlesAppearanceControls />}
			</Stack>
		)
	);
}

export default SubtitlesControls;
