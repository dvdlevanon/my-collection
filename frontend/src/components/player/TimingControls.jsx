import LeftIcon from '@mui/icons-material/NavigateBefore';
import RightIcon from '@mui/icons-material/NavigateNext';
import TimingIcon from '@mui/icons-material/Schedule';
import { Fade, IconButton, Stack, Tooltip } from '@mui/material';
import React from 'react';

function TimingControls({ showSchedule, setShowSchedule, setRelativeTime }) {
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
						<Tooltip title="-250ms">
							<IconButton size="small" onMouseDown={() => setRelativeTime(-0.025)}>
								<LeftIcon sx={{ fontSize: '20px' }} />
							</IconButton>
						</Tooltip>
						<Tooltip title="-1s">
							<IconButton size="small" onClick={() => setRelativeTime(-1)}>
								<LeftIcon sx={{ fontSize: '25px' }} />
							</IconButton>
						</Tooltip>
						<Tooltip title="-10s">
							<IconButton size="small" onClick={() => setRelativeTime(-10)}>
								<LeftIcon sx={{ fontSize: '30px' }} />
							</IconButton>
						</Tooltip>
						<Tooltip title="+1000ms">
							<IconButton size="small" onClick={() => setRelativeTime(10)}>
								<RightIcon sx={{ fontSize: '30px' }} />
							</IconButton>
						</Tooltip>
						<Tooltip title="+1s">
							<IconButton size="small" onClick={() => setRelativeTime(1)}>
								<RightIcon sx={{ fontSize: '25px' }} />
							</IconButton>
						</Tooltip>
						<Tooltip title="+250ms">
							<IconButton size="small" onClick={() => setRelativeTime(0.025)}>
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
