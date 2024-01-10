import { useTheme } from '@emotion/react';
import DoneIcon from '@mui/icons-material/Done';
import PendingIcon from '@mui/icons-material/Pending';
import { Box, CircularProgress, Tooltip, Typography } from '@mui/material';
import { Stack } from '@mui/system';
import React, { useEffect, useReducer } from 'react';
import TasksUtil from '../../utils/tasks-util';
import TimeUtil from '../../utils/time-utils';

function Task({ task }) {
	const [ignored, forceUpdate] = useReducer((x) => x + 1, 0);
	const theme = useTheme();

	useEffect(() => {
		setInterval(forceUpdate, 1000);
	}, []);

	return (
		<Box>
			{(!Boolean(task) && 'No Tasks') || (
				<Stack flexDirection="row" gap={theme.spacing(1)} alignItems="center">
					{TasksUtil.isProcessing(task) && (
						<Tooltip title="Processing">
							<CircularProgress color="bright" size={theme.iconSize(1)} />
						</Tooltip>
					)}
					{TasksUtil.isDone(task) && (
						<Tooltip title="Done">
							<DoneIcon sx={{ fontSize: theme.iconSize(1) }} />
						</Tooltip>
					)}
					{TasksUtil.isPending(task) && (
						<Tooltip title="Pending">
							<PendingIcon color="bright" size={theme.iconSize(1)} />
						</Tooltip>
					)}
					<Stack flexDirection="column">
						<Tooltip title={task.description}>
							<Typography variant="caption" maxWidth={500} noWrap>
								{task.description}
							</Typography>
						</Tooltip>
						{TasksUtil.isProcessing(task) && (
							<Typography variant="caption" color="secondary" sx={{ fontStyle: 'italic' }}>
								Started {TimeUtil.msToTime(Date.now() - task.processingStart)} ago
							</Typography>
						)}
						{TasksUtil.isDone(task) && (
							<Typography variant="caption" color="bright.main" sx={{ fontStyle: 'italic' }}>
								Done in {TimeUtil.msToTime(task.processingEnd - task.processingStart)}
							</Typography>
						)}
						{TasksUtil.isPending(task) && (
							<Typography variant="caption" color="bright.main" sx={{ fontStyle: 'italic' }}>
								Pending for {TimeUtil.msToTime(Date.now() - task.enqueueTime)}
							</Typography>
						)}
					</Stack>
				</Stack>
			)}
		</Box>
	);
}

export default Task;
