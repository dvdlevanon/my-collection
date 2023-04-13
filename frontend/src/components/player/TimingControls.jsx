import LeftIcon from '@mui/icons-material/NavigateBefore';
import RightIcon from '@mui/icons-material/NavigateNext';
import TimingIcon from '@mui/icons-material/Schedule';
import { Fade, IconButton, Stack, Tooltip } from '@mui/material';
import React, { useState } from 'react';

function TimingControls({ showSchedule, setShowSchedule, setRelativeTime }) {
	let [changeTimerId, setChangeTimerId] = useState(0);
	let [pressStartedAt, setPressStartedAt] = useState(0);

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
			gap="20px"
			onMouseEnter={(e) => setShowSchedule(true)}
		>
			{showSchedule && (
				<Fade in={showSchedule}>
					<Stack flexDirection="row">
						<Tooltip title="-1s">
							<IconButton
								size="small"
								onMouseDown={() => onMouseDown(-1)}
								onMouseUp={onMouseUp}
								onMouseLeave={onMouseLeave}
							>
								<LeftIcon sx={{ fontSize: '20px' }} />
							</IconButton>
						</Tooltip>
						<Tooltip title="-10s">
							<IconButton
								size="small"
								onMouseDown={() => onMouseDown(-10)}
								onMouseUp={onMouseUp}
								onMouseLeave={onMouseLeave}
							>
								<LeftIcon sx={{ fontSize: '25px' }} />
							</IconButton>
						</Tooltip>
						<Tooltip title="-30s">
							<IconButton
								size="small"
								onMouseDown={() => onMouseDown(-30)}
								onMouseUp={onMouseUp}
								onMouseLeave={onMouseLeave}
							>
								<LeftIcon sx={{ fontSize: '30px' }} />
							</IconButton>
						</Tooltip>
						<Tooltip title="+30s">
							<IconButton
								size="small"
								onMouseDown={() => onMouseDown(30)}
								onMouseUp={onMouseUp}
								onMouseLeave={onMouseLeave}
							>
								<RightIcon sx={{ fontSize: '30px' }} />
							</IconButton>
						</Tooltip>
						<Tooltip title="+10s">
							<IconButton
								size="small"
								onMouseDown={() => onMouseDown(10)}
								onMouseUp={onMouseUp}
								onMouseLeave={onMouseLeave}
							>
								<RightIcon sx={{ fontSize: '25px' }} />
							</IconButton>
						</Tooltip>
						<Tooltip title="+1s">
							<IconButton
								size="small"
								onMouseDown={() => onMouseDown(1)}
								onMouseUp={onMouseUp}
								onMouseLeave={onMouseLeave}
							>
								<RightIcon sx={{ fontSize: '20px' }} />
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
