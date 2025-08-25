import { useTheme } from '@emotion/react';
import { Divider, Popover } from '@mui/material';
import { Box } from '@mui/system';
import { useQuery } from '@tanstack/react-query';
import { useEffect } from 'react';
import ReactQueryUtil from '../../utils/react-query-util';
import CreateAnnotation from './CreateAnnotation';
import ExistingAnnotations from './ExistingAnnotations';
import PopoverHeader from './PopoverHeader';
import { useTagAnnotationsStore } from './tagAnnotationsStore';

function ManageTagAnnotations({ tag, menu, onClose }) {
	const tagAnnotationsStore = useTagAnnotationsStore();

	const theme = useTheme();
	const availableAnnotations = useQuery(ReactQueryUtil.availableAnnotationsQuery(tag.parentId));

	useEffect(() => {
		tagAnnotationsStore.setTag(tag);
	}, []);

	useEffect(() => {
		tagAnnotationsStore.setAvailableAnnotations(availableAnnotations.data);
	}, [availableAnnotations.data]);

	return (
		<Popover
			onClose={(e, reason) => {
				if (reason == 'backdropClick' || reason == 'escapeKeyDown') {
					onClose(e);
				}
			}}
			open={true}
			anchorReference="anchorPosition"
			anchorPosition={{ top: menu.mouseY, left: menu.mouseX }}
			BackdropProps={{
				open: true,
				invisible: false,
				onClick: (e) => {
					e.preventDefault();
					e.stopPropagation();
					onClose(e);
				},
			}}
			PaperProps={{
				sx: {
					display: 'flex',
					maxWidth: '400px',
					gap: theme.spacing(1),
					flexDirection: 'column',
					padding: theme.spacing(1),
				},
				onClick: (e) => {
					e.preventDefault();
					e.stopPropagation();
				},
			}}
		>
			<Box>
				<PopoverHeader handleClose={onClose} />
				<Divider />
				<ExistingAnnotations />
				<CreateAnnotation />
			</Box>
		</Popover>
	);
}

export default ManageTagAnnotations;
