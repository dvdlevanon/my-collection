import { useTheme } from '@emotion/react';
import CancelIcon from '@mui/icons-material/Cancel';
import StopIcon from '@mui/icons-material/Stop';
import { Box, IconButton, MenuItem, Select, Skeleton, Stack } from '@mui/material';
import React, { useState } from 'react';
import { useQuery } from 'react-query';
import Client from '../../utils/client';
import ReactQueryUtil from '../../utils/react-query-util';
import TagsUtil from '../../utils/tags-util';

function HighlightControls({ onCancel, onDone }) {
	const tagsQuery = useQuery(ReactQueryUtil.TAGS_KEY, Client.getTags);
	const [selectedHighlight, setSelectedHighlight] = useState(-1);
	const [open, setOpen] = useState(false);
	const theme = useTheme();

	const getFakeCategory = () => {
		return {
			id: -1,
			title: 'Select Highlight',
		};
	};

	const getHighlights = (tags) => {
		return TagsUtil.sortByTitle(tags.filter((cur) => TagsUtil.isHighlightsCategory(cur.parentId)));
	};

	return (
		<Stack
			flexDirection="column"
			sx={{
				gap: theme.spacing(1),
				background: '#000',
				padding: theme.multiSpacing(0.5, 1),
				opacity: '0.7',
				borderRadius: theme.spacing(1),
				position: 'absolute',
				right: theme.spacing(2),
				bottom: '100px',
			}}
		>
			{tagsQuery.isSuccess && (
				<Select
					open={open}
					onOpen={() => setOpen(true)}
					onClose={() => setOpen(false)}
					onChange={(e) => setSelectedHighlight(e.target.value)}
					size="small"
					value={selectedHighlight}
					displayEmpty
				>
					<MenuItem key={getFakeCategory().id} value={getFakeCategory().id}>
						{getFakeCategory().title}
					</MenuItem>
					{getHighlights(tagsQuery.data).map((highlight) => {
						return (
							<MenuItem key={highlight.id} value={highlight.id}>
								{highlight.title}
							</MenuItem>
						);
					})}
				</Select>
			)}
			<Stack flexDirection="row" gap={theme.spacing(1)} width="100%">
				<Box sx={{ padding: theme.spacing(0.9) }}>
					<Skeleton
						color={'red'}
						variant="circular"
						animation="pulse"
						width={theme.iconSize(1.2)}
						height={theme.iconSize(1.2)}
						sx={{
							backgroundColor: '#880000',
						}}
					/>
				</Box>
				<IconButton
					onClick={(e) => {
						onDone(selectedHighlight);
					}}
				>
					<StopIcon sx={{ fontSize: theme.iconSize(1.2) }} />
				</IconButton>
				<IconButton
					onClick={(e) => {
						onCancel();
					}}
					sx={{
						marginLeft: 'auto',
					}}
				>
					<CancelIcon sx={{ fontSize: theme.iconSize(1.2) }} />
				</IconButton>
			</Stack>
		</Stack>
	);
}

export default HighlightControls;
