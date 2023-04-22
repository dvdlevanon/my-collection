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
			gap="10px"
			sx={{
				background: '#000',
				padding: '3px 10px',
				opacity: '0.7',
				borderRadius: '10px',
				position: 'absolute',
				right: 20,
				bottom: 100,
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
			<Stack flexDirection="row" gap="10px" width="100%">
				<Box sx={{ padding: '9px' }}>
					<Skeleton
						color={'red'}
						variant="circular"
						animation="pulse"
						width={25}
						height={25}
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
					<StopIcon />
				</IconButton>
				<IconButton
					onClick={(e) => {
						onCancel();
					}}
					sx={{
						marginLeft: 'auto',
					}}
				>
					<CancelIcon />
				</IconButton>
			</Stack>
		</Stack>
	);
}

export default HighlightControls;
