import CloseIcon from '@mui/icons-material/Close';
import PauseIcon from '@mui/icons-material/Pause';
import PlayIcon from '@mui/icons-material/PlayArrow';
import {
	Box,
	CircularProgress,
	Divider,
	IconButton,
	List,
	ListItem,
	Pagination,
	Paper,
	Stack,
	Typography,
} from '@mui/material';
import React from 'react';
import { useQuery, useQueryClient } from 'react-query';
import Client from '../../utils/client';
import ReactQueryUtil from '../../utils/react-query-util';
import Task from './Task';

function Queue({ onClose }) {
	const queryClient = useQueryClient();
	const queueMetadataQuery = useQuery(ReactQueryUtil.QUEUE_METADATA_KEY, Client.getQueueMetadata);
	const [tasksPage, setTasksPage] = React.useState(1);
	const [tasksPageSize, setTasksPageSize] = React.useState(10);
	const tasksQuery = useQuery({
		queryKey: ReactQueryUtil.tasksPageKey(tasksPage, tasksPageSize),
		queryFn: () => Client.getTasks(tasksPage, tasksPageSize),
		keepPreviousData: true,
	});

	const toggleProcessing = () => {
		if (queueMetadataQuery.data.paused) {
			Client.continueProcessingTasks();
		} else {
			Client.pauseProcessingTasks();
		}
	};

	return (
		<Paper elevation={3} sx={{ padding: '10px' }}>
			<Stack flexDirection="column">
				<Stack flexDirection="row" alignItems="center">
					{(queueMetadataQuery.isSuccess && (
						<Typography variant="body1" sx={{ paddingRight: '40px', flexGrow: '1' }}>
							Processing {queueMetadataQuery.data.size} tasks
						</Typography>
					)) || (
						<Box flexGrow={1}>
							<CircularProgress color="bright" size="25px" />
						</Box>
					)}
					{queueMetadataQuery.isSuccess && (
						<IconButton onClick={toggleProcessing}>
							{console.log(queueMetadataQuery.data.paused)}
							{(queueMetadataQuery.data.paused && <PlayIcon />) || <PauseIcon />}
						</IconButton>
					)}
					<IconButton onClick={onClose}>
						<CloseIcon />
					</IconButton>
				</Stack>
				<Divider />
				<List>
					{(!tasksQuery.isSuccess || tasksQuery.data.length == 0) && (
						<ListItem>
							<Task task={null} />
						</ListItem>
					)}
					{tasksQuery.isSuccess &&
						tasksQuery.data.length != 0 &&
						tasksQuery.data.map((task) => {
							return (
								<ListItem key={task.id} sx={{ backgroundColor: 'dark.lighter' }}>
									<Task task={task} />
								</ListItem>
							);
						})}
				</List>
				{queueMetadataQuery.isSuccess && (
					<Pagination
						count={Math.ceil(queueMetadataQuery.data.size / tasksPageSize)}
						page={tasksPage}
						onChange={(e, page) => setTasksPage(page)}
						sx={{ alignSelf: 'center' }}
					/>
				)}
			</Stack>
		</Paper>
	);
}

export default Queue;
