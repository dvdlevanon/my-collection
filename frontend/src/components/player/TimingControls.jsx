import { useTheme } from '@emotion/react';
import LeftIcon from '@mui/icons-material/NavigateBefore';
import RightIcon from '@mui/icons-material/NavigateNext';
import TimingIcon from '@mui/icons-material/Schedule';
import { Fade, IconButton, Stack, Tooltip } from '@mui/material';
import { useState } from 'react';
import { usePlayerStore } from './PlayerStore';

function TimingControls({ setRelativeTime }) {
	const playerStore = usePlayerStore();
	const [changeTimerId, setChangeTimerId] = useState(0);
	const [pressStartedAt, setPressStartedAt] = useState(0);
	const theme = useTheme();

	const onMouseDown = (offset) => {
		setRelativeTime(offset);
		setPressStartedAt(Date.now());
		installLongPressTimer(offset);
	};

	const installLongPressTimer = (offset) => {
		let millisSinceStart = pressStartedAt == 0 ? 0 : Date.now() - pressStartedAt;
		console.log(millisSinceStart);
		setChangeTimerId(
			setTimeout(() => {
				setRelativeTime(offset);
				installLongPressTimer(offset);
			}, Math.max(250 - (millisSinceStart / 1000) * 10, 10))
		);
	};

	const clearLongPress = () => {
		clearTimeout(changeTimerId);
		setChangeTimerId(0);
	};

	const onMouseUp = () => {
		clearLongPress();
	};

	const onMouseLeave = () => {
		clearLongPress();
	};

	return (
		<Stack
			display="flex"
			flexDirection="row"
			alignItems="center"
			gap={theme.spacing(2)}
			onMouseEnter={(e) => playerStore.setShowSchedule(true)}
		>
			{playerStore.showSchedule && (
				<Fade in={playerStore.showSchedule}>
					<Stack flexDirection="row">
						<Tooltip title="-1s">
							<IconButton
								size="small"
								onMouseDown={() => onMouseDown(-1)}
								onMouseUp={onMouseUp}
								onMouseLeave={onMouseLeave}
							>
								<LeftIcon sx={{ fontSize: theme.iconSize(0.7) }} />
							</IconButton>
						</Tooltip>
						<Tooltip title="-10s">
							<IconButton
								size="small"
								onMouseDown={() => onMouseDown(-10)}
								onMouseUp={onMouseUp}
								onMouseLeave={onMouseLeave}
							>
								<LeftIcon sx={{ fontSize: theme.iconSize(1) }} />
							</IconButton>
						</Tooltip>
						<Tooltip title="-30s">
							<IconButton
								size="small"
								onMouseDown={() => onMouseDown(-30)}
								onMouseUp={onMouseUp}
								onMouseLeave={onMouseLeave}
							>
								<LeftIcon sx={{ fontSize: theme.iconSize(1.3) }} />
							</IconButton>
						</Tooltip>
						<Tooltip title="+30s">
							<IconButton
								size="small"
								onMouseDown={() => onMouseDown(30)}
								onMouseUp={onMouseUp}
								onMouseLeave={onMouseLeave}
							>
								<RightIcon sx={{ fontSize: theme.iconSize(1.3) }} />
							</IconButton>
						</Tooltip>
						<Tooltip title="+10s">
							<IconButton
								size="small"
								onMouseDown={() => onMouseDown(10)}
								onMouseUp={onMouseUp}
								onMouseLeave={onMouseLeave}
							>
								<RightIcon sx={{ fontSize: theme.iconSize(1) }} />
							</IconButton>
						</Tooltip>
						<Tooltip title="+1s">
							<IconButton
								size="small"
								onMouseDown={() => onMouseDown(1)}
								onMouseUp={onMouseUp}
								onMouseLeave={onMouseLeave}
							>
								<RightIcon sx={{ fontSize: theme.iconSize(0.7) }} />
							</IconButton>
						</Tooltip>
					</Stack>
				</Fade>
			)}
			<Tooltip title={'Timing controls'}>
				<IconButton>
					<TimingIcon />
				</IconButton>
			</Tooltip>
		</Stack>
	);
}

export default TimingControls;
