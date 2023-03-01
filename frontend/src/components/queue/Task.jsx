import PendingIcon from '@mui/icons-material/Pending';
import { Box, CircularProgress, Tooltip, Typography } from '@mui/material';
import { Stack } from '@mui/system';
import React from 'react';
import TasksUtil from '../../utils/tasks-util';
import TimeUtil from '../../utils/time-utils';

function Task({ task }) {
	return (
		<Box>
			{(!Boolean(task) && 'No Tasks') || (
				<Stack flexDirection="row" gap="10px" alignItems="center">
					{(TasksUtil.isProcessing(task) && (
						<Tooltip title="Processing">
							<CircularProgress color="bright" size="25px" />
						</Tooltip>
					)) || (
						<Tooltip title="Pending">
							<PendingIcon color="bright" size="25px" />
						</Tooltip>
					)}
					<Stack flexDirection="column">
						<Tooltip title={task.title}>
							<Typography variant="caption" maxWidth={500} noWrap>
								{task.description}
							</Typography>
						</Tooltip>
						{(TasksUtil.isProcessing(task) && (
							<Typography variant="caption" color="secondary" sx={{ fontStyle: 'italic' }}>
								Started 5 minutes ago
							</Typography>
						)) || (
							<Typography variant="caption" color="bright.main" sx={{ fontStyle: 'italic' }}>
								Pending for {TimeUtil.msToTime(Date.now() - task.enqueueTime)} minutes
							</Typography>
						)}
					</Stack>
				</Stack>
			)}
		</Box>
	);
}

export default Task;
